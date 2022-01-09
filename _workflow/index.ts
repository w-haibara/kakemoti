import * as cdk from "@aws-cdk/core";
import * as sfn from "@aws-cdk/aws-stepfunctions";
import * as custom from "./custom-task.js";

function pass(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Pass(stack, "Pass State");
}
function pass_chain(stack: cdk.Stack): sfn.IChainable {
  const p1 = new sfn.Pass(stack, "P1");
  const p2 = new sfn.Pass(stack, "P2");
  const p3 = new sfn.Pass(stack, "P3");
  const p4 = new sfn.Pass(stack, "P4");
  const p5 = new sfn.Pass(stack, "P5");
  return p1.next(p2).next(p3).next(p4).next(p5);
}
function pass_result(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Pass(stack, "Pass State(result)", {
    result: sfn.Result.fromObject({
      result: {
        aaa: 111,
        bbb: 222,
      },
    }),
    resultPath: "$.resultpath",
  });
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
function choice_fallback(stack: cdk.Stack): sfn.IChainable {
  const s1 = new sfn.Pass(stack, "State1", {
    result: sfn.Result.fromObject({
      bool: false,
    }),
  });
  const s2 = new sfn.Pass(stack, "State2");
  const s3 = new sfn.Pass(stack, "State3");
  const pass = s1.next(s2);
  const choice = new sfn.Choice(stack, "Choice State")
    .when(sfn.Condition.booleanEquals("$.bool", false), s3)
    .otherwise(pass);
  return s2.next(choice);
}
function task(stack: cdk.Stack): sfn.IChainable {
  return new custom.Task(stack, "Task State", {
    resource: "script:_workflow/script/script1.sh",
  });
}
function task_resultPath(stack: cdk.Stack): sfn.IChainable {
  return new custom.Task(stack, "Task State", {
    resource: "script:_workflow/script/script1.sh",
    resultPath: "$.resultpath",
  });
}
function task_retry(stack: cdk.Stack): sfn.IChainable {
  const task = new custom.Task(stack, "Task State", {
    resource: "script:_workflow/script/script2.sh",
    resultPath: "$.args",
  });
  const chain = new sfn.Parallel(stack, "Chain").branch(task);
  chain.addRetry({
    maxAttempts: 10,
    backoffRate: 0,
    interval: cdk.Duration.seconds(0),
  });
  return chain;
}
function task_catch(stack: cdk.Stack): sfn.IChainable {
  const p1 = new sfn.Pass(stack, "Pass State1");
  const task = new custom.Task(stack, "Task State", {
    resource: "script:...", // invalid resource path
  });
  task.addCatch(p1, {
    errors: ["States.ALL"],
  });
  return task;
}
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
  pass_chain: pass_chain,
  pass_result: pass_result,
  wait: wait,
  succeed: succeed,
  fail: fail,
  choice: choice,
  choice_fallback: choice_fallback,
  task: task,
  task_resultPath: task_resultPath,
  task_retry: task_retry,
  task_catch: task_catch,
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
