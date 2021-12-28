package compiler

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		name       string
		asl        string
		wantStates []States
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
			[]States{{State{"Pass", "Pass State", "",
				&PassState{
					CommonState4: CommonState4{
						CommonState3: CommonState3{
							End: true,
							CommonState2: CommonState2{
								CommonState1: CommonState1{
									Type: "Pass",
								}}}}}}}},
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
					"Default": "State3"
				  },
				  "State2": {
					"Type": "Pass",
					"Next": "Choice State"
				  },
				  "State3": {
					"Type": "Pass",
					"End": true
				  }
				}
			  }`,
			[]States{
				{
					State{"Pass", "State1", "State2",
						&PassState{
							CommonState4: CommonState4{
								CommonState3: CommonState3{
									Next: "State2",
									CommonState2: CommonState2{
										CommonState1: CommonState1{
											Type: "Pass",
										}}}}}},
					State{"Pass", "State2", "Choice State",
						&PassState{
							CommonState4: CommonState4{
								CommonState3: CommonState3{
									Next: "Choice State",
									CommonState2: CommonState2{
										CommonState1: CommonState1{
											Type: "Pass",
										}}}}}},
					State{"Choice", "Choice State", "",
						&ChoiceState{
							Choices: []Choice{{
								Rule: Rule{
									Variable1: "$.bool",
									Variable2: false,
									Operator:  "BooleanEquals",
								},
								BoolExpr: nil,
								Next:     "State1",
							}},
							Default: "State3",
							CommonState2: CommonState2{
								CommonState1: CommonState1{
									Type: "Choice",
								},
							}}},
				},
				{
					State{"Pass", "State3", "",
						&PassState{
							CommonState4: CommonState4{
								CommonState3: CommonState3{
									End: true,
									CommonState2: CommonState2{
										CommonState1: CommonState1{
											Type: "Pass",
										}}}}}},
				},
			},
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
			if diff := cmp.Diff(gotStates, tt.wantStates); diff != "" {
				t.Errorf("Compile() = \n%#v\n want = \n%#v\n", gotStates, tt.wantStates)
				t.Errorf("====== diff ======\n%s\n", diff)
				t.Errorf(pp.Sprint(gotStates))
				t.Errorf(pp.Sprint(tt.wantStates))
			}
		})
	}
}
