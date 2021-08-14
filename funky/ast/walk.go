package ast

import (
	"go/ast"
	"go/token"
	"sort"

	"github.com/luhring/funky/funky/scope"
)

type Visitor interface {
	Visit(node ast.Node, scope scope.Scope) (w Visitor)
}

type inspector func(ast.Node, scope.Scope) bool

func (f inspector) Visit(node ast.Node, scope scope.Scope) Visitor {
	if f(node, scope) {
		return f
	}
	return nil
}

func Inspect(node ast.Node, f func(ast.Node, scope.Scope) bool) {
	InspectWithInitialScope(node, f, scope.New())
}

func InspectWithInitialScope(node ast.Node, f func(ast.Node, scope.Scope) bool, initialScope scope.Scope) {
	Walk(inspector(f), node, initialScope)
}

func Walk(v Visitor, node ast.Node, existingScope scope.Scope) {
	if v = v.Visit(node, existingScope); v == nil {
		return
	}

	s := scope.Copy(existingScope)

	// TODO: consider nil check for node

	switch n := node.(type) {
	case *ast.Package:
		// TODO: Make a decision on how to handle imports

		files := SortedFilesFromPackage(n)

		// determine declarations in package scope
		packageScope := scope.FromFiles(files)

		// Walk child nodes of package -> files
		for _, f := range files {
			Walk(v, f, packageScope)
		}

	case *ast.File:
		// No new scope to gather
		// Walk child nodes of file -> top-level declarations
		for _, decl := range n.Decls {
			Walk(v, decl, s)
		}

	case *ast.FuncDecl:
		funcScope := scope.NewInsideExisting(s)

		// what gets added to scope before further walking?
		// - receiver
		if n.Recv != nil {
			for _, field := range n.Recv.List {
				funcScope = scope.Append(funcScope, field.Names...)
			}
		}

		idents := declarationIdentsFromFuncType(n.Type)
		funcScope = scope.Append(funcScope, idents...)

		// Walk child notes of function -> body
		if funcBody := n.Body; funcBody != nil {
			Walk(v, funcBody, funcScope)
		}

	case *ast.BlockStmt:
		blockStmtScope := scope.NewInsideExisting(s)
		walkStmtList(v, n.List, blockStmtScope)

	case *ast.ForStmt:
		// grab new scope from init stmt
		if n.Init == nil {
			walkBlockStmt(v, n.Body, s)
			break
		}

		initScope := scope.NewInsideExisting(s)
		Walk(v, n.Init, initScope)

		if initStmt, ok := n.Init.(*ast.AssignStmt); ok {
			idents := declarationIdentsFromAssignStmt(initStmt, s)
			s = scope.Append(s, idents...)
		}

		walkBlockStmt(v, n.Body, initScope)

	case *ast.RangeStmt:
		if n.Tok == token.ILLEGAL { // e.g. `for range someSlice { ...`
			walkBlockStmt(v, n.Body, scope.NewInsideExisting(s))
			break
		}

		initScope := scope.NewInsideExisting(s)

		if n.Key != nil {
			Walk(v, n.Key, initScope)
		}

		if n.Value != nil {
			Walk(v, n.Value, initScope)
		}

		if ident := IdentFromExpr(n.Key); ident != nil {
			initScope = scope.Append(initScope, ident)
		}

		if ident := IdentFromExpr(n.Value); ident != nil {
			initScope = scope.Append(initScope, ident)
		}

		walkBlockStmt(v, n.Body, scope.NewInsideExisting(initScope))

	case *ast.IfStmt:
		if n.Init == nil {
			walkIfBlocks(v, n, s)
			break
		}

		initScope := scope.NewInsideExisting(s)
		Walk(v, n.Init, initScope)

		if initStmt, ok := n.Init.(*ast.AssignStmt); ok {
			idents := declarationIdentsFromAssignStmt(initStmt, initScope)
			initScope = scope.Append(initScope, idents...)
		}

		walkIfBlocks(v, n, initScope)

	case *ast.SwitchStmt:
		if n.Init == nil {
			walkBlockStmt(v, n.Body, scope.NewInsideExisting(s))
			break
		}

		initScope := scope.NewInsideExisting(s)
		Walk(v, n.Init, initScope)

		if initStmt, ok := n.Init.(*ast.AssignStmt); ok {
			idents := declarationIdentsFromAssignStmt(initStmt, initScope)
			initScope = scope.Append(initScope, idents...)
		}

		walkBlockStmt(v, n.Body, scope.NewInsideExisting(initScope))

	case *ast.TypeSwitchStmt:
		if n.Init != nil {
			if initStmt, ok := n.Init.(*ast.AssignStmt); ok {
				idents := declarationIdentsFromAssignStmt(initStmt, s)
				s = scope.Append(s, idents...)
			}
		}

		if n.Body != nil {
			Walk(v, n.Body, s)
		}

	case *ast.CaseClause:
		walkStmtList(v, n.Body, scope.NewInsideExisting(s))

	case *ast.FuncLit:
		funcScope := scope.NewInsideExisting(s)

		idents := declarationIdentsFromFuncType(n.Type)
		funcScope = scope.Append(funcScope, idents...)

		if funcBody := n.Body; funcBody != nil {
			Walk(v, funcBody, funcScope)
		}

	case *ast.CompositeLit:
		for _, expr := range n.Elts {
			Walk(v, expr, s)
		}

	case *ast.KeyValueExpr:
		Walk(v, n.Key, s)
		Walk(v, n.Value, s)

	case *ast.DeclStmt:
		Walk(v, n.Decl, s)

	case *ast.GoStmt:
		Walk(v, n.Call, s)

	case *ast.DeferStmt:
		Walk(v, n.Call, s)

	case *ast.ReturnStmt:
		for _, result := range n.Results {
			Walk(v, result, s)
		}

	case *ast.CallExpr:
		for _, arg := range n.Args {
			Walk(v, arg, s)
		}

	case *ast.ParenExpr:
		Walk(v, n.X, s)

	case *ast.ExprStmt:
		Walk(v, n.X, s)

	}

	v.Visit(nil, s)
}

func walkIfBlocks(v Visitor, ifStmt *ast.IfStmt, s scope.Scope) {
	if ifStmt.Body != nil {
		Walk(v, ifStmt.Body, scope.NewInsideExisting(s))
	}

	if ifStmt.Else != nil {
		Walk(v, ifStmt.Else, scope.NewInsideExisting(s))
	}
}

func walkBlockStmt(v Visitor, blockStmt *ast.BlockStmt, s scope.Scope) {
	if blockStmt == nil {
		return
	}

	Walk(v, blockStmt, s)
}

func declarationIdentsFromFuncType(funcType *ast.FuncType) []*ast.Ident {
	if funcType == nil {
		return nil
	}

	var idents []*ast.Ident

	// - input parameters
	if params := funcType.Params; params != nil {
		for _, parameter := range params.List {
			for _, parameterIdent := range parameter.Names {
				idents = append(idents, parameterIdent)
			}
		}
	}

	// - output parameters
	if results := funcType.Results; results != nil {
		for _, param := range results.List {
			idents = append(idents, param.Names...)
		}
	}

	return idents
}

func walkStmtList(v Visitor, list []ast.Stmt, s scope.Scope) {
	// Note: stmts can add to scope for subsequent stmts

	for _, stmt := range list {
		Walk(v, stmt, s)

		// There are two cases where a given statement can add to the scope of subsequent statements in this func block
		switch statement := stmt.(type) {

		// (e.g. `var/const foo string`)
		case *ast.DeclStmt:
			if genDecl, ok := statement.Decl.(*ast.GenDecl); ok {
				for _, spec := range genDecl.Specs {
					if v, ok := spec.(*ast.ValueSpec); ok {
						s = scope.Append(s, v.Names...)
					}
				}
			}

		// (e.g. `a, b := bar()`) <- NEED TO KNOW WHICH LHS ITEMS ARE NEW!
		case *ast.AssignStmt:
			for _, expr := range statement.Rhs {
				Walk(v, expr, s)
			}

			s = scope.Append(s, declarationIdentsFromAssignStmt(statement, s)...)
		}
	}
}

func IdentFromExpr(expr ast.Expr) *ast.Ident {
	if expr == nil {
		return nil
	}

	if ident, ok := expr.(*ast.Ident); ok {
		return ident
	}

	return nil
}

func declarationIdentsFromAssignStmt(statement *ast.AssignStmt, currentScope scope.Scope) []*ast.Ident {
	if statement.Tok != token.DEFINE {
		return nil
	}

	var idents []*ast.Ident

	for _, expr := range statement.Lhs {
		if ident := IdentFromExpr(expr); ident != nil {
			// if ident is already in scope, then it's not being declared by this assignStmt
			if scope.InCurrent(currentScope, ident) {
				continue
			}

			idents = append(idents, ident)
		}
	}

	return idents
}

func SortedFilesFromPackage(p *ast.Package) []*ast.File {
	if p == nil {
		return nil
	}

	var filenames []string

	for filename := range p.Files {
		filenames = append(filenames, filename)
	}

	sort.Strings(filenames)

	var files []*ast.File

	for _, filename := range filenames {
		files = append(files, p.Files[filename])
	}

	return files
}
