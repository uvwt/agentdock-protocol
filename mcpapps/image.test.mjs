import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import vm from 'node:vm';
import { webcrypto } from 'node:crypto';

const html = readFileSync(new URL('./image.html', import.meta.url), 'utf8');
const source = html.match(/<script>([\s\S]*?)<\/script>/)[1];
const result = { content: [{ type:'image', mimeType:'image/png', data:'AQID' }] };

async function waitUntil(predicate, label='condition') {
  const deadline=Date.now()+2000;
  while(!predicate()) {
    if(Date.now()>=deadline) throw Error('timed out waiting for '+label);
    await new Promise(resolve=>setTimeout(resolve,1));
  }
}

function mount(api = {}) {
  const events={}, elements={}, messages=[], documentElement={lang:''};
  for (const id of ['preview','status','share']) elements[id]={hidden:id==='share',disabled:id==='share',textContent:'',addEventListener(name,cb){this[name]=cb;}};
  const parent={postMessage(message){messages.push(message);}};
  const win={parent,openai:api,addEventListener(name,cb){events[name]=cb;}};
  vm.runInNewContext(source,{window:win,document:{documentElement,getElementById:id=>elements[id]},atob,Uint8Array,File,crypto:webcrypto,console,setTimeout,clearTimeout});
  const dispatch=data=>events.message({source:parent,data});
  return {
    elements,api,messages,documentElement,
    initialize(hostCapabilities={},locale='en-US'){
      dispatch({jsonrpc:'2.0',id:'image-init',result:{protocolVersion:'2026-01-26',hostCapabilities,hostInfo:{name:'test',version:'1'},hostContext:{locale}}});
    },
    respond(requestMessage,resultValue={}){
      dispatch({jsonrpc:'2.0',id:requestMessage.id,result:resultValue});
    },
    reject(requestMessage){
      dispatch({jsonrpc:'2.0',id:requestMessage.id,error:{code:-32000,message:'denied'}});
    },
    requests(method){return messages.filter(message=>message.method===method && message.id!=='image-init');},
    globals(){events['openai:set_globals']();},
    receive(value,sender=parent){events.message({source:sender,data:{jsonrpc:'2.0',method:'ui/notifications/tool-result',params:value}});},
    hostContext(locale){dispatch({jsonrpc:'2.0',method:'ui/notifications/host-context-changed',params:{locale}});}
  };
}

test('uses standard model context plus ui/message for automatic image handoff', async () => {
  const ui=mount();
  ui.initialize({updateModelContext:{image:{}},message:{text:{}}});
  ui.receive(result);
  await waitUntil(()=>ui.requests('ui/update-model-context').length===1,'context update');
  const update=ui.requests('ui/update-model-context')[0];
  assert.deepEqual(JSON.parse(JSON.stringify(update.params.content)),[{type:'image',data:'AQID',mimeType:'image/png'}]);
  ui.respond(update);
  await waitUntil(()=>ui.requests('ui/message').length===1,'follow-up message');
  const follow=ui.requests('ui/message')[0];
  assert.equal(follow.params.role,'user');
  assert.equal(follow.params.content[0].type,'text');
  ui.respond(follow);
  await waitUntil(()=>ui.elements.status.textContent==='Image is available to the model.','handoff completion');
  assert.equal(ui.elements.share.hidden,true);
});

test('uses direct standard image message when model-context image updates are unavailable', async () => {
  const ui=mount();
  ui.initialize({message:{text:{},image:{}}});
  ui.receive(result);
  await waitUntil(()=>ui.requests('ui/message').length===1,'image message');
  const message=ui.requests('ui/message')[0];
  assert.equal(message.params.content[0].type,'text');
  assert.deepEqual(JSON.parse(JSON.stringify(message.params.content[1])),{type:'image',data:'AQID',mimeType:'image/png'});
  ui.respond(message);
  await waitUntil(()=>ui.elements.status.textContent==='Image is available to the model.','direct handoff completion');
});

test('uses standard image context with ChatGPT follow-up without uploading a file', async () => {
  let follows=0,uploads=0;
  const ui=mount({
    uploadFile:async()=>{uploads++;return {fileId:'file-real'};},
    setWidgetState:async()=>{},
    sendFollowUpMessage:async()=>{follows++;}
  });
  ui.initialize({updateModelContext:{image:{}}});
  ui.receive(result);
  await waitUntil(()=>ui.requests('ui/update-model-context').length===1,'context update');
  ui.respond(ui.requests('ui/update-model-context')[0]);
  await waitUntil(()=>follows===1,'ChatGPT follow-up');
  assert.equal(uploads,0);
});

test('falls back to ChatGPT file attachment and waits for widget state persistence', async () => {
  const calls=[]; let release;
  const ui=mount({
    uploadFile:async file=>{calls.push(['upload',Array.from(new Uint8Array(await file.arrayBuffer())),file.type]);return {fileId:'file-real'};},
    setWidgetState:()=>new Promise(resolve=>{calls.push(['state-start']);release=()=>{calls.push(['state-done']);resolve();};}),
    sendFollowUpMessage:async()=>calls.push(['follow'])
  });
  ui.initialize({});
  ui.receive(result);
  await waitUntil(()=>typeof release==='function','widget state persistence');
  assert.deepEqual(calls[0],['upload',[1,2,3],'image/png']);
  assert.deepEqual(calls[1],['state-start']);
  assert.equal(calls.some(c=>c[0]==='follow'),false);
  release();
  await waitUntil(()=>calls.some(c=>c[0]==='follow'),'fallback follow-up');
  assert.deepEqual(calls.slice(1),[['state-start'],['state-done'],['follow']]);
});

test('automatic handoff waits for initialization and starts when capabilities arrive', async () => {
  const ui=mount();
  ui.receive(result);
  await waitUntil(()=>ui.elements.preview.src==='data:image/png;base64,AQID','preview');
  assert.equal(ui.requests('ui/message').length,0);
  ui.initialize({message:{image:{}}});
  await waitUntil(()=>ui.requests('ui/message').length===1,'post-init image message');
});

test('repeated globals can activate a newly available ChatGPT fallback without duplicating the image generation', async () => {
  const api={toolResponseMetadata:{mcp_tool_result:result}};
  const ui=mount(api);
  ui.initialize({});
  await waitUntil(()=>ui.elements.preview.src==='data:image/png;base64,AQID','preview');
  assert.match(ui.elements.status.textContent,/preview/i);
  let follows=0;
  api.uploadFile=async()=>({fileId:'file-real'});
  api.setWidgetState=async()=>{};
  api.sendFollowUpMessage=async()=>{follows++;};
  ui.globals();
  await waitUntil(()=>follows===1,'late ChatGPT fallback');
});

test('host request errors expose a retry instead of claiming success', async () => {
  const ui=mount();
  ui.initialize({message:{image:{}}});
  ui.receive(result);
  await waitUntil(()=>ui.requests('ui/message').length===1,'image message');
  ui.reject(ui.requests('ui/message')[0]);
  await waitUntil(()=>ui.elements.status.textContent.includes('failed'),'handoff failure');
  assert.equal(ui.elements.share.hidden,false);
  assert.equal(ui.elements.share.disabled,false);
});

test('untrusted frames, tool errors and active content are not images', () => {
  const ui=mount();
  ui.initialize({});
  ui.receive(result,{});
  assert.equal(ui.elements.preview.src,undefined);
  ui.receive({isError:true,...result});
  assert.match(ui.elements.status.textContent,/error/i);
  ui.receive({content:[{type:'image',mimeType:'image/svg+xml',data:'AQID'}]});
  assert.equal(ui.elements.preview.src,undefined);
});

test('switching images during fallback upload never attaches stale pixels', async () => {
  const finishes=[],states=[];
  const ui=mount({
    uploadFile:()=>new Promise(resolve=>finishes.push(resolve)),
    setWidgetState:async state=>states.push(Array.from(state.imageIds))
  });
  ui.initialize({});
  ui.receive(result);
  await waitUntil(()=>finishes.length===1,'first upload');
  ui.receive({content:[{type:'image',mimeType:'image/png',data:'BAUG'}]});
  await waitUntil(()=>ui.elements.preview.src==='data:image/png;base64,BAUG','replacement image');
  finishes[0]({fileId:'file-old'});
  await waitUntil(()=>finishes.length===2,'replacement upload');
  finishes[1]({fileId:'file-new'});
  await waitUntil(()=>states.length===1,'replacement state');
  assert.deepEqual(states,[['file-new']]);
});

test('localizes the automatic handoff UI from host context', () => {
  const ui=mount();
  ui.initialize({},'zh-CN');
  assert.equal(ui.documentElement.lang,'zh-CN');
  assert.equal(ui.elements.share.textContent,'重试提供图片');
  ui.hostContext('en-US');
  assert.equal(ui.documentElement.lang,'en');
  assert.equal(ui.elements.share.textContent,'Retry image handoff');
});

test('ships a restrictive CSP for the shared image document', () => {
  assert.match(html,/Content-Security-Policy/);
  assert.match(html,/default-src 'none'/);
  assert.match(html,/connect-src 'none'/);
  assert.match(html,/img-src data:/);
});
