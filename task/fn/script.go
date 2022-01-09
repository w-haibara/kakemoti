package fn

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	scriptPrefix = "kakemoti_"
)

func DoScriptTask(ctx context.Context, path string, in Obj) (Obj, string, error) {
	exe, err := exec.LookPath(path)
	if err != nil {
		return nil, "", err
	}

	v, ok := in["args"]
	if !ok {
		return nil, "", errors.New("'args' not found")
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil, "", errors.New("'args' type is not array")
	}
	args := make([]string, 0, len(arr))
	for _, a := range arr {
		args = append(args, fmt.Sprint(a))
	}

	cmd := exec.CommandContext(ctx, exe, args...) // #nosec G204
	out, err := cmd.Output()
	if err != nil {
		return nil, "", err
	}

	output := Obj{}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.HasPrefix(line, scriptPrefix) {
			continue
		}

		s := strings.SplitN(strings.TrimPrefix(line, scriptPrefix), "=", 2)
		if len(s) < 2 {
			continue
		}
		output[s[0]] = s[1]
	}

	return output, "", nil
}
