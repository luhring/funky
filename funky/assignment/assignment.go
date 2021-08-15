package assignment

import (
	"go/ast"
	"go/token"

	funkyAST "github.com/luhring/funky/funky/ast"
)

// Assignment describes an instance of a value being assigned to a variable.
type Assignment struct {
	Ident        *ast.Ident
	Token        token.Token
	VarExpr      ast.Expr
	NewValueExpr ast.Expr
}

func IsDefinition(a Assignment) bool {
	return a.Token == token.DEFINE && !funkyAST.BlankIdentifier(a.Ident)
}

func AssignmentsFromStmt(stmt *ast.AssignStmt) []Assignment {
	var result []Assignment

	for lhsIndex, expr := range stmt.Lhs {
		ident := funkyAST.IdentFromExpr(expr)

		newValueExpr := assignedValueExpr(stmt.Rhs, lhsIndex)

		result = append(result, Assignment{
			Ident:        ident,
			Token:        stmt.Tok,
			VarExpr:      expr,
			NewValueExpr: newValueExpr,
		})
	}

	return result
}

func AssignmentsFromRangeStmtInitializer(stmt *ast.RangeStmt) []Assignment {
	if stmt == nil {
		return nil
	}

	if stmt.Tok == token.ILLEGAL {
		return nil
	}

	var result []Assignment

	for _, expr := range []ast.Expr{stmt.Key, stmt.Value} {
		if expr != nil {
			assignment := assignmentFromRangeInitializedExpr(expr, stmt.Tok)
			result = append(result, assignment)
		}
	}

	return result
}

func assignmentFromRangeInitializedExpr(expr ast.Expr, tok token.Token) Assignment {
	return Assignment{
		Ident:        funkyAST.IdentFromExpr(expr),
		Token:        tok,
		VarExpr:      expr,
		NewValueExpr: nil,
	}
}

func assignedValueExpr(rhs []ast.Expr, lhsIndex int) ast.Expr {
	if len(rhs) == 1 {
		return rhs[0]
	}

	return rhs[lhsIndex]
}
