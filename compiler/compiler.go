package compiler

import (
	"github.com/jingweno/godzilla/ast"
	"github.com/jingweno/godzilla/source"
)

func Compile(node ast.Node) *source.Code {
	code := source.NewCode()
	node.Compile(code)

	return code
}
