package statemachine

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"karage/log"
)

func TestNewStateMachine(t *testing.T) {
	type Want struct {
		States         []string
		Comment        string
		StartAt        string
		Version        string
		TimeoutSeconds int64
	}
	tests := []struct {
		name    string
		asl     string
		want    Want
		wantErr bool
	}{
		{
			name: "top-lebel-1",
			asl: `{
				"Version": "1.0",
				"Comment": "sample comment",
				"StartAt": "a1",
				"TimeoutSeconds": 0,
				"States": {
					"a1": {
						"Type": "Succeed"
					}
				}
			}`,
			want: Want{
				States:         []string{"a1"},
				Comment:        "sample comment",
				StartAt:        "a1",
				Version:        "1.0",
				TimeoutSeconds: 0,
			},
			wantErr: false,
		},
		{
			name: "top-lebel-without-version",
			asl: `{
				"Comment": "sample comment",
				"StartAt": "a1",
				"TimeoutSeconds": 0,
				"States": {
					"a1": {
						"Type": "Succeed"
					}
				}
			}`,
			want: Want{
				States:         []string{"a1"},
				Comment:        "sample comment",
				StartAt:        "a1",
				Version:        "1.0",
				TimeoutSeconds: 0,
			},
			wantErr: false,
		},
		{
			name: "top-lebel-without-timeoutseconds",
			asl: `{
				"Comment": "sample comment",
				"StartAt": "a1",
				"States": {
					"a1": {
						"Type": "Succeed"
					}
				}
			}`,
			want: Want{
				States:         []string{"a1"},
				Comment:        "sample comment",
				StartAt:        "a1",
				Version:        "1.0",
				TimeoutSeconds: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		logger := log.NewLogger()
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStateMachine(bytes.NewBufferString(tt.asl), logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStateMachine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, v := range tt.want.States {
				_, ok := got.States[v]
				if !ok {
					t.Errorf("States: invalid key: %s", v)
					return
				}
			}

			if *got.Comment != tt.want.Comment {
				t.Errorf("Comment: %s (got) != %s (want)", *got.Comment, tt.want.Comment)
				return
			}

			if *got.StartAt != tt.want.StartAt {
				t.Errorf("StartAt: %s (got) != %s (want)", *got.StartAt, tt.want.StartAt)
				return
			}

			if *got.Version != tt.want.Version {
				t.Errorf("Version: %s (got) != %s (want)", *got.Version, tt.want.Version)
				return
			}

			if got.TimeoutSeconds == nil {
				t.Errorf("TimeoutSeconds is nil")
				return
			}

			if *got.TimeoutSeconds != tt.want.TimeoutSeconds {
				t.Errorf("TimeoutSeconds: %d (got) != %d (want)", *got.TimeoutSeconds, tt.want.TimeoutSeconds)
				return
			}
		})
	}
}

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
			name: "pass-end",
			asl: `{
				"StartAt": "pass1",
				"States": {
					"pass1": {
						"Type": "Pass",
						"end": true
					}
				}
			}`,
			input:   "{}",
			want:    "{}",
			wantErr: false,
		},
		{
			name: "pass-next",
			asl: `{
				"StartAt": "pass1",
				"States": {
					"pass1": {
						"Type": "Pass",
						"Next": "pass2"
					},
					"pass2": {
						"Type": "Pass",
						"End": true
					}
				}
			}`,
			input: `{
				"abc":"123"
			}`,
			want: `{
				"abc":"123"
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
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Start() returns an error: %v", err)
				}
				return
			}

			if string(got) == tt.want {
				return
			}

			g := new(interface{})
			if err := json.Unmarshal(got, g); err != nil {
				t.Errorf("invalid json format: %v \n%q", err, got)
				return
			}

			if reflect.DeepEqual(g, want) {
				t.Errorf("Start() = %q, want %q", got, tt.want)
			}
		})
	}
}
