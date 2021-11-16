package statemachine

import (
	"errors"
	"fmt"
	"testing"
)

func Test_statemachineError_Is(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "match",
			fields: fields{
				name: "sample error",
			},
			args: args{
				target: newStateMachineError("sample error"),
			},
			want: true,
		},
		{
			name: "unmatch",
			fields: fields{
				name: "sample error",
			},
			args: args{
				target: newStateMachineError("AAA BBB CCC"),
			},
			want: false,
		},
		{
			name: "States.All",
			fields: fields{
				name: statesAll,
			},
			args: args{
				target: newStateMachineError("AAA BBB CCC"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := statemachineError{
				name: tt.fields.name,
			}
			if got := errors.Is(tt.args.target, e); got != tt.want {
				t.Errorf("errors.Is(tt.args.target, e) = %v, want %v", got, tt.want)
			}

			if got := errors.Is(fmt.Errorf("wrap:%w", tt.args.target), e); got != tt.want {
				t.Errorf(`errors.Is(fmt.Errorf("wrap:%%w", tt.args.target) = %v, want %v`, got, tt.want)
			}
		})
	}
}
