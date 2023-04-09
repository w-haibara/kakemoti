package fn

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

const (
	scriptInputPrefix  = "KAKEMOTI_IN"
	scriptOutputPrefix = "KAKEMOTI_OUT"
	scriptErrorPrefix  = "KAKEMOTI_ERR"
)

func DoScriptTask(ctx context.Context, path string, in Obj) (Obj, string, error) {
	exe, err := exec.LookPath(path)
	if err != nil {
		return nil, "", err
	}

	args, ok := in["args"]
	if !ok {
		return nil, "", errors.New("'args' not found")
	}

	cmd := exec.CommandContext(ctx, exe) // #nosec G204

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, marshalArgs(args)...)

	out, err := cmd.Output()
	if err != nil {
		return nil, "", err
	}

	output := Obj{}
	stateserror := ""
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, scriptErrorPrefix+"=") {
			stateserror = strings.TrimPrefix(line, scriptErrorPrefix+"=")
			continue
		}

		if !strings.HasPrefix(line, scriptOutputPrefix) {
			continue
		}

		s := strings.SplitN(strings.TrimPrefix(line, scriptOutputPrefix), "=", 2)
		if len(s) < 2 {
			continue
		}

		s[0] = strings.TrimPrefix(s[0], "_")

		output[s[0]] = s[1]
	}

	return output, stateserror, nil
}

func marshalArgs(args interface{}) []string {
	result := []string{}
	for k, v := range marshalArgsToMap(scriptInputPrefix, args) {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

func marshalArgsToMap(key string, args interface{}) Obj {
	switch args.(type) {
	case []Obj:
		break
	default:
		if reflect.TypeOf(args).Kind() == reflect.Slice {
			v := reflect.ValueOf(args)
			s := make([]string, v.Len())
			for i := 0; i < v.Len(); i++ {
				s[i] = fmt.Sprint(v.Index(i))
			}
			args = s
		}
	}

	switch args := args.(type) {
	case Obj:
		result := make(Obj)
		for k1, arg := range args {
			for k2, v := range marshalArgsToMap(fmt.Sprintf("%s_%v", key, k1), arg) {
				result[k2] = v
			}
		}
		return result
	case []Obj:
		result := make(Obj)
		for i, arg := range args {
			for k, v := range marshalArgsToMap(key+"_"+strconv.Itoa(i), arg) {
				result[k] = v
			}
		}
		return result
	case []string:
		result := make(Obj)
		for i, arg := range args {
			result[key+"_"+strconv.Itoa(i)] = arg
		}
		return result
	default:
		return Obj{
			key: fmt.Sprint(args),
		}
	}
}
