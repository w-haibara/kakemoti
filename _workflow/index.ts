import * as cdk from "@aws-cdk/core";
import * as sfn from "@aws-cdk/aws-stepfunctions";

function render(sm: sfn.IChainable) {
  return new cdk.Stack().resolve(
    new sfn.StateGraph(sm.startState, "Graph").toGraphJson()
  );
}

function case1(): sfn.Chain {
  const stack = new cdk.Stack();
  const pass1 = new sfn.Pass(stack, "Pass State 1");
  const pass2 = new sfn.Pass(stack, "Pass State 2");
  return pass1.next(pass2)
  /*
  const stack = new cdk.Stack();
  const pass = new sfn.Pass(stack, "Pass State");
  const succeed = new sfn.Succeed(stack, "Succeed State");
  const fail = new sfn.Fail(stack, "Fail State");
  const parallel = new sfn.Parallel(stack, "Parallel State")
    .branch(succeed)
    .branch(fail);
  return pass.next(parallel);
  */
}

const args = process.argv.slice(2);
if (args.length == 0) {
  console.error("not enough args");
  process.exit(1);
}

switch (args[0]) {
  case "case1":
    console.log(JSON.stringify(render(case1()), null, "  "));
    break;
  default:
    console.error("unknown type");
    process.exit(1);
}
