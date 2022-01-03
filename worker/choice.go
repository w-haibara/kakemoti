package worker

import (
	"context"
	"errors"
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/w-haibara/kakemoti/compiler"
)

func (w Workflow) evalChoice(ctx context.Context, state *compiler.ChoiceState, input interface{}) (string, interface{}, error) {
	for _, choice := range state.Choices {
		/*
			if choice.BoolExpr != nil {
			// case "And":
			// case "Not":
			// case "Or":
			}
		*/

		if choice.Rule != nil {
			switch choice.Rule.Operator {
			case "BooleanEquals":
				next, output, err := BooleanEquals(choice, input)
				if err != nil {
					return "", nil, err
				}
				if next == "" {
					continue
				}
				return next, output, nil
			case "BooleanEqualsPath":
				panic("Not Implemented")
			case "IsBoolean":
				panic("Not Implemented")
			case "IsNull":
				panic("Not Implemented")
			case "IsNumeric":
				panic("Not Implemented")
			case "IsPresent":
				panic("Not Implemented")
			case "IsString":
				panic("Not Implemented")
			case "IsTimestamp":
				panic("Not Implemented")
			case "NumericEquals":
				panic("Not Implemented")
			case "NumericEqualsPath":
				panic("Not Implemented")
			case "NumericGreaterThan":
				panic("Not Implemented")
			case "NumericGreaterThanPath":
				panic("Not Implemented")
			case "NumericGreaterThanEquals":
				panic("Not Implemented")
			case "NumericGreaterThanEqualsPath":
				panic("Not Implemented")
			case "NumericLessThan":
				panic("Not Implemented")
			case "NumericLessThanPath":
				panic("Not Implemented")
			case "NumericLessThanEquals":
				panic("Not Implemented")
			case "NumericLessThanEqualsPath":
				panic("Not Implemented")
			case "StringEquals":
				panic("Not Implemented")
			case "StringEqualsPath":
				panic("Not Implemented")
			case "StringGreaterThan":
				panic("Not Implemented")
			case "StringGreaterThanPath":
				panic("Not Implemented")
			case "StringGreaterThanEquals":
				panic("Not Implemented")
			case "StringGreaterThanEqualsPath":
				panic("Not Implemented")
			case "StringLessThanStringLessThanPath":
				panic("Not Implemented")
			case "StringLessThanEqualsStringLessThanEqualsPath":
				panic("Not Implemented")
			case "StringMatches":
				panic("Not Implemented")
			case "TimestampEquals":
				panic("Not Implemented")
			case "TimestampEqualsPath":
				panic("Not Implemented")
			case "TimestampGreaterThan":
				panic("Not Implemented")
			case "TimestampGreaterThanPath":
				panic("Not Implemented")
			case "TimestampGreaterThanEquals":
				panic("Not Implemented")
			case "TimestampGreaterThanEqualsPath":
				panic("Not Implemented")
			case "TimestampLessThan":
				panic("Not Implemented")
			case "TimestampLessThanPath":
				panic("Not Implemented")
			case "TimestampLessThanEquals":
				panic("Not Implemented")
			case "TimestampLessThanEqualsPath":
				panic("Not Implemented")
			}
		}
	}

	return state.Default, input, nil
}

func BooleanEquals(choice compiler.Choice, input interface{}) (string, interface{}, error) {
	path, ok := choice.Rule.Variable1.(string)
	if !ok {
		return "", nil, errors.New("type of choice.Rule.Variable1 is not string")
	}

	p, err := jp.ParseString(path)
	if err != nil {
		return "", nil, fmt.Errorf("jp.ParseString(path) failed: %w", err)
	}
	nodes := p.Get(input)

	if len(nodes) != 1 {
		return "", nil, fmt.Errorf("invalid length of path.Get(input) result")
	}

	v1, ok := nodes[0].(bool)
	if !ok {
		return "", nil, fmt.Errorf("invalid type of path.Get(input) result")
	}

	v2, ok := choice.Rule.Variable2.(bool)
	if !ok {
		return "", nil, errors.New("type of choice.Rule.Variable2 is not bool")
	}

	if v1 == v2 {
		return choice.Next, input, nil
	}

	return "", nil, nil
}
