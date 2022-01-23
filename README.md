# kakemoti
kakemoti is a simple tool that executes workflows defined in the [Amazon States Language](https://states-language.net/). It breaks up large scripts into pieces to improve readability, observability and serviceability.

# TODO
- [x] Top-level fields
  - [x] States
  - [x] StartAt
  - [x] Comment
  - [x] Version
  - [x] TimeoutSeconds
- [x] States
  - [x] Pass State
  - [x] Task State
  - [x] Choice State
    - [x] Boolean expression
    - [x] Data-test expression
      - [x] StringEquals
      - [x] StringEqualsPath
      - [x] StringLessThan
      - [x] StringLessThanPath
      - [x] StringGreaterThan
      - [x] StringGreaterThanPath
      - [x] StringLessThanEquals
      - [x] StringLessThanEqualsPath
      - [x] StringGreaterThanEquals
      - [x] StringGreaterThanEqualsPath
      - [x] StringMatches
      - [x] NumericEquals
      - [x] NumericEqualsPath
      - [x] NumericLessThan
      - [x] NumericLessThanPath
      - [x] NumericGreaterThan
      - [x] NumericGreaterThanPath
      - [x] NumericLessThanEquals
      - [x] NumericLessThanEqualsPath
      - [x] NumericGreaterThanEquals
      - [x] NumericGreaterThanEqualsPath
      - [x] BooleanEquals
      - [x] BooleanEqualsPath
      - [x] TimestampEquals
      - [x] TimestampEqualsPath
      - [x] TimestampLessThan
      - [x] TimestampLessThanPath
      - [x] TimestampGreaterThan
      - [x] TimestampGreaterThanPath
      - [x] TimestampLessThanEquals
      - [x] TimestampLessThanEqualsPath
      - [x] TimestampGreaterThanEquals
      - [x] TimestampGreaterThanEqualsPath
      - [x] IsNull
      - [x] IsPresent
      - [x] IsNumeric
      - [x] IsString
      - [x] IsBoolean
      - [x] IsTimestamp
  - [x] Wait State
  - [x] Succeed State
  - [x] Fail State
  - [x] Parallel State
  - [x] Map State
    - [x] Map State input/output processing
    - [x] Map State concurrency
    - [x] Map State Iterator definition
- [x] Transitions
- [x] Timestamps
- [x] Data
- [x] The Context Object
- [x] Paths
- [x] Reference Paths
- [x] Payload Template
- [x] Intrinsic Functions
  - [x] States.Format
  - [x] States.StringToJson
  - [x] States.JsonToString
  - [x] States.Array
- [x] Input and Output Processing
  - [x] InputPath
  - [x] Parameters
  - [x] ResultSelector
  - [x] ResultPath
  - [x] OutputPath
- [ ] Errors
  - [x] States.ALL
  - [ ] States.HeartbeatTimeout
  - [ ] States.Timeout
  - [x] States.TaskFailed
  - [x] States.Permissions
  - [x] States.ResultPathMatchFailure
  - [x] States.ParameterPathFailure
  - [x] States.BranchFailed
  - [ ] States.NoChoiceMatched
  - [ ] States.IntrinsicFailure
