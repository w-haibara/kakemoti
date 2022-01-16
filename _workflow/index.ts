import { Stack, Duration, aws_stepfunctions as sfn } from "aws-cdk-lib";
import { ScriptTask } from "./script-task.js";
const path = require("path");
const fs = require("fs");

function pass(stack: Stack): sfn.IChainable {
  return new sfn.Pass(stack, "Pass State");
}
function pass_chain(stack: Stack): sfn.IChainable {
  const p1 = new sfn.Pass(stack, "P1");
  const p2 = new sfn.Pass(stack, "P2");
  const p3 = new sfn.Pass(stack, "P3");
  const p4 = new sfn.Pass(stack, "P4");
  const p5 = new sfn.Pass(stack, "P5");
  return p1.next(p2).next(p3).next(p4).next(p5);
}
function pass_result(stack: Stack): sfn.IChainable {
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
function pass_parameters(stack: Stack): sfn.IChainable {
  return new sfn.Pass(stack, "Pass State(parameter)", {
    parameters: {
      parameters: {
        aaa: 111,
        bbb: 222,
      },
    },
  });
}
function pass_intrinsic(stack: Stack): sfn.IChainable {
  return new sfn.Pass(stack, "Pass State(intrinsic)", {
    parameters: {
      parameters: {
        aaa: 111,
        intrinsic: {
          "format.$": "States.Format('Hello, my name is {}.', $.name)",
          "stringToJson.$": "States.StringToJson($.string)",
          "jsonToString.$": "States.JsonToString($.json)",
          "array.$": "States.Array('start', $.json.aaa, $.json.bbb, 'end')",
        },
      },
    },
  });
}
function wait(stack: Stack): sfn.IChainable {
  return new sfn.Wait(stack, "Wait State", {
    time: sfn.WaitTime.duration(Duration.seconds(1)),
  });
}
function succeed(stack: Stack): sfn.IChainable {
  return new sfn.Succeed(stack, "Succeed State");
}
function fail(stack: Stack): sfn.IChainable {
  return new sfn.Fail(stack, "Fail State");
}
function choice(stack: Stack): sfn.IChainable {
  return new sfn.Choice(stack, "Choice State")
    .when(sfn.Condition.booleanEquals("$.bool", true), succeed(stack))
    .otherwise(fail(stack));
}
function choice_fallback(stack: Stack): sfn.IChainable {
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
function choice_bool(stack: Stack): sfn.IChainable {
  const ok = new sfn.Pass(stack, "OK", {
    result: sfn.Result.fromObject([
      {
        args: {
          Payload: "OK",
        },
      },
    ]),
    resultPath: sfn.JsonPath.DISCARD,
  });
  const ng = new sfn.Pass(stack, "NG", {
    result: sfn.Result.fromObject([
      {
        args: {
          Payload: "NG",
        },
      },
    ]),
    resultPath: sfn.JsonPath.DISCARD,
  });

  const t = sfn.Condition.booleanEquals("$.true", true); // true
  const f = sfn.Condition.booleanEquals("$.false", true); // false

  const cond1 = sfn.Condition.and(t, f); // false
  const cond2 = sfn.Condition.or(t, f); // true
  const cond3 = sfn.Condition.not(cond1); // true
  const cond4 = sfn.Condition.and(cond2, cond3); // true
  const cond5 = sfn.Condition.and(cond2, cond3, cond4); // true
  return new sfn.Choice(stack, "Choice State").when(cond5, ok).otherwise(ng);
}
function choice_data_test(stack: Stack): sfn.IChainable {
  const ok = new sfn.Pass(stack, "OK", {
    result: sfn.Result.fromObject([
      {
        args: {
          Payload: "OK",
        },
      },
    ]),
    resultPath: sfn.JsonPath.DISCARD,
  });
  const ng = new sfn.Pass(stack, "NG", {
    result: sfn.Result.fromObject([
      {
        args: {
          Payload: "NG",
        },
      },
    ]),
    resultPath: sfn.JsonPath.DISCARD,
  });

  const cond1 = sfn.Condition.booleanEquals("$.bool", true);
  const cond2 = sfn.Condition.not(
    sfn.Condition.booleanEqualsJsonPath("$.bool", "$.object.bool")
  );
  const cond3 = sfn.Condition.isBoolean("$.bool");
  const cond4 = sfn.Condition.not(sfn.Condition.isBoolean("$.string"));
  const cond5 = sfn.Condition.isNull("$.null");
  const cond6 = sfn.Condition.not(sfn.Condition.isNull("$.string"));
  const cond7 = sfn.Condition.isNumeric("$.int");
  const cond8 = sfn.Condition.not(sfn.Condition.isNumeric("$.string"));
  return new sfn.Choice(stack, "Choice State")
    .when(sfn.Condition.and(cond1, cond2, cond3, cond4, cond5, cond6, cond7, cond8), ok)
    .otherwise(ng);
}
function task(stack: Stack): sfn.IChainable {
  return new ScriptTask(stack, "Task State", {
    scriptPath: "_workflow/script/script1.sh",
  });
}
function task_filter(stack: Stack): sfn.IChainable {
  return new ScriptTask(stack, "Task State", {
    scriptPath: "_workflow/script/script1.sh",
    inputPath: "$.inputpath",
    parameters: sfn.TaskInput.fromObject({
      aaa: 111,
      "old.$": "$.args",
      args: ["param0", "param1", "param2"],
    }),
    resultSelector: {
      bbb: 222,
      "resultselector.$": "$",
    },
    resultPath: "$.resultpath.outputpath",
    outputPath: "$.resultpath",
  });
}
function task_retry(stack: Stack): sfn.IChainable {
  const task = new ScriptTask(stack, "Task State", {
    scriptPath: "_workflow/script/script2.sh",
    resultPath: "$.args",
  });
  const chain = new sfn.Parallel(stack, "Chain").branch(task);
  chain.addRetry({
    maxAttempts: 10,
    backoffRate: 0,
    interval: Duration.seconds(0),
  });
  return chain;
}
function task_catch(stack: Stack): sfn.IChainable {
  const p1 = new sfn.Pass(stack, "Pass State1");
  const task = new ScriptTask(stack, "Task State", {
    scriptPath: "::", // invalid resource path
  });
  task.addCatch(p1, {
    errors: ["States.ALL"],
  });
  return task;
}
function task_ctx(stack: Stack): sfn.IChainable {
  return new ScriptTask(stack, "Task State", {
    scriptPath: "_workflow/script/script1.sh",
    resultSelector: {
      ctx: {
        "ctx_aaa.$": "$$.aaa",
      },
    },
  });
}
function parallel(stack: Stack): sfn.IChainable {
  return new sfn.Parallel(stack, "Parallel State")
    .branch(pass(stack))
    .branch(succeed(stack));
}
/*
function map(stack: Stack): sfn.IChainable {
  return new sfn.Pass(stack, "Pass State");
}
*/

const workflows: ((stack: Stack) => sfn.IChainable)[] = [
  pass,
  pass_chain,
  pass_result,
  pass_parameters,
  pass_intrinsic,
  wait,
  succeed,
  fail,
  choice,
  choice_fallback,
  choice_bool,
  choice_data_test,
  task,
  task_filter,
  task_retry,
  task_catch,
  task_ctx,
  parallel,
];

function list() {
  workflows.forEach(function (elm) {
    console.log(elm.name);
  });
}

function render(sm: sfn.IChainable) {
  return new Stack().resolve(
    new sfn.StateGraph(sm.startState, "Graph").toGraphJson()
  );
}

function toString(sm: sfn.IChainable): string {
  return JSON.stringify(render(sm), null, "  ");
}

const args = process.argv.slice(2);
if (args.length == 0) {
  console.error("not enough args");
  process.exit(1);
}

if (args[0] == "list") {
  list();
  process.exit(0);
}

const filename = path.basename(args[0]);
const asl = path.basename(args[0], ".asl.json");
const abspath = path.join(__dirname, "asl", filename);

const stack = new Stack();
const fn = workflows.find((elm) => elm.name == asl);
if (fn === undefined) {
  console.error("unknown key:", asl);
} else {
  fs.writeFile(abspath, toString(fn(stack)), function (err: any) {
    if (err) {
      console.log(err);
    } else {
      console.log("write to", args[0]);
    }
  });
}
