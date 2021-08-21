package ast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

func Render(node ast.Node, fset *token.FileSet) string {
	var buf bytes.Buffer

	err := format.Node(&buf, fset, node)
	if err != nil {
		return ""
	}

	return buf.String()
}
