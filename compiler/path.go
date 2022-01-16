package compiler

import (
	"strings"

	"github.com/ohler55/ojg/jp"
)

type Path struct {
	Expr          jp.Expr
	IsContextPath bool
}

func NewPath(path string) (Path, error) {
	result := Path{}

	if strings.HasPrefix(path, "$$") {
		path = strings.TrimPrefix(path, "$")
		result.IsContextPath = true
	}

	p, err := jp.ParseString(path)
	if err != nil {
		return Path{}, err
	}
	result.Expr = p

	return result, nil
}

func MustNewPath(path string) Path {
	p, err := NewPath(path)
	if err != nil {
		panic(err.Error())
	}
	return p
}

func (p Path) String() string {
	return p.Expr.String()
}
