# store4

[![MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE.md)
[![Build Status](https://github.com/jimsmart/store4/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/jimsmart/store4/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/jimsmart/store4/branch/master/graph/badge.svg)](https://codecov.io/gh/jimsmart/store4)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimsmart/store4)](https://goreportcard.com/report/github.com/jimsmart/store4)
[![Used By](https://img.shields.io/sourcegraph/rrc/github.com/jimsmart/store4.svg)](https://sourcegraph.com/github.com/jimsmart/store4)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/jimsmart/store4)

store4 is a [Go](https://golang.org) package providing a fast in-memory quad store, with graph and subject views.

## Installation

```bash
go get github.com/jimsmart/store4
```

```go
import "github.com/jimsmart/store4"
```

### Dependencies

- Standard library.
- [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/) if you wish to run the tests.

## Example

```go
import "github.com/jimsmart/store4"

s := store4.NewQuadStore()
s.Add("Alice", "knows", "Bob", "")
s.Add("Alice", "knows", "Charlie", "")
s.Add("Charlie", "knows", "Bob", "")
s.Add("Charlie", "age", 23, "")

// Find everyone that Alice knows, in any graph.
x := s.FindObjects("Alice", "knows", "*")

// Find everyone who knows Charlie, in the unnamed/default graph.
y := s.FindSubjects("knows", "Charlie", "")

// Iterate over all quads.
s.ForEach(func(s, p string, o interface{}, g string) {
    // ...
})

// Iterate over quads matching given pattern.
s.ForEachWith("*", "*", "Bob", "*", func(s, p string, o interface{}, g string) {
    // ...
})

// Remove all statements about Charlie, from all graphs.
s.Remove("Charlie", "*", "*", "*")
```

See GoDocs for more detailed examples.

## Documentation

GoDocs [https://godoc.org/github.com/jimsmart/store4](https://godoc.org/github.com/jimsmart/store4)

## Testing

Package store4 is extensively tested:

- 200+ Gingko specs (see **_test.go*)
- Example code for most methods, with verified output (see **_examples_test.go*)

To run the tests execute `go test` inside the project folder.

For a full coverage report, try:

```bash
go test -coverprofile=coverage.out && go tool cover -html=coverage.out
```

## License

Package store4 is copyright 2016 by Jim Smart and released under the [MIT License](LICENSE.md)

### Additional credits

The internals of QuadStore draw heavily from the implementation of N3Store, a component of [N3.js](https://github.com/RubenVerborgh/N3.js). The N3.js library is copyrighted by Ruben Verborgh and released under the MIT License.

## History

- v0.0.1 (2021-04-19) Use Go modules. Enable CI using GitHub Actions. Remove Travis.
