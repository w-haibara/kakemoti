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
		{"quoted string", "a('x')", nil, "a", []interface{}{"x"}},
		{"int", "a(1)", nil, "a", []interface{}{1}},
		{"float", "a(3.14)", nil, "a", []interface{}{3.14}},
		{"null", "a(null)", nil, "a", []interface{}{"null"}},
		{"path", "a($.aaa)", map[string]interface{}{"aaa": 111}, "a", []interface{}{111}},
		{"all", "abc_XYZ.123('x', 1, 3.14, null, $.aaa)", map[string]interface{}{"aaa": 111}, "abc_XYZ.123", []interface{}{"x", 1, 3.14, "null", 111}},
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
