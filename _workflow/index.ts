import * as cdk from "@aws-cdk/core";
import * as sfn from "@aws-cdk/aws-stepfunctions";

function pass(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Pass(stack, "Pass State");
}
function wait(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Wait(stack, "Wait State", {
    time: sfn.WaitTime.duration(cdk.Duration.seconds(1)),
  });
}
function succeed(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Succeed(stack, "Succeed State");
}
function fail(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Fail(stack, "Fail State");
}
function choice(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Choice(stack, "Choice State")
    .when(sfn.Condition.booleanEquals("$.bool", true), succeed(stack))
    .otherwise(fail(stack));
}
/*
function task(stack: cdk.Stack): sfn.IChainable {
  return undefined;
}
*/
function parallel(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Parallel(stack, "Parallel State")
    .branch(pass(stack))
    .branch(succeed(stack));
}
/*
function map(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Pass(stack, "Pass State");
}
*/

const workflows = {
  pass: pass,
  wait: wait,
  succeed: succeed,
  fail: fail,
  choice: choice,
  parallel: parallel,
};

function render(sm: sfn.IChainable) {
  return new cdk.Stack().resolve(
    new sfn.StateGraph(sm.startState, "Graph").toGraphJson()
  );
}

function print(sm: sfn.IChainable) {
  console.log(JSON.stringify(render(sm), null, "  "));
}

const args = process.argv.slice(2);
if (args.length == 0) {
  console.error("not enough args");
  process.exit(1);
}

const stack = new cdk.Stack();
for (const [key, wf] of Object.entries(workflows)) {
  if (key == args[0]) {
    print(wf(stack));
    process.exit(0);
  }
}

console.error("unknown key:", args[0]);
