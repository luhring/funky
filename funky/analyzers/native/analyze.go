package native

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	funkyAST "github.com/luhring/funky/funky/ast"
	"github.com/luhring/funky/funky/finding"
	"github.com/luhring/funky/funky/mutation"
)

func AnalyzeDirectory(path string) error {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, path, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	for _, f := range findings(packages) {
		fmt.Println(finding.Report(f, fset))
	}

	return nil
}

func findings(packages map[string]*ast.Package) []finding.Finding {
	var findings []finding.Finding

	for _, packageNode := range packages {
		if packageNode == nil {
			continue
		}

		packageFiles := funkyAST.SortedFilesFromPackage(packageNode)
		mutations := mutation.FindInFiles(packageFiles)
		mutationFindings := mutation.Findings(mutations)

		findings = append(findings, mutationFindings...)
	}

	return findings
}
