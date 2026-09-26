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

function element(id, hidden=false) {
  const attrs=new Map();
  const classes=new Set();
  return {
    id,
    hidden,
    disabled:id==='share',
    textContent:'',
    naturalWidth:0,
    naturalHeight:0,
    style:{},
    addEventListener(name,cb){this[name]=cb;},
    setAttribute(name,value){attrs.set(name,String(value));},
    getAttribute(name){return attrs.get(name)??null;},
    classList:{
      toggle(name,force){
        const next=force===undefined?!classes.has(name):Boolean(force);
        if(next) classes.add(name); else classes.delete(name);
        return next;
      },
      contains(name){return classes.has(name);}
    },
    getBoundingClientRect(){return {height:68};}
  };
}

function mount(api = {}) {
  const events={},messages=[];
  const elements={
    shell:element('shell'),
    toggle:element('toggle'),
    metadata:element('metadata'),
    panel:element('panel',true),
    preview:element('preview',true),
    share:element('share',true)
  };
  elements.toggle.setAttribute('aria-expanded','false');
  elements.shell.getBoundingClientRect=()=>({height:elements.panel.hidden?68:420});

  const style={setProperty(){},colorScheme:''};
  const documentElement={lang:'',dataset:{},style,scrollHeight:0};
  const parent={postMessage(message){messages.push(message);}};
  const win={parent,openai:api,addEventListener(name,cb){events[name]=cb;}};
  const document={
    documentElement,
    body:{scrollHeight:0},
    getElementById:id=>elements[id]
  };

  vm.runInNewContext(source,{
    window:win,document,atob,Uint8Array,File,crypto:webcrypto,console,setTimeout,clearTimeout
  });

  const dispatch=data=>events.message({source:parent,data});
  return {
    elements,api,messages,documentElement,
    initialize(hostCapabilities={},locale='en-US',extraContext={}){
      dispatch({jsonrpc:'2.0',id:'image-init',result:{
        protocolVersion:'2026-01-26',
        hostCapabilities,
        hostInfo:{name:'test',version:'1'},
        hostContext:{locale,...extraContext}
      }});
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
    hostContext(params){dispatch({jsonrpc:'2.0',method:'ui/notifications/host-context-changed',params});},
    expand(){elements.toggle.click();},
    collapse(){elements.toggle.click();}
  };
}

async function loadImage(ui,value=result,width=320,height=240) {
  ui.receive(value);
  const expected=value.content.find(item=>item.type==='image')?.data;
  if(expected) await waitUntil(()=>ui.elements.preview.src?.endsWith(expected),'preview');
  ui.elements.preview.naturalWidth=width;
  ui.elements.preview.naturalHeight=height;
  ui.elements.preview.onload?.();
  await new Promise(resolve=>setTimeout(resolve,1));
}

test('uses the shared compact card shape and stays folded until the user expands it', async () => {
  const ui=mount();
  ui.initialize({updateModelContext:{image:{}},message:{text:{}}});
  await loadImage(ui);

  assert.match(html,/class="compact-title">Image<\/span>/);
  assert.match(html,/class="brand">AgentDock<\/span>/);
  assert.equal(ui.elements.toggle.getAttribute('aria-expanded'),'false');
  assert.equal(ui.elements.panel.hidden,true);
  assert.equal(ui.elements.share.hidden,true);
  assert.equal(ui.elements.metadata.textContent,'320 × 240 · PNG · Awaiting confirmation');
  assert.equal(ui.requests('ui/update-model-context').length,0);
  assert.equal(ui.requests('ui/message').length,0);

  ui.expand();
  assert.equal(ui.elements.toggle.getAttribute('aria-expanded'),'true');
  assert.equal(ui.elements.toggle.classList.contains('expanded'),true);
  assert.equal(ui.elements.panel.hidden,false);
  assert.equal(ui.elements.share.hidden,false);
  assert.equal(ui.elements.share.disabled,false);

  ui.collapse();
  assert.equal(ui.elements.toggle.getAttribute('aria-expanded'),'false');
  assert.equal(ui.elements.panel.hidden,true);
  assert.equal(ui.elements.share.hidden,true);
});

test('waits for expanded user confirmation before standard model context plus ui/message handoff', async () => {
  const ui=mount();
  ui.initialize({updateModelContext:{image:{}},message:{text:{}}});
  await loadImage(ui);
  ui.expand();

  const action=ui.elements.share.click();
  await waitUntil(()=>ui.requests('ui/update-model-context').length===1,'context update');
  const update=ui.requests('ui/update-model-context')[0];
  assert.deepEqual(JSON.parse(JSON.stringify(update.params.content)),[{type:'image',data:'AQID',mimeType:'image/png'}]);
  ui.respond(update);
  await waitUntil(()=>ui.requests('ui/message').length===1,'follow-up message');
  const follow=ui.requests('ui/message')[0];
  assert.equal(follow.params.role,'user');
  assert.equal(follow.params.content[0].type,'text');
  ui.respond(follow);
  await action;

  assert.equal(ui.elements.metadata.textContent,'320 × 240 · PNG · Provided to model');
  assert.equal(ui.elements.share.hidden,true);
});

test('uses direct standard image message only after user confirmation', async () => {
  const ui=mount();
  ui.initialize({message:{text:{},image:{}}});
  await loadImage(ui);
  assert.equal(ui.requests('ui/message').length,0);
  ui.expand();

  const action=ui.elements.share.click();
  await waitUntil(()=>ui.requests('ui/message').length===1,'image message');
  const message=ui.requests('ui/message')[0];
  assert.equal(message.params.content[0].type,'text');
  assert.deepEqual(JSON.parse(JSON.stringify(message.params.content[1])),{type:'image',data:'AQID',mimeType:'image/png'});
  ui.respond(message);
  await action;
  assert.match(ui.elements.metadata.textContent,/Provided to model$/);
});

test('prefers standard image context over ChatGPT file upload after user confirmation', async () => {
  let follows=0,uploads=0;
  const ui=mount({
    uploadFile:async()=>{uploads++;return {fileId:'file-real'};},
    setWidgetState:async()=>{},
    sendFollowUpMessage:async()=>{follows++;}
  });
  ui.initialize({updateModelContext:{image:{}}});
  await loadImage(ui);
  assert.equal(ui.requests('ui/update-model-context').length,0);
  assert.equal(follows,0);
  assert.equal(uploads,0);
  ui.expand();

  const action=ui.elements.share.click();
  await waitUntil(()=>ui.requests('ui/update-model-context').length===1,'context update');
  ui.respond(ui.requests('ui/update-model-context')[0]);
  await action;
  assert.equal(follows,1);
  assert.equal(uploads,0);
});

test('falls back to ChatGPT file attachment only after user confirmation and waits for widget state persistence', async () => {
  const calls=[]; let release;
  const ui=mount({
    uploadFile:async file=>{calls.push(['upload',Array.from(new Uint8Array(await file.arrayBuffer())),file.type]);return {fileId:'file-real'};},
    setWidgetState:()=>new Promise(resolve=>{calls.push(['state-start']);release=()=>{calls.push(['state-done']);resolve();};}),
    sendFollowUpMessage:async()=>calls.push(['follow'])
  });
  ui.initialize({});
  await loadImage(ui);
  assert.deepEqual(calls,[]);
  ui.expand();

  const action=ui.elements.share.click();
  await waitUntil(()=>typeof release==='function','widget state persistence');
  assert.deepEqual(calls[0],['upload',[1,2,3],'image/png']);
  assert.deepEqual(calls[1],['state-start']);
  assert.equal(calls.some(c=>c[0]==='follow'),false);
  release();
  await action;
  assert.deepEqual(calls.slice(1),[['state-start'],['state-done'],['follow']]);
});

test('image can arrive before initialization without starting a handoff', async () => {
  const ui=mount();
  await loadImage(ui);
  assert.equal(ui.elements.panel.hidden,true);
  assert.equal(ui.elements.share.hidden,true);
  assert.equal(ui.requests('ui/message').length,0);

  ui.initialize({message:{image:{}}});
  assert.equal(ui.elements.share.hidden,true);
  ui.expand();
  await waitUntil(()=>ui.elements.share.hidden===false,'post-init confirmation button');
  assert.equal(ui.requests('ui/message').length,0);

  const action=ui.elements.share.click();
  await waitUntil(()=>ui.requests('ui/message').length===1,'post-click image message');
  ui.respond(ui.requests('ui/message')[0]);
  await action;
});

test('late ChatGPT fallback changes folded metadata without automatic upload', async () => {
  const api={toolResponseMetadata:{mcp_tool_result:result}};
  const ui=mount(api);
  ui.initialize({});
  await waitUntil(()=>ui.elements.preview.src==='data:image/png;base64,AQID','preview');
  ui.elements.preview.naturalWidth=320;
  ui.elements.preview.naturalHeight=240;
  ui.elements.preview.onload();
  await waitUntil(()=>/Preview only$/.test(ui.elements.metadata.textContent),'preview-only metadata');

  let follows=0,uploads=0;
  api.uploadFile=async()=>{uploads++;return {fileId:'file-real'};};
  api.setWidgetState=async()=>{};
  api.sendFollowUpMessage=async()=>{follows++;};
  ui.globals();

  await waitUntil(()=>/Awaiting confirmation$/.test(ui.elements.metadata.textContent),'ready metadata');
  assert.equal(ui.elements.panel.hidden,true);
  assert.equal(ui.elements.share.hidden,true);
  assert.equal(uploads,0);
  assert.equal(follows,0);

  ui.expand();
  await waitUntil(()=>ui.elements.share.hidden===false,'late confirmation button');
  await ui.elements.share.click();
  assert.equal(uploads,1);
  assert.equal(follows,1);
});

test('host request errors expose retry only inside the expanded panel', async () => {
  const ui=mount();
  ui.initialize({message:{image:{}}});
  await loadImage(ui);
  ui.expand();

  const action=ui.elements.share.click();
  await waitUntil(()=>ui.requests('ui/message').length===1,'image message');
  ui.reject(ui.requests('ui/message')[0]);
  await action;

  assert.match(ui.elements.metadata.textContent,/Unable to provide image$/);
  assert.equal(ui.elements.share.textContent,'Try again');
  assert.equal(ui.elements.share.hidden,false);
  assert.equal(ui.elements.share.disabled,false);

  ui.collapse();
  assert.equal(ui.elements.share.hidden,true);
});

test('untrusted frames, tool errors and active content are not images', async () => {
  const ui=mount();
  ui.initialize({});
  ui.receive(result,{});
  assert.equal(ui.elements.preview.src,undefined);
  ui.receive({isError:true,...result});
  await waitUntil(()=>/Image unavailable$/.test(ui.elements.metadata.textContent),'tool error');
  ui.receive({content:[{type:'image',mimeType:'image/svg+xml',data:'AQID'}]});
  await waitUntil(()=>/Image unavailable$/.test(ui.elements.metadata.textContent),'invalid image');
  assert.equal(ui.elements.preview.src,undefined);
});

test('switching images during fallback upload never attaches replacement pixels automatically', async () => {
  const finishes=[],states=[];
  const ui=mount({
    uploadFile:()=>new Promise(resolve=>finishes.push(resolve)),
    setWidgetState:async state=>states.push(Array.from(state.imageIds))
  });
  ui.initialize({});
  await loadImage(ui);
  ui.expand();

  const first=ui.elements.share.click();
  await waitUntil(()=>finishes.length===1,'first upload');
  await loadImage(ui,{content:[{type:'image',mimeType:'image/png',data:'BAUG'}]},640,480);
  finishes[0]({fileId:'file-old'});
  await first;

  assert.deepEqual(states,[]);
  assert.equal(finishes.length,1);
  assert.equal(ui.elements.metadata.textContent,'640 × 480 · PNG · Awaiting confirmation');
  await waitUntil(()=>ui.elements.share.hidden===false,'replacement confirmation button');

  const second=ui.elements.share.click();
  await waitUntil(()=>finishes.length===2,'replacement upload');
  finishes[1]({fileId:'file-new'});
  await second;
  assert.deepEqual(states,[['file-new']]);
});

test('keeps Image untranslated while localizing only status and action text', async () => {
  const ui=mount();
  ui.initialize({message:{image:{}}},'zh-CN');
  await loadImage(ui);

  assert.equal(ui.documentElement.lang,'zh-CN');
  assert.match(html,/class="compact-title">Image<\/span>/);
  assert.equal(ui.elements.metadata.textContent,'320 × 240 · PNG · 等待确认');
  assert.equal(ui.elements.panel.hidden,true);

  ui.expand();
  assert.equal(ui.elements.share.textContent,'让模型查看图片');

  ui.hostContext({locale:'en-US'});
  assert.equal(ui.documentElement.lang,'en');
  assert.equal(ui.elements.metadata.textContent,'320 × 240 · PNG · Awaiting confirmation');
  assert.equal(ui.elements.share.textContent,'Let model inspect image');
});

test('ships the same transparent compact-shell styling and restrictive CSP', () => {
  assert.match(html,/body\{margin:0;padding:0;background:transparent/);
  assert.match(html,/\.compact-toggle\{width:100%;min-height:68px;border:0;border-radius:0/);
  assert.match(html,/\.detail-panel\{border-top:1px solid var\(--ad-border\)/);
  assert.doesNotMatch(html,/preview-shell/);
  assert.doesNotMatch(html,/height:580/);
  assert.match(html,/Content-Security-Policy/);
  assert.match(html,/default-src 'none'/);
  assert.match(html,/connect-src 'none'/);
  assert.match(html,/img-src data:/);
  assert.match(html,/modelContent:"The image is available in the current context\. Use the image itself when answering the user\'s request\."/);
  assert.match(html,/followup:"Inspect the image in the current context and continue answering the user\'s request about it\. If the user did not ask a specific question, briefly describe the image\. Do not call view_image again\."/);
  assert.doesNotMatch(html,/previous user message|substitute OCR|Retry image handoff|Handoff failed|Image decode failed/);
});
