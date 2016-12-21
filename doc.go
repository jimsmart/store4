// Package store4 provides a fast in-memory string-based quad store
// (graphs of triples).
//
// Simple
//
// Package store4 does not feature any specific support
// for RDF — its API is purely string-based, bare-boned and minimal.
//
//  // Add some quads to a store.
//  s := store4.NewQuadStore()
//  s.Add("s", "p", "o", "g")
//  s.Add("Alice", "knows", "Bob", "")
//  s.Add("Alice", "knows", "Charlie", "")
//  s.Add("Charlie", "knows", "Bob", "")
//  s.Add("Bob", "isa", "Cat", "")
//  s.Add("play", "with", "strings!", "₍˄·͈༝·͈˄₎◞ ̑̑ෆ⃛")
//  s.Add("_:bn", "<urn:foo:bar>", `"baz"`, "<qux>")
//
//  // Find everyone that Alice knows, in any graph.
//  list := s.FindObjects("Alice", "knows", "*")
//  fmt.Println(list)
//  // Output (exact order may vary):
//  // [Bob Charlie]
//
//  // Find everyone who knows Charlie, in the 'unnamed' graph.
//  x := s.FindSubjects("knows", "Charlie", "")
//  fmt.Println(x)
//  // Output:
//  // [Alice]
//
//  // Remove all statements about Charlie, from all graphs.
//  s.Remove("Charlie", "*", "Charlie", "*")
//
// Callbacks make it easy to work with and query the contents of the
// store without allocating lists for results.
//
//  // Iterate over all quads.
//  s.ForEach(func(s, p, o, g string) {
//      // ...
//  })
//
//  // Iterate over quads matching given pattern.
//  s.ForEachWith("*", "*", "Bob", "*", func(s, p, o, g string) {
//      fmt.Println(s, p, o, g)
//  })
//  // Output (exact order may vary):
//  // Alice knows Bob
//  // Charlie knows Bob
//
// For cancellable iterators see Some and Every, and
// their filtering counterparts SomeWith and EveryWith.
//
// Internals
//
// Each graph is indexed by SPO, POS and OSP.
// Each index is composed of three layers of native Go maps.
//
// Additionally, the store only holds a single reference to each
// string term it has seen.
//
// Concurrency
//
// Package store4 is not concurrency safe while being modified.
//
// Dependencies
//
// Standard library.
//
// Ginkgo and Gomega to run tests.
//
// License
//
// Package store4 is free software / open source software, released
// under the MIT License.
//
// Additional credits
//
// The internals of QuadStore draw heavily from the implementation
// of N3Store, a component of N3.js.
//
// The N3.js library is
// copyrighted by Ruben Verborgh and released under the MIT License.
// https://github.com/RubenVerborgh/N3.js
package store4
