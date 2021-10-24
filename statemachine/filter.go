package statemachine

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spyzhov/ajson"
)

func filterByInputPath(input *ajson.Node, path string) (*ajson.Node, error) {
	return filterByJSONPath(input, path)
}

func filterByOutputPath(output *ajson.Node, path string) (*ajson.Node, error) {
	return filterByJSONPath(output, path)
}

func filterByJSONPath(input *ajson.Node, path string) (*ajson.Node, error) {
	switch path {
	case "", "$":
		return input, nil
	}

	nodes, err := input.JSONPath(path)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, ErrInvalidJsonPath
	}

	return nodes[0], nil
}

func filterByParameters(input *ajson.Node, parameters *json.RawMessage) (*ajson.Node, error) {
	n, err := repraceByRawJSON(input, parameters)
	if err != nil {
		panic(err.Error())
	}

	if errors.Is(err, ErrInvalidRawJSON) {
		err = fmt.Errorf("invalid Task Parameters: %v", err)
	}

	return n, err
}

func filterByResultSelector(output *ajson.Node, selector *json.RawMessage) (*ajson.Node, error) {
	n, err := repraceByRawJSON(output, selector)
	if errors.Is(err, ErrInvalidRawJSON) {
		err = fmt.Errorf("invalid ResulutSelector: %v", err)
	}

	return n, err
}

func repraceByRawJSON(node *ajson.Node, raw *json.RawMessage) (*ajson.Node, error) {
	if node == nil {
		return node, nil
	}

	if raw == nil {
		return node, nil
	}

	b, err := raw.MarshalJSON()
	if err != nil {
		return nil, err
	}

	root, err := ajson.Unmarshal(b)
	if err != nil {
		return nil, err
	}

	if !root.IsObject() {
		return nil, ErrInvalidRawJSON
	}

	m, err := root.GetObject()
	if err != nil {
		return nil, err
	}

	for k, n := range m {
		if strings.HasSuffix(k, ".$") {
			if !n.IsString() {
				continue
			}

			nodes, err := node.JSONPath(n.MustString())
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
