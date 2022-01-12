package fn

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDoStatesFormat(t *testing.T) {
	tests := []struct {
		name    string
		args    []interface{}
		want    interface{}
		wantErr bool
	}{
		{"string", []interface{}{"aaa={}", "bbb"}, "aaa=bbb", false},
		{"escaped", []interface{}{"aaa=\\{}{}\\{}", "bbb"}, "aaa=\\{}bbb\\{}", false},
		{"int", []interface{}{"aaa={}", 111}, "aaa=111", false},
		{"float", []interface{}{"aaa={}", 3.14}, "aaa=3.14", false},
		{"all", []interface{}{"{}, {}, {}", "bbb", 111, 3.14}, "bbb, 111, 3.14", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoStatesFormat(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoStatesFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoStatesFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoStatesStringToJson(t *testing.T) {
	tests := []struct {
		name    string
		args    []interface{}
		want    interface{}
		wantErr bool
	}{
		{"basic1", []interface{}{`{"aaa":111}`}, map[string]interface{}{"aaa": 111.0}, false},
		{"basic2", []interface{}{`{"aaa":111, "bbb":{"ccc": "xxx"}}`}, map[string]interface{}{"aaa": 111.0, "bbb": map[string]interface{}{"ccc": "xxx"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoStatesStringToJson(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoStatesStringToJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("DoStatesStringToJson() failed: \n%s", d)
			}
		})
	}
}
