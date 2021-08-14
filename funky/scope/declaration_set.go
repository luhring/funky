package scope

import (
	"go/ast"
)

type declsMap map[string]*ast.Ident

type declarationSet struct {
	decls declsMap
}

func newDeclarationSet() declarationSet {
	return declarationSet{decls: make(declsMap)}
}

func (s declarationSet) copy() declarationSet {
	newDecls := make(declsMap)

	for k, v := range s.decls {
		newDecls[k] = v
	}

	return declarationSet{
		decls: newDecls,
	}
}

func (s declarationSet) idents() []*ast.Ident {
	var result []*ast.Ident

	for _, ident := range s.decls {
		result = append(result, ident)
	}

	return result
}

func (s declarationSet) contains(i *ast.Ident) bool {
	if i == nil {
		return false
	}

	_, contains := s.decls[i.Name]
	return contains
}

func (s declarationSet) lookup(name string) *ast.Ident {
	if ident, ok := s.decls[name]; ok {
		return ident
	}

	return nil
}

func add(s declarationSet, identities ...*ast.Ident) declarationSet {
	result := s.copy()

	for _, identity := range identities {
		result.decls[identity.Name] = identity
	}

	return result
}
