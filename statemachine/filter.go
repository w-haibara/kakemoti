package statemachine

import (
	"github.com/google/uuid"
	"github.com/spyzhov/ajson"
)

func insertNode(n1, n2 *ajson.Node, path string) (*ajson.Node, error) {
	cmds, err := ajson.ParseJSONPath(path)
	if err != nil {
		return nil, err
	}

	if len(cmds) < 1 {
		return nil, ErrInvalidJsonPath
	}

	if cmds[0] != "$" {
		return nil, ErrInvalidJsonPath
	}

	root := n1.Clone()
	node := root
	for i, cmd := range cmds[1:] {
		n := map[string]*ajson.Node{}
		if i+3 == len(cmds) {
			n[cmds[len(cmds)-1]] = n2
		}

		cur := ajson.ObjectNode(uuid.New().String(), n)
		if err := node.AppendObject(cmd, cur); err != nil {
			return nil, err
		}

		if i+3 == len(cmds) {
			break
		}

		node = cur
	}

	return root, nil
}
