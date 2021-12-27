package compiler

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/k0kubun/pp"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		name       string
		asl        string
		wantStates States
		wantErr    bool
	}{
		{
			"basic",
			`{
	"StartAt": "Pass State",
	"States": {
	  "Pass State": {
		"Type": "Pass",
		"End": true
	  }
	}
}`,
			States{
				State{"Pass", "Pass State", "",
					&PassState{
						CommonState4: CommonState4{
							CommonState3: CommonState3{
								End: true,
								CommonState2: CommonState2{
									CommonState1: CommonState1{
										Type: "Pass",
									},
								},
							},
						}},
					nil,
				},
			},
			false,
		},
		{
			"choice(fallback)",
			`{
				"StartAt": "State1",
				"States": {
				  "State1": {
					"Type": "Pass",
					"Next": "State2"
				  },
				  "Choice State": {
					"Type": "Choice",
					"Choices": [
					  {
						"Variable": "$.bool",
						"BooleanEquals": false,
						"Next": "State1"
					  }
					],
					"Default": "State2"
				  },
				  "State2": {
					"Type": "Pass",
					"Next": "Choice State"
				  }
				}
			  }`,
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compile(context.TODO(), bytes.NewBufferString(tt.asl))
			if (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotStates := got.States
			if !reflect.DeepEqual(gotStates, tt.wantStates) {
				t.Errorf("Compile() = \n%#v\n want = \n%#v\n", gotStates, tt.wantStates)
				t.Errorf("====== diff ======\n%s\n", diff.LineDiff(pp.Sprint(gotStates), pp.Sprint(tt.wantStates)))
			}
		})
	}
}
