# funky 🎸

A Go (golang) linter that finds mutations.

## Installation

```
go install github.com/luhring/funky/cmd/funky@main
```

## Usage

```
funky ./directory-with-go-files
```

(**Note:** Funky doesn't yet analyze directories _recursively_. 😳)

### Example output

```
add.go:134:3: mutation: "destfi" was assigned a new value: nil
add.go:155:5: mutation: "d" was assigned a new value: filepath.Join(dest, path.Base(url.Path))
add.go:157:7: mutation: "err" was assigned a new value: addURL(d, src, hostOwner, options.Hasher)
add.go:175:4: mutation: "err" was assigned a new value: os.Stat(esrc)
```

## What is a "mutation"?

A mutation is when a variable's value changes. In the Go language, this means an assignment of a value to a variable anywhere **other than** where that variable is declared.

For example, here's some Go code that has a few mutations:

```go
package main

var x = 1 // Not a mutation. This is where `x` is declared.

func main() {
    x = 7 // MUTATION! (... of `x`, which is declared as a package-level variable.)

    var output string // Not a mutation. This is where `output` is declared, and it's implicitly being assigned the zero-value of the `string` type, which is "".

    if someCondition {
    	output = "new value" // MUTATION!
    }

    print(output)
}
```


## The mission: functional programming for Go

Funky's objective is to take the approaches of functional programming and apply them to the Go language.

The first task is to alert developers to the presence of **mutations**.

Mutations introduce complexity in code. Because of this, they make code more likely to introduce bugs and more difficult to comprehend. Funky alerts you to the mutations in your code so you can spot opportunities for making your code less complex and more predictable — presumably, by adjusting your implementation to avoid mutations.

### "...but Go isn't a functional programming language"

That's true (to some extent). But more developers are realizing the benefits of applying coding practices _taken from the functional paradigm_ to more languages than just the esoteric FP languages.

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
