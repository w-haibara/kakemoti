package fn

import (
	"context"
	"encoding/json"
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

func TestDoStatesJsonToString(t *testing.T) {
	tests := []struct {
		name    string
		args    []interface{}
		want    interface{}
		wantErr bool
	}{
		{"basic1", []interface{}{map[string]interface{}{"aaa": 111}}, `{"aaa":111}`, false},
		{"basic2", []interface{}{map[string]interface{}{"aaa": 111, "bbb": map[string]interface{}{"ccc": "xxx"}}}, `{"aaa":111, "bbb":{"ccc": "xxx"}}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoStatesJsonToString(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoStatesJsonToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if d := jsonDiff(t, got, tt.want); d != "" {
				t.Errorf("DoStatesStringToJson() failed: \n%s", d)
			}
		})
	}
}

func jsonDiff(t *testing.T, arg1, arg2 interface{}) string {
	s1, ok := arg1.(string)
	if !ok {
		t.Fatal("arg1.(string) failed")
	}
	s2, ok := arg2.(string)
	if !ok {
		t.Fatal("arg2.(string) failed")
	}

	var v1, v2 interface{}
	if err := json.Unmarshal([]byte(s1), &v1); err != nil {
		t.Fatal("json.Unmarshal(b1, &v1) failed:", err)
	}
	if err := json.Unmarshal([]byte(s2), &v2); err != nil {
		t.Fatal("json.Unmarshal(b2, &v2) failed:", err)
	}
	return cmp.Diff(v1, v2)
}
