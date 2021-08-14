package finding

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Type string

type Finding interface {
	fmt.Stringer
	Type() Type
	Node() ast.Node
	Location(*token.FileSet) Location
	Message(*token.FileSet) string
}

func Report(f Finding, fset *token.FileSet) string {
	return fmt.Sprintf("%s: %s: %s", f.Location(fset), f.Type(), f.Message(fset))
}
