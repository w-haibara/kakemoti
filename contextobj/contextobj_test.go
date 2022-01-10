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
			got := Get(ctx)
			if got == nil {
				t.Error("Get(ctx) returns nil")
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
