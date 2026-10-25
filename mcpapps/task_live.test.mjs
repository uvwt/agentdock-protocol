import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import test from 'node:test';
import vm from 'node:vm';

const html=readFileSync(new URL('app.html',import.meta.url),'utf8');
const start=html.indexOf('function taskRecord');
const end=html.indexOf('function contextItems',start);

if(start<0||end<=start)throw new Error('Task live helpers are missing from app.html');

const context={
  isObject:value=>value!==null&&typeof value==='object'&&!Array.isArray(value),
};
vm.runInNewContext(
  `${html.slice(start,end)}
  globalThis.taskRecord=taskRecord;
  globalThis.taskCardAction=taskCardAction;
  globalThis.taskReview=taskReview;
  globalThis.taskRefreshCall=taskRefreshCall;
  globalThis.taskRefreshResult=taskRefreshResult;`,
  context,
);

const plain=value=>JSON.parse(JSON.stringify(value));

test('routes live task reads directly or through the original Nexus node',()=>{
  assert.deepEqual(
    plain(context.taskRefreshCall('tsk_1','')),
    {name:'task_manage',arguments:{action:'get',task_id:'tsk_1'}},
  );
  assert.deepEqual(
    plain(context.taskRefreshCall('tsk_1','node_1')),
    {name:'task_manage',arguments:{action:'get',task_id:'tsk_1',node_id:'node_1'}},
  );
});

test('renders task state instead of the internal get action',()=>{
  const data={action:'get',task:{id:'tsk_1',status:'active'}};
  assert.equal(context.taskCardAction(context.taskRecord(data)),'active');
});

test('reads final review from the authoritative task object',()=>{
  const data={
    action:'get',
    task:{id:'tsk_1',status:'active',final_review:{status:'pass',summary:'verified'}},
  };
  assert.deepEqual(plain(context.taskReview(context.taskRecord(data),data)),{status:'pass',summary:'verified'});
});

test('uses active and blocked refresh intervals and stops after completion',()=>{
  assert.deepEqual(
    plain(context.taskRefreshResult({structuredContent:{action:'get',task:{status:'active'}}})),
    {kind:'render',data:{action:'get',task:{status:'active'}},delay:2000,terminal:false},
  );
  assert.equal(
    context.taskRefreshResult({structuredContent:{action:'get',task:{status:'blocked'}}}).delay,
    10000,
  );
  assert.deepEqual(
    plain(context.taskRefreshResult({structuredContent:{action:'get',task:{status:'completed'}}})),
    {kind:'render',data:{action:'get',task:{status:'completed'}},delay:0,terminal:true},
  );
});

test('backs off retryable errors and stops non-retryable task reads',()=>{
  assert.deepEqual(
    plain(context.taskRefreshResult({isError:true,structuredContent:{code:'TEMPORARY',retryable:true}})),
    {kind:'retry',delay:5000},
  );
  assert.deepEqual(
    plain(context.taskRefreshResult({isError:true,structuredContent:{code:'TASK_NOT_FOUND',retryable:false}})),
    {kind:'stop',delay:5000},
  );
});
