package mutation

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/luhring/funky/funky/assignment"
	funkyAST "github.com/luhring/funky/funky/ast"
	"github.com/luhring/funky/funky/finding"
	"github.com/luhring/funky/funky/scope"
)

var Type finding.Type = "mutation"

// Enforce that Mutation implements the Finding type
var _ finding.Finding = (*Mutation)(nil)

type Mutation struct {
	node           ast.Node
	mutatedVarExpr ast.Expr
	newValueExpr   ast.Expr
}

func (m Mutation) Message(fset *token.FileSet) string {
	var newValue string

	if m.newValueExpr != nil {
		newValue = funkyAST.Render(m.newValueExpr, fset)
	} else {
		newValue = "[unable to render expression]" // e.g. in the case of range stmt initializers
	}

	return fmt.Sprintf("%q was assigned a new value: %s", funkyAST.Render(m.mutatedVarExpr, fset), newValue)
}

func (m Mutation) Type() finding.Type {
	return Type
}

func (m Mutation) Node() ast.Node {
	return m.node
}

func (m Mutation) Location(fset *token.FileSet) finding.Location {
	return finding.Location(fset.Position(m.node.Pos()).String())
}

func (m Mutation) String() string {
	return fmt.Sprintf("%q mutated", varName(m))
}

func FindInFiles(files []*ast.File) []Mutation {
	var mutations []Mutation

	packageScope := scope.FromFiles(files)

	for _, file := range files {
		funkyAST.InspectWithInitialScope(file, func(node ast.Node, s scope.Scope) bool {
			if node == nil {
				return false
			}

			switch stmt := node.(type) {
			case *ast.AssignStmt:
				assignments := assignment.AssignmentsFromStmt(stmt)
				mutations = append(mutations, mutationsFromAssignments(assignments, stmt, s)...)

			case *ast.RangeStmt:
				assignments := assignment.AssignmentsFromRangeStmtInitializer(stmt)
				mutations = append(mutations, mutationsFromAssignments(assignments, stmt, scope.NewInsideExisting(s))...)
			}

			return true
		}, packageScope)
	}

	return mutations
}

func mutationsFromAssignments(assignments []assignment.Assignment, n ast.Node, s scope.Scope) []Mutation {
	var mutations []Mutation

	for _, a := range assignments {
		if isMutation(a, s) {
			mutation := Mutation{
				node:           n,
				mutatedVarExpr: a.VarExpr,
				newValueExpr:   a.NewValueExpr,
			}
			mutations = append(mutations, mutation)
		}
	}

	return mutations
}

func isMutation(a assignment.Assignment, s scope.Scope) bool {
	return !isDeclaration(a, s) && !funkyAST.BlankIdentifier(a.Ident)
}

func isDeclaration(a assignment.Assignment, s scope.Scope) bool {
	return !funkyAST.BlankIdentifier(a.Ident) &&
		!isFromImportedPackage(a.Ident, s) &&
		assignment.IsDefinition(a) &&
		!scope.InCurrent(s, a.Ident)
}

func isFromImportedPackage(expr ast.Expr, s scope.Scope) bool {
	ident := funkyAST.SelectorIdentFromExpr(expr)
	if ident == nil {
		return false
	}

	return scope.HasImport(s, ident.Name)
}

// varName returns the name of the variable being mutated
func varName(m Mutation) string {
	if ident := funkyAST.IdentFromExpr(m.mutatedVarExpr); ident != nil {
		return ident.Name
	}

	return "[unnamed variable]"
}

func Findings(mutations []Mutation) []finding.Finding {
	var findings []finding.Finding

	for _, m := range mutations {
		findings = append(findings, m)
	}

	return findings
}
