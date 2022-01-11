package worker

import (
	"context"
	"reflect"
	"testing"
)

func Test_parseIntrinsicFunction(t *testing.T) {
	tests := []struct {
		name     string
		fnstr    string
		input    interface{}
		wantName string
		wantArgs []interface{}
	}{
		{"basic1", "a('x')", nil, "a", []interface{}{"x"}},
		{"basic2", "a(1)", nil, "a", []interface{}{1}},
		{"basic3", "a(3.14)", nil, "a", []interface{}{3.14}},
		{"basic4", "a('x', 1, 3.14)", nil, "a", []interface{}{"x", 1, 3.14}},
		//{"", "States.Format('Hello, my name is {}.', $.name)", 0, "", []interface{}{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotArgs, err := parseIntrinsicFunction(context.Background(), tt.fnstr, tt.input)
			if err != nil {
				t.Errorf("parseIntrinsicFunction() error = %v", err)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("parseIntrinsicFunction() got = %#v, want %#v", gotName, tt.wantName)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("parseIntrinsicFunction() got = %#v, want %#v", gotArgs, tt.wantArgs)
			}
		})
	}
}
