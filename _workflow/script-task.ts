import { Construct } from "constructs";
import { aws_stepfunctions as sfn } from "aws-cdk-lib";

export interface ScriptTaskProps extends sfn.TaskStateBaseProps {
  readonly parameters?: sfn.TaskInput | undefined;
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
      Parameters: this.props.parameters
        ? this.props.parameters.value
        : undefined,
    };
  }
}
