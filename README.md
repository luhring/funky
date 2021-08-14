# funky 🎸

a Go linter for functional programming

## Installation

```
go install github.com/luhring/funky/cmd/funky@main
```

## Usage

```
funky ./directory-with-go-files
```

(Funky doesn't yet analyze directories **recursively** 😳 — this feature is at the top of the roadmap.)

## Why do I need this?

Like any opinionated linter, Funky **might add no value** for you!

Funky is a reaction to the recent rise in popularity of functional programming (FP). If you're interested in functional programming, and you're writing code in Go, Funky intends to help you assess your code's FP fitness and see opportunities for improvement.

### "...but Go isn't a functional programming language"

That's true (to some extent). But more and more, developers are realizing the benefits of applying coding practices _taken from the functional paradigm_ to more languages than just the canonical or esoteric FP ones. In particular, Funky is initially focused on helping Go developers detect and avoid **mutations** and **side effects** in their Go code — both of which can lead to bugs and difficulty in code comprehension.


## Roadmap

- [ ] support for recursive directory analysis (e.g. `./...`)
- [ ] CI pipeline
- [ ] release pipeline
- [ ] **feature:** avoiding mutations
  - [x] mutation detection
  - [ ] failure on mutation detection
  - [ ] configurable exceptions to mutation detection-based failing
- [ ] **feature:** avoiding side effects
  - [ ] side effect detection
  - [ ] failure on side effect detection
  - [ ] configurable exceptions to side effect detection-based failing

## Concepts

### Mutation

_"an instance of a variable assigned a value anywhere besides where that variable is declared"_

For example, here's how Funky views the following Go code:

```go
package main

var x int // not a mutation

func main() {
    x = 7 // mutation (... of x, which is declared as a package-level var)
    
    x := 8 // not a mutation (this is a new declaration of x)
    
    x = 9 // mutation (...of the NEW x, which is declared as a function-scoped var just above)
}
```
