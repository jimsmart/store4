// Package store4 provides a fast in-memory quad store, with graph and subject views.
//
// QuadStore API
//
// The QuadStore API is based around subject-predicate-object-graph quads.
//
//  // Add some quads to a store.
//  s := store4.NewQuadStore()
//  s.Add("s", "p", "o", "g")
//  s.Add("Alice", "knows", "Bob", "")
//  s.Add("Alice", "knows", "Charlie", "")
//  s.Add("Charlie", "knows", "Bob", "")
//
//  // Find everyone that Alice knows, in any graph.
//  x := s.FindObjects("Alice", "knows", "*")
//
//  // Find everyone who knows Charlie, in the 'unnamed' graph.
//  y := s.FindSubjects("knows", "Charlie", "")
//
//  // Remove all statements about Charlie, from all graphs.
//  s.Remove("Charlie", "*", "Charlie", "*")
//
// Callbacks make it easy to work with and query the contents of the
// store without allocating lists for results.
//
//  // Iterate over all quads.
//  s.ForEach(func(s, p string, o interface{}, g string) {
//      // ...
//  })
//
//  // Iterate over quads matching given pattern.
//  s.ForEachWith("*", "*", "Bob", "*", func(s, p string, o interface{}, g string) {
//      // ...
//  })
//
// For cancellable iterators see Some and Every, and
// their filtering counterparts SomeWith and EveryWith.
//
// QuadStore also features callback hooks for both OnAdd and OnRemove,
// which can be used to integrate external features such as logging or
// inference.
//
// GraphView API
//
// The GraphView API is based around subject-predicate-object triples.
//
//  // Get a view over the 'unnamed' graph.
//  g := s.GraphView("")
//
//  // Find everyone that Alice knows in this graph view.
//  x := g.FindObjects("Alice", "knows")
//
//  // Add a triple to the graph (adds a quad to the store).
//  g.Add("Charlie", "knows", "Bob")
//
// Obtain a GraphView by calling QuadStore.GraphView, QuadStore.GraphViews
// — or NewGraph.
//
// SubjectView API
//
// The SubjectView API is based around predicate-object (property/value) tuples.
//
//  // Get a view over the subject Alice in 'unnamed' graph.
//  v := s.SubjectView("Alice", "")
//
//  // Find all predicates for this subject view.
//  props := v.FindPredicates("*")
//
//  // Add a tuple to this view (adds a quad to the store).
//  v.Add("knows", "Eve")
//
//  // Get a map of all predicate-object mappings for this subject view.
//  m := v.Map()
//
// SubjectViews are returned by calls to Query, SubjectView and SubjectViews.
//
// Implementation
//
// Inside QuadStore each graph is indexed by SPO, POS and OSP,
// with each index being composed of three layers of native Go maps.
//
// Internally, the store uses numeric identifiers for index keys,
// and only holds a single reference to each term.
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
