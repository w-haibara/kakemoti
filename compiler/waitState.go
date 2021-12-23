package compiler

type WaitState struct {
	CommonState3
	Seconds       *int64  `json:"Seconds"`
	Timestamp     *string `json:"Timestamp"`
	SecondsPath   *string `json:"SecondsPath"`
	TimestampPath *string `json:"TimestampPath"`
}
