package statemachine

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/spyzhov/ajson"
)

func filterByInputPath(input *ajson.Node, path string) (*ajson.Node, error) {
	switch path {
	case "", "$":
		return input, nil
	}

	return filterNode(input, path)
}

func filterByOutputPath(output *ajson.Node, path string) (*ajson.Node, error) {
	switch path {
	case "", "$":
		return output, nil
	}

	return filterNode(output, path)
}

func filterByResultSelector(output *ajson.Node, selector *json.RawMessage) (*ajson.Node, error) {
	if selector == nil {
		return output, nil
	}

	b, err := selector.MarshalJSON()
	if err != nil {
		return nil, err
	}

	root, err := ajson.Unmarshal(b)
	if err != nil {
		return nil, err
	}

	if !root.IsObject() {
		return nil, ErrInvalidResultSelector
	}

	m, err := root.GetObject()
	if err != nil {
		return nil, err
	}

	for k, node := range m {
		if strings.HasSuffix(k, ".$") {
			if !node.IsString() {
				continue
			}

			nodes, err := output.JSONPath(node.MustString())
			if err != nil {
				continue
			}

			if len(nodes) == 0 {
				continue
			}

			delete(m, k)
			k = strings.TrimSuffix(k, ".$")
			m[k] = nodes[0]
		}
	}

	return ajson.ObjectNode("", m), nil
}

func filterNode(input *ajson.Node, path string) (*ajson.Node, error) {
	nodes, err := input.JSONPath(path)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, ErrInvalidInputPath
	}

	return nodes[0], nil
}

func filterByResultPath(input, result *ajson.Node, path string) (*ajson.Node, error) {
	switch path {
	case "", "$":
		return result, nil
	case "null":
		return input, nil
	}

	return insertNode(input, result, path)
}

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
