package scope

import "go/ast"

type Scope struct {
	current declarationSet
	outer   declarationSet
}

func New() Scope {
	return Scope{
		current: newDeclarationSet(),
	}
}

func Copy(s Scope) Scope {
	return Scope{
		current: s.current.copy(),
		outer:   s.outer.copy(),
	}
}

func Append(s Scope, identities ...*ast.Ident) Scope {
	return Scope{
		current: add(s.current, identities...),
		outer:   s.outer.copy(),
	}
}

func NewInsideExisting(existing Scope) Scope {
	// flatten existing scope and set to new scope's outer
	outer := existing.outer.copy()
	outer = add(outer, existing.current.idents()...)

	current := newDeclarationSet()

	return Scope{
		current: current,
		outer:   outer,
	}
}

func FromFiles(files []*ast.File) Scope {
	scope := Scope{
		current: newDeclarationSet(),
	}

	for _, file := range files {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				scope = Append(scope, d.Name)
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					if v, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range v.Names {
							scope = Append(scope, name)
						}
					}
				}
			}
		}
	}

	return scope
}

func InCurrent(s Scope, identity *ast.Ident) bool {
	return s.current.contains(identity)
}
