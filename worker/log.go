package worker

import (
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
