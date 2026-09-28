import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import test from 'node:test';
import vm from 'node:vm';

const html=readFileSync(new URL('app.html',import.meta.url),'utf8');
const start=html.indexOf('function acpTextContent');
const end=html.indexOf('function renderACP',start);

if(start<0||end<=start)throw new Error('ACP transcript helpers are missing from app.html');

const context={
  isObject:value=>value!==null&&typeof value==='object'&&!Array.isArray(value),
  scalar:value=>value===null||['string','number','boolean'].includes(typeof value),
  t:(key,vars={})=>key==='moreItems'?`${vars.count} more`:key,
};
vm.runInNewContext(`${html.slice(start,end)};globalThis.acpHistoryMessages=acpHistoryMessages;globalThis.acpNamedValues=acpNamedValues;globalThis.acpChangeSummary=acpChangeSummary;globalThis.acpPromptText=acpPromptText;`,context);

const project=events=>JSON.parse(JSON.stringify(context.acpHistoryMessages(events)));

test('projects ACP history chunks into user and assistant messages',()=>{
  const messages=project([
    {type:'user_message_chunk',update:{messageId:'user-1',content:{type:'text',text:'hello '}}},
    {type:'user_message_chunk',update:{messageId:'user-1',content:{type:'text',text:'world'}}},
    {type:'agent_thought_chunk',update:{content:{type:'text',text:'private reasoning'}}},
    {type:'agent_message_chunk',update:{content:{type:'text',text:'answer '}}},
    {type:'agent_message_chunk',update:{content:'done'}},
    {type:'agent_message_chunk',update:{content:{type:'image',data:'ignored'}}},
  ]);

  assert.deepEqual(messages,[
    {role:'user',content:'hello world'},
    {role:'assistant',content:'answer done'},
  ]);
});

test('keeps consecutive messages with distinct message ids separate',()=>{
  const messages=project([
    {type:'agent_message_chunk',update:{messageId:'assistant-1',content:{type:'text',text:'first'}}},
    {type:'agent_message_chunk',update:{messageId:'assistant-2',content:{type:'text',text:'second'}}},
  ]);

  assert.deepEqual(messages,[
    {role:'assistant',content:'first'},
    {role:'assistant',content:'second'},
  ]);
});

test('summarizes ACP configuration names with their current values',()=>{
  const summary=context.acpNamedValues([
    {id:'model',name:'Model',currentValue:'gpt-5'},
    {id:'safe',name:'Safe mode',currentValue:false},
    {optionId:'review',name:'Review'},
  ]);

  assert.equal(summary,'Model: gpt-5 · Safe mode: false · Review');
});

test('summarizes an ACP update with before and after values',()=>{
  const summary=context.acpChangeSummary({
    field:'config_option',id:'model',label:'Model',before:'gpt-4.1',after:'gpt-5',
  });

  assert.equal(summary,'Model: gpt-4.1 → gpt-5');
});

test('summarizes ACP prompt content blocks for the prompt card',()=>{
  const summary=context.acpPromptText([
    {type:'text',text:'请检查这个登录流程。'},
    {type:'image',mimeType:'image/png',data:'AA=='},
  ]);

  assert.equal(summary,'请检查这个登录流程。\n[image]');
});
