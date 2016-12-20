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
type TripleCallbackFn func(s, p, o string)

// TripleTestFn is the function signature used to implement
// callback functions performing triple tests.
// A response of true means that the test has been passed.
//
// Used with calls to GraphView's Every, EveryWith, Some and SomeWith.
type TripleTestFn func(s, p, o string) bool

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

// GraphView returns a proxy-façade that provides a
// triple-based API for working with graphs in the store.
func (s *QuadStore) GraphView(graph string) *GraphView {
	return &GraphView{
		Graph:     graph,
		QuadStore: s,
	}
}

// NewGraph returns a GraphView over an unnamed graph with a
// newly created QuadStore as its backing.
//
// It is shorthand for NewQuadStore(args).GraphView("")
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
func (g *GraphView) Add(subject, predicate, object string) bool {
	return g.QuadStore.Add(subject, predicate, object, g.Graph)
}

// Count returns a count of triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) Count(subject, predicate, object string) uint64 {
	return g.QuadStore.Count(subject, predicate, object, g.Graph)
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
// Acting like the 'for all' quantifier in maths, it should
// be noted that Every returns true for an empty graph.
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
// Acting like the 'for all' quantifier in maths, it should
// be noted that EveryWith returns true for an empty graph.
// By extension, if the given parameters cause the iteration
// set to be empty, then EveryWith also returns true.
func (g *GraphView) EveryWith(subject, predicate, object string, fn TripleTestFn) bool {
	return g.QuadStore.EveryWith(subject, predicate, object, g.Graph, adaptTripleTestFn(fn))
}

// FindObjects returns a list of distinct object terms for all
// triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) FindObjects(subject, predicate string) []string {
	return g.QuadStore.FindObjects(subject, predicate, g.Graph)
}

// FindPredicates returns a list of distinct predicate terms for all
// triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) FindPredicates(subject, object string) []string {
	return g.QuadStore.FindPredicates(subject, object, g.Graph)
}

// FindSubjects returns a list of distinct subject terms for all
// triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) FindSubjects(predicate, object string) []string {
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
func (g *GraphView) ForEachWith(subject, predicate, object string, fn TripleCallbackFn) {
	g.QuadStore.ForEachWith(subject, predicate, object, g.Graph, adaptTripleCallbackFn(fn))
}

// ForObjects executes the given callback once for each distinct object term
// for all triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) ForObjects(subject, predicate string, fn StringCallbackFn) {
	g.QuadStore.ForObjects(subject, predicate, g.Graph, fn)
}

// ForPredicates executes the given callback once for each distinct predicate term
// for all triples in the graph that graph the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) ForPredicates(subject, object string, fn StringCallbackFn) {
	g.QuadStore.ForPredicates(subject, object, g.Graph, fn)
}

// ForSubjects executes the given callback once for each distinct subject term
// for all triples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) ForSubjects(predicate, object string, fn StringCallbackFn) {
	g.QuadStore.ForSubjects(predicate, object, g.Graph, fn)
}

// Remove quads from the underlying QuadStore,
// with the given subject, predicate and object values
// and this GraphView's Graph value.
// Returns true if quads were removed,
// or false if no matching quads exist.
//
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) Remove(subject, predicate, object string) bool {
	return g.QuadStore.Remove(subject, predicate, object, g.Graph)
}

// Size returns the total count of triples in the graph.
func (g *GraphView) Size() uint64 {
	gimpl, ok := g.QuadStore.graphs[g.Graph]
	if !ok {
		return 0
	}
	return gimpl.size
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
func (g *GraphView) SomeWith(subject, predicate, object string, fn TripleTestFn) bool {
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
			sort.Strings(objects)
			for _, object := range objects {
				buf.WriteByte('[')
				buf.WriteString(subject)
				buf.WriteByte(' ')
				buf.WriteString(predicate)
				buf.WriteByte(' ')
				buf.WriteString(object)
				buf.WriteByte(']')
				buf.WriteByte('\n')
			}
		}
	}
	return buf.String()
}

func adaptTripleCallbackFn(fn TripleCallbackFn) QuadCallbackFn {
	return func(s, p, o, g string) {
		fn(s, p, o)
	}
}

func adaptTripleTestFn(fn TripleTestFn) QuadTestFn {
	return func(s, p, o, g string) bool {
		return fn(s, p, o)
	}
}