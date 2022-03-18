package worker

import (
	"fmt"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/w-haibara/kakemoti/compiler"
)

func workflowFields(w Workflow) log.Fields {
	return log.Fields{
		"id":      w.ID,
		"startat": w.StartAt,
		"timeout": w.TimeoutSeconds,
		"line":    Line(),
	}
}

func errorFields(err error) log.Fields {
	return log.Fields{
		"Error": err,
		"Line":  LineN(4),
	}
}

func stateFields(s compiler.State) log.Fields {
	return log.Fields{
		"Type": s.Common().Type,
		"Name": s.Name(),
		"Next": s.Next(),
		"Line": Line(),
	}
}

func Line() string {
	return LineN(3)
}

func LineN(n int) string {
	_, path, line, ok := runtime.Caller(n)
	if !ok {
		return "---"
	}

	_, file := filepath.Split(path)

	return fmt.Sprintf("%s:%d", file, line)
}
