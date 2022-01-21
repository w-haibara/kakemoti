package contextobj

import (
	"context"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	bg := context.Background()
	tests := []struct {
		name string
		want map[string]interface{}
	}{
		{"int", map[string]interface{}{"aaa": 111}},
		{"string", map[string]interface{}{"aaa": "bbb"}},
		{"slice", map[string]interface{}{"aaa": []string{"x", "y"}}},
		{"struct", map[string]interface{}{"aaa": struct {
			a int
			b string
		}{
			99,
			"xx",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := New(bg)
			for k, v := range tt.want {
				ctx = Set(ctx, k, v)
			}
			got := GetAll(ctx)
			if got == nil {
				t.Error("GetAll(ctx) returns nil")
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll(ctx) failed: got %#v, want %#v", got, tt.want)
			}

			for k, v := range tt.want {
				ctx = Set(ctx, k, v)
			}
			got = make(map[string]interface{})
			for k := range tt.want {
				v, ok := Get(ctx, k)
				if !ok {
					t.Error("Get(ctx, k) returns not ok")
				}
				got[k] = v
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get(ctx, k) failed: got %#v, want %#v", got, tt.want)
			}
		})
	}
}
