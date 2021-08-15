package mutation

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	funkyAST "github.com/luhring/funky/funky/ast"
	"github.com/luhring/funky/funky/finding"
)

func TestFindInFiles(t *testing.T) {
	expected := []testableMutation{
		{
			location:         "testdata/mixed/main.go:13:2",
			variableName:     "x",
			newValueRendered: "\"some new value for x!\"",
		},
		{
			location:         "testdata/mixed/main.go:24:3",
			variableName:     "x",
			newValueRendered: "\"changed\"",
		},
		{
			location:         "testdata/mixed/main.go:33:2",
			variableName:     "x",
			newValueRendered: "2",
		},
		{
			location:         "testdata/mixed/main.go:34:2",
			variableName:     "x",
			newValueRendered: "3",
		},
		{
			location:         "testdata/mixed/main.go:41:5",
			variableName:     "x",
			newValueRendered: "\"test\"",
		},
		{
			location:         "testdata/mixed/main.go:43:3",
			variableName:     "x",
			newValueRendered: "\"inside if\"",
		},
		{
			location:         "testdata/mixed/main.go:46:4",
			variableName:     "y",
			newValueRendered: "\"some\"",
		},
		{
			location:         "testdata/mixed/main.go:52:3",
			variableName:     "x",
			newValueRendered: "\"inside else\"",
		},
		{
			location:         "testdata/mixed/main.go:59:9",
			variableName:     "x",
			newValueRendered: "\"something\"",
		},
		{
			location:         "testdata/mixed/main.go:61:3",
			variableName:     "x",
			newValueRendered: "\"inside case\"",
		},
		{
			location:         "testdata/mixed/main.go:63:3",
			variableName:     "x",
			newValueRendered: "\"inside default\"",
		},
		{
			location:         "testdata/mixed/main.go:74:3",
			variableName:     "x",
			newValueRendered: "\"again\"",
		},
		{
			location:         "testdata/mixed/main.go:88:3",
			variableName:     "x",
			newValueRendered: "200",
		},
		{
			location:         "testdata/mixed/main.go:92:6",
			variableName:     "x",
			newValueRendered: "1",
		},
		{
			location:         "testdata/mixed/main.go:110:3",
			variableName:     "thing",
			newValueRendered: "\"anotherValue\"",
		},
		{
			location:         "testdata/mixed/main.go:112:3",
			variableName:     "i",
			newValueRendered: "0",
		},
		{
			location:         "testdata/mixed/main.go:115:2",
			variableName:     "thing",
			newValueRendered: "",
		},
		{
			location:         "testdata/mixed/main.go:133:4",
			variableName:     "a",
			newValueRendered: "\"1\"",
		},
		{
			location:         "testdata/mixed/main.go:134:4",
			variableName:     "b",
			newValueRendered: "\"2\"",
		},
		{
			location:         "testdata/mixed/main.go:141:2",
			variableName:     "genDecl",
			newValueRendered: "\"zzz\"",
		},
		{
			location:         "testdata/mixed/main.go:154:3",
			variableName:     "v",
			newValueRendered: "\"value\"",
		},
		{
			location:         "testdata/mixed/main.go:171:3",
			variableName:     "x",
			newValueRendered: "5",
		},
		{
			location:         "testdata/mixed/main.go:184:3",
			variableName:     "x",
			newValueRendered: "1",
		},
		{
			location:         "testdata/mixed/main.go:188:3",
			variableName:     "x",
			newValueRendered: "3",
		},
		{
			location:         "testdata/mixed/main.go:198:3",
			variableName:     "x",
			newValueRendered: "1",
		},
		{
			location:         "testdata/mixed/main.go:201:3",
			variableName:     "x",
			newValueRendered: "3",
		},
		{
			location:         "testdata/mixed/main.go:206:2",
			variableName:     "other.MyVar",
			newValueRendered: "\"changed\"",
		},
		{
			location:         "testdata/mixed/main.go:214:2",
			variableName:     "a",
			newValueRendered: "5",
		},
		{
			location:         "testdata/mixed/main.go:215:2",
			variableName:     "result",
			newValueRendered: "7",
		},
		{
			location:         "testdata/mixed/main.go:217:2",
			variableName:     "a",
			newValueRendered: "10",
		},
		{
			location:         "testdata/mixed/main.go:228:2",
			variableName:     "other.name",
			newValueRendered: "\"lore\"",
		},
		{
			location:         "testdata/mixed/main.go:231:2",
			variableName:     "some.name",
			newValueRendered: "\"barney\"",
		},
	}

	fset := token.NewFileSet()
	packages := loadGoSourceTestFixture(t, fset, "mixed")
	mainPackage := packages["main"]

	files := funkyAST.SortedFilesFromPackage(mainPackage)

	mutations := FindInFiles(files)

	actual := mapToTestableMutations(mutations, fset)

	assertEqualTestableMutationSets(
		t,
		newTestableMutationSet(expected),
		newTestableMutationSet(actual),
	)
}

type testableMutationSet map[testableMutation]struct{}

func newTestableMutationSet(mutations []testableMutation) testableMutationSet {
	set := make(testableMutationSet)

	for _, mutation := range mutations {
		set[mutation] = struct{}{}
	}

	return set
}

func assertEqualTestableMutationSets(t *testing.T, expected, actual testableMutationSet) {
	t.Helper()

	if expected == nil {
		t.Fatalf("expected was nil")
	}

	if actual == nil {
		t.Fatalf("actual was nil")
	}

	for mutation := range expected {
		if _, exists := actual[mutation]; !exists {
			t.Errorf("expected mutation was absent: %s", reportMutation(mutation))
		}
	}

	for mutation := range actual {
		if _, exists := expected[mutation]; !exists {
			t.Errorf("unexpected mutation was present: %s", reportMutation(mutation))
		}
	}
}

func reportMutation(m testableMutation) string {
	return fmt.Sprintf("%s —— %q -> %s", m.location, m.variableName, m.newValueRendered)
}

// TODO: move to more general location
func loadGoSourceTestFixture(t testing.TB, fset *token.FileSet, fixtureDirectory string) map[string]*ast.Package {
	t.Helper()
	const errMessage = "unable to load Go source test fixture: %v"

	dir := "testdata/" + fixtureDirectory
	packages, err := parser.ParseDir(fset, dir, nil, parser.AllErrors)
	if err != nil {
		t.Fatalf(errMessage, err)
	}

	return packages
}

type testableMutation struct {
	location         finding.Location
	variableName     string
	newValueRendered interface{}
}

func mapToTestableMutations(mutations []Mutation, fset *token.FileSet) []testableMutation {
	var result []testableMutation

	for _, mutation := range mutations {
		result = append(result, mapToTestableMutation(mutation, fset))
	}

	return result
}

func mapToTestableMutation(mutation Mutation, fset *token.FileSet) testableMutation {
	location := finding.Location(fset.Position(mutation.node.Pos()).String())

	return testableMutation{
		location:         location,
		variableName:     funkyAST.Render(mutation.mutatedVarExpr, fset),
		newValueRendered: funkyAST.Render(mutation.newValueExpr, fset),
	}
}
