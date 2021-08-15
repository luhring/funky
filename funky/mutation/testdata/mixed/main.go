package main

import (
	"strings"

	"github.com/luhring/funky/funky/mutation/testdata/mixed/other"
)

func declarationStatement() {
	var x string
	print(x)

	x = "some new value for x!" // mutation
}

func shadowingInsideFuncLiteral() {
	var x string
	print(x)

	test := func() { // not a mutation
		x := "value" // not a mutation
		print(x)

		x = "changed" // mutation
	}
	print(test)
}

func twoMutations() {
	var x = 1
	print(x)

	x = 2 // mutation
	x = 3 // mutation
}

func ifStatement() {
	var x, y string
	print(x, y)

	if x = "test"; strings.Contains(x, "t") { // mutation

		x = "inside if" // mutation

		if y := "thing"; strings.Contains(y, "t") { // not a mutation
			y = "some" // mutation

			x := "" // not a mutation
			print(x)
		}
	} else {
		x = "inside else" // mutation
	}
}

func switchStatement() {
	var x string

	switch x = "something"; x { // mutation
	case "a":
		x = "inside case" // mutation
	default:
		x = "inside default" // mutation
	}

	switch x := "newVariable"; x { // not a mutation
	default:
		x := 1 // not a mutation
		print(x)
	}

	switch x := "new"; x {
	default:
		x = "again" // mutation
	}
}

func forStatement() {
	var x int
	print(x)

	for x := 1; x < 10; x++ { // not a mutation
		x := 200 // not a mutation
		print(x)
	}

	for x := 1; x < 10; x++ { // not a mutation
		x = 200 // mutation
		print(x)
	}

	for x = 1; x < 10; x++ { // mutation
		x := 2 // not a mutation
		print(x)
	}
}

func rangeStatement() {
	var thing string
	print(thing)

	var things []string

	for i, thing := range things {
		print(thing)

		thing := "newValue" // not a mutation
		print(i, thing)

		thing = "anotherValue" // mutation

		i = 0 // mutation
	}

	for _, thing = range things { // mutation
		print(thing)

		thing := "newValue" // not a mutation
		print(thing)
	}
}

func structLiteralWithFuncLiteral() {
	z := struct {
		a string
		b func(a string)
	}{
		a: "",
		b: func(a string) {
			b := "" // not a mutation
			print(a, b)

			a = "1" // mutation
			b = "2" // mutation
		},
	}
	print(z)
}

func usingGenDecls() {
	genDecl = "zzz"            // mutation
	genDecl, v := "xxx", "yyy" // no mutations
	print(genDecl, v)
}

var genDecl string

func typeSwitchStatement() {
	var x interface{}
	print(x)

	switch v := x.(type) {
	case string:
		v = "value" // mutation
		print(v)
	default:
		x := "brand new"
		print(x)
	}
}

func switchStmtWithCaseExpressions() {
	var x = 1

	switch {
	case false:
		x := 9
		print(x)

	case x > 0:
		x = 5 // mutation
	}
}

func exprStmt() {
	f := func(g func()) {
		g()
	}

	x := 0
	print(x)

	f(func() {
		x = 1 // mutation

		x := 2

		x = 3 // mutation
		print(x)
	})
}

func blockStmt() {
	var x int
	print(x)

	{
		x = 1 // mutation
		x, y := 2, 3
		print(x, y)
		x = 3 // mutation
	}
}

func otherPackage() {
	other.MyVar = "changed" // mutation
	MyVar := "value"        // not a mutation
	print(MyVar)
}

func funcWithParameters(a int) (result int) {
	print(a)

	a = 5      // mutation
	result = 7 // mutation

	a, b := 10, 11 // a is mutated, b isn't (because it's declared here)
	print(b)
	return
}

func main() {

}
