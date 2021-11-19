package statemachine

import (
	"errors"
	"fmt"
	"os"
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
				target: newASLError("sample error"),
			},
			want: true,
		},
		{
			name: "unmatch",
			fields: fields{
				name: "sample error",
			},
			args: args{
				target: newASLError("AAA BBB CCC"),
			},
			want: false,
		},
		{
			name: "States.All",
			fields: fields{
				name: statesAll,
			},
			args: args{
				target: newASLError("AAA BBB CCC"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := aslError{
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

func Test_aslError_Unwrap(t *testing.T) {
	type fields struct {
		name  string
		cause error
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "with cause",
			fields: fields{
				name:  "internal error",
				cause: os.ErrInvalid,
			},
			want: os.ErrInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newASLErrorWithCause(
				tt.fields.name,
				tt.fields.cause,
			)
			if err := e.Unwrap(); err != tt.want {
				t.Errorf("aslError.Unwrap() got = %v, want %v", err, tt.want)
			}
		})
	}
}
