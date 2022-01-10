import * as cdk from "@aws-cdk/core";
import * as core from "@aws-cdk/core";
import * as sfn from "@aws-cdk/aws-stepfunctions";

export interface ScriptTaskProps {
  readonly comment?: string;
  readonly inputPath?: string;
  readonly outputPath?: string;
  readonly resultPath?: string;
  readonly resultSelector?: { [name: string]: any };
  readonly parameters?: { [name: string]: any };
  readonly scriptPath: string;
  readonly timeout?: cdk.Duration;
  readonly heartbeat?: cdk.Duration;
}

export class ScriptTask extends sfn.State implements sfn.INextable {
  public readonly endStates: sfn.INextable[];
  private readonly taskProps: ScriptTaskProps;

  constructor(scope: core.Construct, id: string, props: ScriptTaskProps) {
    super(scope, id, props);
    this.taskProps = props;
    this.endStates = [this];
  }

  public addRetry(props: sfn.RetryProps = {}): ScriptTask {
    super._addRetry(props);
    return this;
  }

  public addCatch(handler: sfn.IChainable, props: sfn.CatchProps = {}): ScriptTask {
    super._addCatch(handler.startState, props);
    return this;
  }

  public next(next: sfn.IChainable): sfn.Chain {
    super.makeNext(next.startState);
    return sfn.Chain.sequence(this, next);
  }

  public toStateJson(): object {
    return {
      ...this.renderNextEnd(),
      ...this.renderRetryCatch(),
      ...this.renderInputOutput(),
      Type: "Task",
      Comment: this.taskProps.comment,
      Resource: "script:"+this.taskProps.scriptPath,
      Parameters:
        this.taskProps.parameters &&
        sfn.FieldUtils.renderObject(this.taskProps.parameters),
      ResultPath: sfn.renderJsonPath(this.resultPath),
      TimeoutSeconds:
        this.taskProps.timeout && this.taskProps.timeout.toSeconds(),
      HeartbeatSeconds:
        this.taskProps.heartbeat && this.taskProps.heartbeat.toSeconds(),
    };
  }

  private renderParameters(): any {
    return sfn.FieldUtils.renderObject({
      Parameters: this.parameters,
    });
  }
}