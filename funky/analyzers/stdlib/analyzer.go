package stdlib

import (
	"flag"
	"fmt"
	"go/token"

	"github.com/luhring/funky/funky/finding"
	"github.com/luhring/funky/funky/mutation"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = analysis.Analyzer{
	Name:  "funky",
	Doc:   "enforces functional programming best practices such as avoiding mutations",
	Flags: flag.FlagSet{},
	Run: func(pass *analysis.Pass) (interface{}, error) {
		return run(pass)
	},
	RunDespiteErrors: true,
	Requires:         nil,
	ResultType:       nil,
	FactTypes:        nil,
}

func run(pass *analysis.Pass) (interface{}, error) {
	mutations := mutation.FindInFiles(pass.Files)

	for _, m := range mutations {
		diagnostic := diagnostic(m, pass.Fset)

		if report := pass.Report; report != nil {
			report(diagnostic)
		}
	}

	return nil, nil
}

func diagnostic(f finding.Finding, fset *token.FileSet) analysis.Diagnostic {
	return analysis.Diagnostic{
		Pos:      f.Node().Pos(),
		End:      f.Node().End(),
		Category: string(f.Type()),
		Message:  message(f, fset),
	}
}

func message(f finding.Finding, fset *token.FileSet) string {
	return fmt.Sprintf("%s: %s", f.Type(), f.Message(fset))
}
