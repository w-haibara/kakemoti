package fn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func DoStatesFormat(ctx context.Context, args []interface{}) (interface{}, error) {
	var ErrStatesFormatFailed = errors.New("DoStatesFormat() failed")
	const (
		str1 = "{}"
		str2 = "\\{}"
	)

	if len(args) < 2 {
		return nil, ErrStatesFormatFailed
	}

	f, ok := args[0].(string)
	if !ok {
		return nil, ErrStatesFormatFailed
	}

	if strings.Count(f, str1)-strings.Count(f, str2) != len(args)-1 {
		return nil, ErrStatesFormatFailed
	}

	result := f
	for _, arg := range args[1:] {
		switch reflect.ValueOf(arg).Kind() {
		case reflect.Map, reflect.Struct:
			return nil, ErrStatesFormatFailed
		}

		str := result
		n := 0
		for {
			n1 := strings.Index(str, str2)
			n2 := strings.Index(str, str1)
			if n1 < 0 || n2 < n1 {
				break
			}
			n = n1 + 3
			str = str[n:]
		}

		result = result[:n] + strings.Replace(result[n:], str1, fmt.Sprint(arg), 1)
	}

	return result, nil
}

func DoStatesStringToJson(ctx context.Context, args []interface{}) (interface{}, error) {
	var ErrStatesStringToJsonFailed = errors.New("DoStatesStringToJson() failed")

	if len(args) < 1 {
		return nil, ErrStatesStringToJsonFailed
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, ErrStatesStringToJsonFailed
	}

	var v interface{}
	if err := json.Unmarshal([]byte(str), &v); err != nil {
		return nil, ErrStatesStringToJsonFailed
	}

	return v, nil
}

func DoStatesJsonToString(ctx context.Context, args []interface{}) (interface{}, error) {
	var ErrStatesJsonToStringFailed = errors.New("DoStatesJsonToString() failed")

	if len(args) < 1 {
		return nil, ErrStatesJsonToStringFailed
	}

	b, err := json.Marshal(args[0])
	if err != nil {
		return nil, ErrStatesJsonToStringFailed
	}

	return string(b), nil
}
