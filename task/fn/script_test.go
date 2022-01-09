package fn

import (
	"reflect"
	"testing"

	"github.com/k0kubun/pp"
)

func Test_marshalArgsToMap(t *testing.T) {
	type args struct {
		key  string
		args interface{}
	}
	tests := []struct {
		name string
		args args
		want Obj
	}{
		{
			"string",
			args{scriptInputPrefix, "aaa"},
			Obj{scriptInputPrefix: "aaa"},
		},
		{
			"int",
			args{scriptInputPrefix, 123},
			Obj{scriptInputPrefix: "123"},
		},
		{
			"string slice",
			args{scriptInputPrefix, []string{"111", "222", "333"}},
			Obj{
				scriptInputPrefix + "_0": "111",
				scriptInputPrefix + "_1": "222",
				scriptInputPrefix + "_2": "333",
			},
		},
		{
			"float64 slice",
			args{scriptInputPrefix, []float64{1, 2.0, 3.1}},
			Obj{
				scriptInputPrefix + "_0": "1",
				scriptInputPrefix + "_1": "2",
				scriptInputPrefix + "_2": "3.1",
			},
		},
		{
			"float32 slice",
			args{scriptInputPrefix, []float32{1, 2.0, 3.1}},
			Obj{
				scriptInputPrefix + "_0": "1",
				scriptInputPrefix + "_1": "2",
				scriptInputPrefix + "_2": "3.1",
			},
		},
		{
			"int slice",
			args{scriptInputPrefix, []int{1, 2, 3}},
			Obj{
				scriptInputPrefix + "_0": "1",
				scriptInputPrefix + "_1": "2",
				scriptInputPrefix + "_2": "3",
			},
		},
		{
			"int32 slice",
			args{scriptInputPrefix, []int32{1, 2, 3}},
			Obj{
				scriptInputPrefix + "_0": "1",
				scriptInputPrefix + "_1": "2",
				scriptInputPrefix + "_2": "3",
			},
		},
		{
			"Obj",
			args{scriptInputPrefix, Obj{
				"aaa": "111",
				"bbb": "222",
				"ccc": 333,
				"ddd": 444,
			}},
			Obj{
				scriptInputPrefix + "_aaa": "111",
				scriptInputPrefix + "_bbb": "222",
				scriptInputPrefix + "_ccc": "333",
				scriptInputPrefix + "_ddd": "444",
			},
		},
		{
			"Obj slice",
			args{scriptInputPrefix, []Obj{
				{"aaa": "111", "123": 456},
				{"bbb": "222", "234": 567},
			}},
			Obj{
				scriptInputPrefix + "_0_aaa": "111",
				scriptInputPrefix + "_0_123": "456",
				scriptInputPrefix + "_1_bbb": "222",
				scriptInputPrefix + "_1_234": "567",
			},
		},
		{
			"Nested Obj",
			args{scriptInputPrefix, Obj{
				"aaa": Obj{
					"bbb": "222",
					"ccc": "333",
				},
			}},
			Obj{
				scriptInputPrefix + "_aaa_bbb": "222",
				scriptInputPrefix + "_aaa_ccc": "333",
			},
		},
		{
			"Nested Obj slice",
			args{scriptInputPrefix, Obj{
				"aaa": []Obj{
					{"bbb": "222", "ccc": "333"},
					{"bbb": "222", "ccc": "333"},
				},
			}},
			Obj{
				scriptInputPrefix + "_aaa_0_bbb": "222",
				scriptInputPrefix + "_aaa_0_ccc": "333",
				scriptInputPrefix + "_aaa_1_bbb": "222",
				scriptInputPrefix + "_aaa_1_ccc": "333",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := marshalArgsToMap(tt.args.key, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("marshalArgsToMap() = %v, want %v", got, tt.want)
				t.Errorf(pp.Sprintln(got))
			}
		})
	}
}
