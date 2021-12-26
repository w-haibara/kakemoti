import * as cdk from "@aws-cdk/core";
import * as sfn from "@aws-cdk/aws-stepfunctions";

function pass(stack: cdk.Stack): sfn.IChainable {
  return new sfn.Pass(stack, "Pass State 1");
}

const workflows = {
  pass: pass,
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
