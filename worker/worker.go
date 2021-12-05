package worker

import (
	"bytes"
	"context"

	"github.com/k0kubun/pp"
	"github.com/w-haibara/kuirejo/compiler"
)

func Exec(ctx context.Context, w compiler.Workflow, input *bytes.Buffer) ([]byte, error) {
	_, _ = pp.Println(w)

	return nil, nil
}
