import * as cdk from "@aws-cdk/core";
import * as core from "@aws-cdk/core";
import * as sfn from "@aws-cdk/aws-stepfunctions";

export interface TaskProps {
  readonly comment?: string;
  readonly inputPath?: string;
  readonly outputPath?: string;
  readonly resultPath?: string;
  readonly parameters?: { [name: string]: any };
  readonly resource?: string;
  readonly timeout?: cdk.Duration;
  readonly heartbeat?: cdk.Duration;
}

export class Task extends sfn.State implements sfn.INextable {
  public readonly endStates: sfn.INextable[];
  private readonly taskProps: TaskProps;


  constructor(scope: core.Construct, id: string, props: TaskProps = {}) {
    super(scope, id, props);
    if (props.resource === undefined) {
      throw new Error('Task Resource is not found');
    }
    this.taskProps = props;
    this.endStates = [this];
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
      Resource: this.taskProps.resource,
      Parameters:
        this.taskProps.parameters &&
        sfn.FieldUtils.renderObject(this.taskProps.parameters),
      ResultPath: sfn.renderJsonPath(this.resultPath),
      TimeoutSeconds: this.taskProps.timeout && this.taskProps.timeout.toSeconds(),
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
