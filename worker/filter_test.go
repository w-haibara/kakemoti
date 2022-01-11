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
		{"basic4", "abc_XYZ.123('x', 1, 3.14)", nil, "abc_XYZ.123", []interface{}{"x", 1, 3.14}},
		{"path1", "a($.aaa)", map[string]interface{}{"aaa": 111}, "a", []interface{}{111}},
		{"path1", "a(1, $.aaa)", map[string]interface{}{"aaa": 111}, "a", []interface{}{1, 111}},
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
