package fn

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

type Obj map[string]interface{}

func DoScriptTask(ctx context.Context, path string, in Obj) (Obj, error) {
	exe, err := exec.LookPath(path)
	if err != nil {
		return nil, err
	}

	v, ok := in["args"]
	if !ok {
		return nil, errors.New("'args' not found")
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil, errors.New("'args' type is not array")
	}
	args := make([]string, 0, len(arr))
	for _, a := range arr {
		args = append(args, fmt.Sprint(a))
	}

	cmd := exec.CommandContext(ctx, exe, args...) // #nosec G204
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return Obj{"Output": string(out)}, nil
}
