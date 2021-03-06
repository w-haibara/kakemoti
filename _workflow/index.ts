import { Stack, Duration, aws_stepfunctions as sfn } from "aws-cdk-lib";
import { ScriptTask } from "./script-task.js";
const path = require("path");
const fs = require("fs");

type def = [sfn.IChainable, number];

function pass(stack: Stack): def {
  return [new sfn.Pass(stack, "Pass State"), 0];
}
function pass_chain(stack: Stack): def {
  const p1 = new sfn.Pass(stack, "P1");
  const p2 = new sfn.Pass(stack, "P2");
  const p3 = new sfn.Pass(stack, "P3");
  const p4 = new sfn.Pass(stack, "P4");
  const p5 = new sfn.Pass(stack, "P5");
  const v = p1.next(p2).next(p3).next(p4).next(p5);
  return [v, 0];
}
function pass_result(stack: Stack): def {
  const v = new sfn.Pass(stack, "Pass State(result)", {
    result: sfn.Result.fromObject({
      result: {
        aaa: 111,
        bbb: 222,
      },
    }),
    resultPath: "$.resultpath",
  });
  return [v, 0];
}
function pass_parameters(stack: Stack): def {
  const v = new sfn.Pass(stack, "Pass State(parameter)", {
    parameters: {
      parameters: {
        aaa: 111,
        bbb: 222,
      },
    },
  });
  return [v, 0];
}
function pass_intrinsic(stack: Stack): def {
  const v = new sfn.Pass(stack, "Pass State(intrinsic)", {
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
  return [v, 0];
}
function pass_ctxobj(stack: Stack): def {
  const v = new sfn.Pass(stack, "Pass State", {
    parameters: {
      "ctx_aaa.$": "$$.aaa",
    },
    resultPath: "$.result.ctx",
    outputPath: "$.result",
  });
  return [v, 0];
}
function wait(stack: Stack): def {
  const v = new sfn.Wait(stack, "Wait State", {
    time: sfn.WaitTime.duration(Duration.seconds(1)),
  });
  return [v, 0];
}
function succeed(stack: Stack): def {
  return [new sfn.Succeed(stack, "Succeed State"), 0];
}
function fail(stack: Stack): def {
  return [new sfn.Fail(stack, "Fail State"), 0];
}
function choice(stack: Stack): def {
  const v = new sfn.Choice(stack, "Choice State")
    .when(sfn.Condition.booleanEquals("$.bool", true), succeed(stack)[0])
    .otherwise(fail(stack)[0]);
  return [v, 0];
}
function choice_fallback(stack: Stack): def {
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
  return [s2.next(choice), 0];
}
function choice_bool(stack: Stack): def {
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
  const v = new sfn.Choice(stack, "Choice State").when(cond5, ok).otherwise(ng);
  return [v, 0];
}
function choice_data_test(stack: Stack): def {
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

  const v = new sfn.Choice(stack, "Choice State")
    .when(
      sfn.Condition.and(
        /**
         * 1. StringEquals, StringEqualsPath
         */
        // stringEquals
        sfn.Condition.stringEquals("$.string", "hello"),
        // stringEqualsPath
        sfn.Condition.stringEqualsJsonPath("$.string", "$.string"),
        sfn.Condition.not(
          sfn.Condition.stringEqualsJsonPath("$.string", "$.largestring")
        ),

        /**
         * 2. StringLessThan, StringLessThanPath
         */
        // stringLessThan
        sfn.Condition.stringLessThan("$.string", "zzzzzzzzz"),
        sfn.Condition.not(sfn.Condition.stringLessThan("$.string", "a")),
        // stringLessThanPath
        sfn.Condition.stringLessThanJsonPath("$.string", "$.largestring"),
        sfn.Condition.not(
          sfn.Condition.stringLessThanJsonPath("$.string", "$.smallstring")
        ),

        /**
         * 3. StringGreaterThan, StringGreaterThanPath
         */
        // stringGreaterThan
        sfn.Condition.stringGreaterThan("$.string", "a"),
        sfn.Condition.not(
          sfn.Condition.stringGreaterThan("$.string", "zzzzzzzzzzzzz")
        ),
        // stringGreaterThanPath
        sfn.Condition.stringGreaterThanJsonPath("$.string", "$.smallstring"),
        sfn.Condition.not(
          sfn.Condition.stringGreaterThanJsonPath("$.string", "$.largestring")
        ),

        /**
         * 4. StringLessThanEquals, StringLessThanEqualsPath
         */
        // stringLessThanEquals
        sfn.Condition.stringLessThanEquals("$.string", "zzzzzzzzz"),
        sfn.Condition.not(sfn.Condition.stringLessThanEquals("$.string", "a")),
        sfn.Condition.stringLessThanEquals("$.string", "hello"),
        // stringLessThanEqualsPath
        sfn.Condition.stringLessThanEqualsJsonPath("$.string", "$.largestring"),
        sfn.Condition.stringLessThanEqualsJsonPath("$.string", "$.string"),
        sfn.Condition.not(
          sfn.Condition.stringLessThanEqualsJsonPath(
            "$.string",
            "$.smallstring"
          )
        ),

        /**
         * 5. StringGreaterThanEquals, StringGreaterThanEqualsPath
         */
        // stringGreaterThanEquals
        sfn.Condition.stringGreaterThanEquals("$.string", "a"),
        sfn.Condition.stringGreaterThanEquals("$.string", "hello"),
        sfn.Condition.not(
          sfn.Condition.stringGreaterThanEquals("$.string", "zzzzzzzzzzzzz")
        ),
        // stringGreaterThanPathEquals
        sfn.Condition.stringGreaterThanEqualsJsonPath(
          "$.string",
          "$.smallstring"
        ),
        sfn.Condition.stringGreaterThanEqualsJsonPath("$.string", "$.string"),
        sfn.Condition.not(
          sfn.Condition.stringGreaterThanEqualsJsonPath(
            "$.string",
            "$.largestring"
          )
        ),

        /**
         * 6. StringMatches
         */
        sfn.Condition.stringMatches("$.string", "hello"),
        sfn.Condition.stringMatches("$.string", "*"),
        sfn.Condition.stringMatches("$.string", "*llo"),
        sfn.Condition.stringMatches("$.string", "hel*"),
        sfn.Condition.stringMatches("$.string", "*h*e*l*l*o*"),
        sfn.Condition.stringMatches("$.wildslash", "a\\*b\\\\c"),
        sfn.Condition.not(sfn.Condition.stringMatches("$.string", "*xxx*")),

        /**
         * 7. NumericEquals, NumericEqualsPath
         */
        // numericEquals
        sfn.Condition.numberEquals("$.number", 3.14),
        // numericEqualsPath
        sfn.Condition.numberEqualsJsonPath("$.number", "$.number"),
        sfn.Condition.not(
          sfn.Condition.numberEqualsJsonPath("$.number", "$.largenumber")
        ),

        /**
         * 8. NumericLessThan, NumericLessThanPath
         */
        // numericLessThan
        sfn.Condition.numberLessThan("$.number", 10000),
        sfn.Condition.not(sfn.Condition.numberLessThan("$.number", 0)),
        // numericLessThanPath
        sfn.Condition.numberLessThanJsonPath("$.number", "$.largenumber"),
        sfn.Condition.not(
          sfn.Condition.numberLessThanJsonPath("$.number", "$.smallnumber")
        ),

        /**
         * 9. NumericGreaterThan, NumericGreaterThanPath
         */
        // numericGreaterThan
        sfn.Condition.numberGreaterThan("$.number", 0),
        sfn.Condition.not(sfn.Condition.numberGreaterThan("$.number", 10000)),
        // numericGreaterThanPath
        sfn.Condition.numberGreaterThanJsonPath("$.number", "$.smallnumber"),
        sfn.Condition.not(
          sfn.Condition.numberGreaterThanJsonPath("$.number", "$.largenumber")
        ),

        /**
         * 10. NumericLessThanEquals, NumericLessThanEqualsPath
         */
        // numericLessThanEquals
        sfn.Condition.numberLessThanEquals("$.number", 10000),
        sfn.Condition.numberLessThanEquals("$.number", 3.14),
        sfn.Condition.not(sfn.Condition.numberLessThanEquals("$.number", 0)),
        // numericLessThanEqualsPath
        sfn.Condition.numberLessThanEqualsJsonPath("$.number", "$.largenumber"),
        sfn.Condition.numberLessThanEqualsJsonPath("$.number", "$.number"),
        sfn.Condition.not(
          sfn.Condition.numberLessThanEqualsJsonPath(
            "$.number",
            "$.smallnumber"
          )
        ),

        /**
         * 11. NumericGreaterThanEquals, NumericGreaterThanEqualsPath
         */
        // numericGreaterThanEquals
        sfn.Condition.numberGreaterThanEquals("$.number", 0),
        sfn.Condition.numberGreaterThanEquals("$.number", 3.14),
        sfn.Condition.not(
          sfn.Condition.numberGreaterThanEquals("$.number", 10000)
        ),
        // numericGreaterThanEqualsPath
        sfn.Condition.numberGreaterThanEqualsJsonPath(
          "$.number",
          "$.smallnumber"
        ),
        sfn.Condition.numberGreaterThanEqualsJsonPath("$.number", "$.number"),
        sfn.Condition.not(
          sfn.Condition.numberGreaterThanEqualsJsonPath(
            "$.number",
            "$.largenumber"
          )
        ),

        /**
         * 12. BooleanEquals, BooleanEqualsPath
         */
        // booleanEquals
        sfn.Condition.booleanEquals("$.bool", true),
        // booleanEqualsPath
        sfn.Condition.not(
          sfn.Condition.booleanEqualsJsonPath("$.bool", "$.object.bool")
        ),

        /**
         * 13. TimestampEquals, TimestampEqualsPath
         */
        // timestampEquals
        sfn.Condition.timestampEquals("$.timestamp", "2016-03-14T01:59:00Z"),
        // timestampEqualsPath
        sfn.Condition.timestampEqualsJsonPath("$.timestamp", "$.timestamp"),
        sfn.Condition.not(
          sfn.Condition.timestampEqualsJsonPath(
            "$.timestamp",
            "$.largetimestamp"
          )
        ),

        /**
         * 14. TimestampLessThan, TimestampLessThanPath
         */
        // timestampLessThan
        sfn.Condition.timestampLessThan("$.timestamp", "2030-01-23T01:23:00Z"),
        sfn.Condition.not(
          sfn.Condition.timestampLessThan("$.timestamp", "1999-11-11T11:11:11Z")
        ),
        // timestampLessThanPath
        sfn.Condition.timestampLessThanJsonPath(
          "$.timestamp",
          "$.largetimestamp"
        ),
        sfn.Condition.not(
          sfn.Condition.timestampLessThanJsonPath(
            "$.timestamp",
            "$.smalltimestamp"
          )
        ),

        /**
         * 15. TimestampGreaterThan, TimestampGreaterThanPath
         */
        // timestampGreaterThan
        sfn.Condition.timestampGreaterThan(
          "$.timestamp",
          "1999-11-11T11:11:11Z"
        ),
        sfn.Condition.not(
          sfn.Condition.timestampGreaterThan(
            "$.timestamp",
            "2030-01-23T01:23:00Z"
          )
        ),
        // timestampGreaterThanPath
        sfn.Condition.timestampGreaterThanJsonPath(
          "$.timestamp",
          "$.smalltimestamp"
        ),
        sfn.Condition.not(
          sfn.Condition.timestampGreaterThanJsonPath(
            "$.timestamp",
            "$.largetimestamp"
          )
        ),

        /**
         * 16. TimestampLessThanEquals, TimestampLessThanEqualsPath
         */
        // timestampLessThanEquals
        sfn.Condition.timestampLessThanEquals(
          "$.timestamp",
          "2030-01-23T01:23:00Z"
        ),
        sfn.Condition.timestampLessThanEquals(
          "$.timestamp",
          "2016-03-14T01:59:00Z"
        ),
        sfn.Condition.not(
          sfn.Condition.timestampLessThanEquals(
            "$.timestamp",
            "1999-11-11T11:11:11Z"
          )
        ),
        // timestampLessThanEqualsPath
        sfn.Condition.timestampLessThanEqualsJsonPath(
          "$.timestamp",
          "$.largetimestamp"
        ),
        sfn.Condition.timestampLessThanEqualsJsonPath(
          "$.timestamp",
          "$.timestamp"
        ),
        sfn.Condition.not(
          sfn.Condition.timestampLessThanEqualsJsonPath(
            "$.timestamp",
            "$.smalltimestamp"
          )
        ),

        /**
         * 17. TimestampGreaterThanEquals, TimestampGreaterThanEqualsPath
         */
        // timestampGreaterThanEquals
        sfn.Condition.timestampGreaterThanEquals(
          "$.timestamp",
          "1999-11-11T11:11:11Z"
        ),
        sfn.Condition.timestampGreaterThanEquals(
          "$.timestamp",
          "2016-03-14T01:59:00Z"
        ),
        sfn.Condition.not(
          sfn.Condition.timestampGreaterThanEquals(
            "$.timestamp",
            "2030-01-23T01:23:00Z"
          )
        ),
        // timestampGreaterThanPath
        sfn.Condition.timestampGreaterThanEqualsJsonPath(
          "$.timestamp",
          "$.smalltimestamp"
        ),
        sfn.Condition.timestampGreaterThanEqualsJsonPath(
          "$.timestamp",
          "$.timestamp"
        ),
        sfn.Condition.not(
          sfn.Condition.timestampGreaterThanEqualsJsonPath(
            "$.timestamp",
            "$.largetimestamp"
          )
        ),

        /**
         * 18. IsNull
         */
        sfn.Condition.isNull("$.null"),
        sfn.Condition.not(sfn.Condition.isNull("$.string")),

        /**
         * 19. IsPresent
         */
        sfn.Condition.isPresent("$.bool"),
        sfn.Condition.not(sfn.Condition.isPresent("$.non.existing.path")),

        /**
         * 20. IsNumeric
         */
        sfn.Condition.isNumeric("$.int"),
        sfn.Condition.not(sfn.Condition.isNumeric("$.string")),

        /**
         * 21. IsString
         */
        sfn.Condition.isString("$.string"),
        sfn.Condition.not(sfn.Condition.isString("$.bool")),

        /**
         * 22. IsBoolean
         */
        sfn.Condition.isBoolean("$.bool"),
        sfn.Condition.not(sfn.Condition.isBoolean("$.string")),

        /**
         * 23. IsTimestamp
         */
        sfn.Condition.isTimestamp("$.timestamp"),
        sfn.Condition.not(sfn.Condition.isTimestamp("$.bool")),
        sfn.Condition.not(sfn.Condition.isTimestamp("$.string"))
      ),
      ok
    )
    .otherwise(ng);
  return [v, 0];
}
function task(stack: Stack): def {
  const v = new ScriptTask(stack, "Task State", {
    scriptPath: "_workflow/script/script1.sh",
  });
  return [v, 0];
}
function task_filter(stack: Stack): def {
  const v = new ScriptTask(stack, "Task State", {
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
  return [v, 0];
}
function task_retry(stack: Stack): def {
  const task = new ScriptTask(stack, "Task State", {
    scriptPath: "_workflow/script/script2.sh",
    resultPath: "$.args",
  });
  const chain = new sfn.Parallel(stack, "Chain").branch(task);
  chain.addRetry({
    maxAttempts: 10,
    backoffRate: 1,
    interval: Duration.seconds(0),
  });
  return [chain, 0];
}
function task_catch(stack: Stack): def {
  const p1 = new sfn.Pass(stack, "Pass State1");
  const task = new ScriptTask(stack, "Task State", {
    scriptPath: "::", // invalid resource path
  });
  task.addCatch(p1, {
    errors: ["States.ALL"],
  });
  return [task, 0];
}
function task_ctxobj(stack: Stack): def {
  const v = new ScriptTask(stack, "Task State", {
    scriptPath: "_workflow/script/script1.sh",
    resultSelector: {
      ctx: {
        "ctx_aaa.$": "$$.aaa",
      },
    },
  });
  return [v, 0];
}
function parallel(stack: Stack): def {
  const v = new sfn.Parallel(stack, "Parallel State")
    .branch(pass(stack)[0])
    .branch(succeed(stack)[0]);
  return [v, 0];
}
function parallel_ctxobj(stack: Stack): def {
  const s1 = new sfn.Pass(stack, "S1", {
    parameters: {
      "ctx_aaa.$": "$$.aaa",
    },
    resultPath: "$.result.ctx",
    outputPath: "$.result",
  });
  const s2 = new sfn.Pass(stack, "S2", {
    parameters: {
      "ctx_aaa.$": "$$.aaa",
    },
    resultPath: "$.result.ctx",
    outputPath: "$.result",
  });
  const s3 = new sfn.Pass(stack, "S3", {
    parameters: {
      "ctx_aaa.$": "$$.aaa",
    },
    resultPath: "$.result.ctx",
    outputPath: "$.result",
  });

  const v = new sfn.Parallel(stack, "Parallel State")
    .branch(s1)
    .branch(s2)
    .branch(s3);
  return [v, 0];
}
function map(stack: Stack): def {
  const map = new sfn.Map(stack, "Map State", {
    maxConcurrency: 1,
    itemsPath: sfn.JsonPath.stringAt("$.inputForMap"),
  });
  map.iterator(new sfn.Pass(stack, "Pass State"));
  return [map, 0];
}
function map_concurrency(stack: Stack): def {
  const map = new sfn.Map(stack, "Map State", {
    maxConcurrency: 3,
    itemsPath: sfn.JsonPath.stringAt("$.inputForMap"),
  });
  map.iterator(new sfn.Pass(stack, "Pass State"));
  return [map, 0];
}
function map_ctxobj(stack: Stack): def {
  const map = new sfn.Map(stack, "Map State", {
    maxConcurrency: 3,
    itemsPath: sfn.JsonPath.stringAt("$.inputForMap"),
  });
  map.iterator(
    new sfn.Pass(stack, "Pass State", {
      parameters: {
        "index.$": "$$.Map.Item.Index",
        "value.$": "$$.Map.Item.Value",
      },
    })
  );
  return [map, 0];
}
function map_ctxobj2(stack: Stack): def {
  const map = new sfn.Map(stack, "Map State", {
    maxConcurrency: 3,
    itemsPath: sfn.JsonPath.stringAt("$.inputForMap"),
  });
  map.iterator(
    new sfn.Pass(stack, "Pass State", {
      parameters: {
        "ctx_aaa.$": "$$.aaa",
      },
    })
  );
  return [map, 0];
}

const workflows: ((stack: Stack) => def)[] = [
  pass,
  pass_chain,
  pass_result,
  pass_parameters,
  pass_intrinsic,
  pass_ctxobj,
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
  task_ctxobj,
  parallel,
  parallel_ctxobj,
  map,
  map_concurrency,
  map_ctxobj,
  map_ctxobj2,
];

function list() {
  workflows.forEach(function (elm) {
    console.log(elm.name);
  });
}

function render(sm: def) {
  const graph = new sfn.StateGraph(sm[0].startState, "Graph");
  graph.timeout = Duration.seconds(sm[1]);
  return new Stack().resolve(graph.toGraphJson());
}

function toString(sm: def): string {
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
