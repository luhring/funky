package ast

import "go/ast"

func BlankIdentifier(ident *ast.Ident) bool {
	if ident == nil {
		return false
	}

	return ident.Name == "_"
}

func SelectorIdentFromExpr(expr ast.Expr) *ast.Ident {
	if expr == nil {
		return nil
	}

	if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
		if ident, ok := selectorExpr.X.(*ast.Ident); ok {
			return ident
		}

		return SelectorIdentFromExpr(selectorExpr.X)
	}

	return nil
}
