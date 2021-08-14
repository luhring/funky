package main

import (
	"github.com/luhring/funky/funky/analyzers/stdlib"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(&stdlib.Analyzer)
}
