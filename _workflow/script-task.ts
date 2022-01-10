import { Construct } from "constructs";
import { aws_stepfunctions as sfn } from "aws-cdk-lib";

export interface ScriptTaskProps extends sfn.TaskStateBaseProps {
  readonly payload?: sfn.TaskInput;
  readonly scriptPath: string;
}

export class ScriptTask extends sfn.TaskStateBase {
  protected readonly taskMetrics?: undefined;
  protected readonly taskPolicies?: undefined;

  constructor(
    scope: Construct,
    id: string,
    private readonly props: ScriptTaskProps
  ) {
    super(scope, id, props);
  }

  protected _renderTask(): any {
    return {
      Resource: "script:" + this.props.scriptPath,
      Parameters: sfn.FieldUtils.renderObject({
        Payload: this.props.payload
          ? this.props.payload.value
          : sfn.TaskInput.fromJsonPathAt("$").value,
      }),
    };
  }
}
