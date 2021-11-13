package statemachine

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"karage/log"
)

func TestStart(t *testing.T) {
	const (
		timeout = 10
	)

	tests := []struct {
		name    string
		asl     string
		input   string
		want    string
		wantErr bool
	}{
		{
			name: "minimal",
			asl: `{
				"Comment": "A simple minimal example of the States language",
				"StartAt": "Hello World",
				"States": {
				"Hello World": {
				  "Type": "Task",
				  "Resource": "script:../workflows/task-script1/task.sh",
				  "End": true
				}
			  }
			}`,
			input: `{
				"args": ["1", "2", "3"]
			}`,
			want: `{
				"result": "args: 1, 2, 3"
			}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ctx := context.TODO()
		logger := log.NewLogger()

		want := new(interface{})
		if err := json.Unmarshal([]byte(tt.want), want); err != nil {
			t.Errorf("[tt.want] is invalid json format: %v\n%q", err, tt.want)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := Start(ctx,
				bytes.NewBufferString(tt.asl),
				bytes.NewBufferString(tt.input),
				timeout, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			g := new(interface{})
			if err := json.Unmarshal(got, g); err != nil {
				t.Errorf("invalid json format: %v \n%q", err, got)
				return
			}

			if reflect.DeepEqual(g, want) {
				t.Errorf("Start() = %v, want %v", g, tt.want)
			}
		})
	}
}
