# store4

Package store4 provides a fast in-memory string-based quad store, written in [Go](https://golang.org).

## Installation
```bash
$ go get github.com/jimsmart/store4
```

```go
import "github.com/jimsmart/store4"
```

Package store4 has no external dependencies (except [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/), if you wish to run the tests).

## Example

```go
import "github.com/jimsmart/store4"

s := store4.NewQuadStore()
s.Add("Alice", "knows", "Bob", "")
s.Add("Alice", "knows", "Charlie", "")
s.Add("Charlie", "knows", "Bob", "")

// Find everyone that Alice knows, in any graph.
list := s.FindObjects("Alice", "knows", "*")
fmt.Println(list)
// Output (exact order may vary):
// [Bob Charlie]

// Find everyone who knows Charlie, in the unnamed/default graph.
x := s.FindSubjects("knows", "Charlie", "")
fmt.Println(x)
// Output:
// [Alice]

// Iterate over all quads.
s.ForEach(func(s, p, o, g string) {
    // ...
})

// Iterate over quads matching given pattern.
s.ForEachWith("*", "*", "Bob", "*", func(s, p, o, g string) {
    // ...
})

// Remove all statements about Charlie, from all graphs.
s.Remove("Charlie", "*", "Charlie", "*")
```

## Documentation

Full API documentation is on [GoDoc](https://godoc.org/github.com/jimsmart/store4)

## Tests

To run the tests, execute `go test` inside the project folder. For a full coverage report, try:
```bash
$ go test -coverprofile=coverage.out && go tool cover -html=coverage.out
```

# License

Package store4 is copyrighted by Jim Smart and released under the [MIT License](LICENSE.md)

## Additional credits

The internals of store4 draw heavily from the implementation of N3Store, a component of [N3.js](https://github.com/RubenVerborgh/N3.js). The N3.js library is copyrighted by Ruben Verborgh and released under the MIT License.


