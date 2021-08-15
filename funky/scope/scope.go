package scope

import (
	"go/ast"
	"path"
)

type Scope struct {
	current declarationSet
	outer   declarationSet
	imports []string
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
		imports: s.imports,
	}
}

func Append(s Scope, identities ...*ast.Ident) Scope {
	return Scope{
		current: add(s.current, identities...),
		outer:   s.outer.copy(),
		imports: s.imports,
	}
}

func WithImports(s Scope, imports []*ast.ImportSpec) Scope {
	var imported []string

	for _, spec := range imports {
		imported = append(imported, importName(spec))
	}

	return Scope{
		current: s.current.copy(),
		outer:   s.outer.copy(),
		imports: imported,
	}
}

func importName(spec *ast.ImportSpec) string {
	if spec.Name != nil {
		return spec.Name.Name
	}

	return path.Base(spec.Path.Value)
}

func NewInsideExisting(existing Scope) Scope {
	// flatten existing scope and set to new scope's outer
	outer := existing.outer.copy()
	outer = add(outer, existing.current.idents()...)

	current := newDeclarationSet()

	return Scope{
		current: current,
		outer:   outer,
		imports: existing.imports,
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

func HasImport(s Scope, name string) bool {
	for _, scopeImport := range s.imports {
		if scopeImport == name {
			return true
		}
	}

	return false
}
