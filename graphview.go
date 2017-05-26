package store4

import (
	"bytes"
	"fmt"
	"sort"
)

// TripleCallbackFn is the function signature used to implement
// callback functions that receive a triple.
//
// Used with calls to GraphView's ForEach and ForEachWith.
type TripleCallbackFn func(s, p string, o interface{})

// TripleTestFn is the function signature used to implement
// callback functions performing triple tests.
// A response of true means that the test has been passed.
//
// Used with calls to GraphView's Every, EveryWith, Some and SomeWith.
type TripleTestFn func(s, p string, o interface{}) bool

// GraphView provides a graph-centric API
// for working with subject-predicate-object triples.
//
// GraphView is a convenience façade that simply
// proxies calls to its associated QuadStore.
//
// Returned by calls to NewGraph and QuadStore.GraphView.
type GraphView struct {
	Graph     string
	QuadStore *QuadStore
}

// GraphView returns a GraphView for the given graph name.
func (s *QuadStore) GraphView(graph string) *GraphView {
	return &GraphView{
		Graph:     graph,
		QuadStore: s,
	}
}

// TODO(js) Consistent terminology: pattern vs parameters (here) and term vs parameters (elsewhere).
// Also: subject(etc) vs subject term, graph vs graph name.

// GraphViews returns a list of GraphViews for graphs in the store
// that contain triples that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) GraphViews(subject, predicate string, object interface{}) []*GraphView {
	var out []*GraphView
	s.ForGraphs(subject, predicate, object, func(graph string) {
		g := &GraphView{
			Graph:     graph,
			QuadStore: s,
		}
		out = append(out, g)
	})
	return out
}

// NewGraph returns a GraphView over an unnamed graph with a
// newly created QuadStore as its backing, optionally
// initialising it with triples.
//
// NewGraph is shorthand for NewQuadStore().GraphView("")
// to quickly provide a default graph to use.
//
// If you wish to create graph that is not unnamed, or create
// a graph associated with an existing store, see QuadStore.GraphView.
//
// Initial triples can be provided using the following types:
//  [][3]string
//  [3]string
//
// See NewQuadStore for a greater variety of initialisation options.
func NewGraph(args ...interface{}) *GraphView {
	g := NewQuadStore().GraphView("")
	for _, arg := range args {
		switch arg := arg.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T\n", arg))
		case [3]string:
			// Single string triple.
			g.Add(arg[0], arg[1], arg[2])
		case [][3]string:
			// Slice of string triples.
			for _, q := range arg {
				g.Add(q[0], q[1], q[2])
			}
		}
	}
	return g
}

// Add a quad to the underlying QuadStore,
// with the given subject, predicate and object values
// and this GraphView's Graph value.
// Returns true if the quad was a new quad,
// or false if the quad already existed.
//
// If any of the given terms are "*" (an asterisk),
// then this method will panic. (The asterisk is reserved
// for wildcard operations throughout the API).
func (g *GraphView) Add(subject, predicate string, object interface{}) bool {
	return g.QuadStore.Add(subject, predicate, object, g.Graph)
}

// Count returns a count of triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) Count(subject, predicate string, object interface{}) uint64 {
	return g.QuadStore.Count(subject, predicate, object, g.Graph)
}

// Empty returns true if the GraphView has no contents.
func (g *GraphView) Empty() bool {
	return g.Size() == 0
}

// Every tests whether all triples in the graph pass the test
// implemented by the given function.
//
// The given callback is
// executed once for each triple present in the graph until
// Every finds one where the callback returns false. If such
// an element is found, iteration is immediately halted and
// Every returns false. Otherwise, if the callback returns
// true for all triples, then Every returns true.
//
// If no triples match the given terms, or the graph is empty,
// then Every returns false. Note that this differs from
// the interpretation of 'every' in some other languages,
// which may return true for an empty iteration set.
func (g *GraphView) Every(fn TripleTestFn) bool {
	return g.QuadStore.EveryWith("*", "*", "*", g.Graph, adaptTripleTestFn(fn))
}

// EveryWith tests whether all triples in the graph that match the
// given terms pass the test implemented by the given function.
//
// The given callback is
// executed once for each matching triple in the graph until
// EveryWith finds one where the callback returns false. If such
// an element is found, iteration is immediately halted and
// EveryWith returns false. Otherwise, if the callback returns
// true for all triples, then EveryWith returns true.
//
// If no triples match the given terms, or the graph is empty,
// then EveryWith returns false. Note that this differs from
// the interpretation of 'every' in some other languages,
// which may return true for an empty iteration set.
func (g *GraphView) EveryWith(subject, predicate string, object interface{}, fn TripleTestFn) bool {
	return g.QuadStore.EveryWith(subject, predicate, object, g.Graph, adaptTripleTestFn(fn))
}

// FindObjects returns a list of distinct object terms for all
// triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) FindObjects(subject, predicate string) []interface{} {
	return g.QuadStore.FindObjects(subject, predicate, g.Graph)
}

// FindPredicates returns a list of distinct predicate terms for all
// triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) FindPredicates(subject string, object interface{}) []string {
	return g.QuadStore.FindPredicates(subject, object, g.Graph)
}

// FindSubjects returns a list of distinct subject terms for all
// triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) FindSubjects(predicate string, object interface{}) []string {
	return g.QuadStore.FindSubjects(predicate, object, g.Graph)
}

// ForEach executes the given callback once for each triple in the graph.
func (g *GraphView) ForEach(fn TripleCallbackFn) {
	g.QuadStore.ForEachWith("*", "*", "*", g.Graph, adaptTripleCallbackFn(fn))
}

// ForEachWith executes the given callback once for each triple in the graph
// that matches the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) ForEachWith(subject, predicate string, object interface{}, fn TripleCallbackFn) {
	g.QuadStore.ForEachWith(subject, predicate, object, g.Graph, adaptTripleCallbackFn(fn))
}

// ForObjects executes the given callback once for each distinct object term
// for all triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) ForObjects(subject, predicate string, fn ObjectCallbackFn) {
	g.QuadStore.ForObjects(subject, predicate, g.Graph, fn)
}

// ForPredicates executes the given callback once for each distinct predicate term
// for all triples in the graph that graph the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) ForPredicates(subject string, object interface{}, fn StringCallbackFn) {
	g.QuadStore.ForPredicates(subject, object, g.Graph, fn)
}

// ForSubjects executes the given callback once for each distinct subject term
// for all triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) ForSubjects(predicate string, object interface{}, fn StringCallbackFn) {
	g.QuadStore.ForSubjects(predicate, object, g.Graph, fn)
}

// Remove quads from the underlying QuadStore,
// with the given subject, predicate and object values
// and this GraphView's Graph value.
// Returns the number of quads that were removed.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) Remove(subject, predicate string, object interface{}) uint64 {
	return g.QuadStore.Remove(subject, predicate, object, g.Graph)
}

// Size returns the total count of triples in the graph.
func (g *GraphView) Size() uint64 {
	return g.QuadStore.Count("*", "*", "*", g.Graph)
}

// Some tests whether some triple in the graph passes the test
// implemented by the given function.
//
// The given callback is
// executed once for each triple present in the graph until
// Some finds one where the callback returns true. If such
// an element is found, iteration is immediately halted and
// Some returns true. Otherwise, if the callback returns
// false for all triples, then Some returns false.
func (g *GraphView) Some(fn TripleTestFn) bool {
	return g.QuadStore.SomeWith("*", "*", "*", g.Graph, adaptTripleTestFn(fn))
}

// SomeWith tests whether some triple matching the given pattern
// passes the test implemented by the given function.
//
// The given callback is
// executed once for each triple matching the given pattern until
// SomeWith finds one where the callback returns true. If such
// an element is found, iteration is immediately halted and
// SomeWith returns true. Otherwise, if the callback returns
// false for all triples, then SomeWith returns false.
func (g *GraphView) SomeWith(subject, predicate string, object interface{}, fn TripleTestFn) bool {
	return g.QuadStore.SomeWith(subject, predicate, object, g.Graph, adaptTripleTestFn(fn))
}

// String returns the contents of the graph in a human-readable format.
func (g *GraphView) String() string {
	var buf bytes.Buffer
	name := g.Graph
	if len(name) > 0 {
		buf.WriteString(name)
		buf.WriteByte('\n')
	}
	subjects := g.FindSubjects("*", "*")
	sort.Strings(subjects)
	for _, subject := range subjects {
		predicates := g.FindPredicates(subject, "*")
		sort.Strings(predicates)
		for _, predicate := range predicates {
			objects := g.FindObjects(subject, predicate)
			sortObjects(objects)
			for _, object := range objects {
				buf.WriteByte('[')
				buf.WriteString(subject)
				buf.WriteByte(' ')
				buf.WriteString(predicate)
				buf.WriteByte(' ')
				buf.WriteString(fmt.Sprint(object))
				buf.WriteByte(']')
				buf.WriteByte('\n')
			}
		}
	}
	return buf.String()
}

func adaptTripleCallbackFn(fn TripleCallbackFn) QuadCallbackFn {
	return func(s, p string, o interface{}, g string) {
		fn(s, p, o)
	}
}

func adaptTripleTestFn(fn TripleTestFn) QuadTestFn {
	return func(s, p string, o interface{}, g string) bool {
		return fn(s, p, o)
	}
}
