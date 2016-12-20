package store4

import (
	"bytes"
	"fmt"
	"sort"
)

// TupleCallbackFn is the function signature used to implement
// callback functions that receive a tuple.
//
// Used with calls to Projection's ForEach and ForEachWith.
type TupleCallbackFn func(p, o string)

// TupleTestFn is the function signature used to implement
// callback functions performing tuple tests.
// A response of true means that the test has been passed.
//
// Used with calls to Projection's Every, EveryWith, Some and SomeWith.
type TupleTestFn func(p, o string) bool

// Projection provides an 'object-oriented' view of a particular subject that is expressed in the store.
//
// Returned from calls to Query, Projection and Projections on both QuadStore and Graph.
type Projection struct {
	Subject   string
	Graph     string
	QuadStore *QuadStore
}

func (s *QuadStore) Projection(subject, graph string) (*Projection, bool) {
	if subject == "*" {
		panic("Unexpected use of wildcard '*' for subject")
	}

	haltFn := func(s, p, o, g string) bool {
		return true
	}

	ok := s.SomeWith(subject, "*", "*", graph, haltFn)

	p := &Projection{
		Subject:   subject,
		Graph:     graph,
		QuadStore: s,
	}
	return p, ok
}

func (s *QuadStore) Projections(predicate, object, graph string) []*Projection {
	var out []*Projection
	s.ForSubjects(predicate, object, graph, func(subject string) {
		p := &Projection{
			Subject:   subject,
			Graph:     graph,
			QuadStore: s,
		}
		out = append(out, p)
	})
	return out
}

func (g *Graph) Projection(subject string) (*Projection, bool) {
	return g.QuadStore.Projection(subject, g.Name)
}

func (g *Graph) Projections(predicate, object string) []*Projection {
	return g.QuadStore.Projections(predicate, object, g.Name)
}

// Query retrieves a list of all Projections in the store that match the given pattern and graph.
//
// Pattern is a collection of predicate-object tuples, expressed as any of the following types:
// map[string]string, map[string][]string, [][2]string, or a single [2]string.
func (s *QuadStore) Query(pattern interface{}, graph string) []*Projection {
	// Convert given pattern into query list.
	var poList [][2]string
	switch pattern := pattern.(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T\n", pattern))
	case map[string]string:
		for p, o := range pattern {
			poList = append(poList, [2]string{p, o})
		}
	case map[string][]string:
		for p, objects := range pattern {
			for _, o := range objects {
				poList = append(poList, [2]string{p, o})
			}
		}
	case [][2]string:
		poList = pattern
	case [2]string:
		poList = append(poList, pattern)
	}
	// Nothing to do?
	if len(poList) == 0 {
		return nil
	}

	haltFn := func(s, p, o, g string) bool {
		return true
	}

	var out []*Projection
	s.ForSubjects(poList[0][0], poList[0][1], graph, func(subject string) {
		// Got a match for first q entry,
		// check if it also satisfies all other q entries.
		for _, po := range poList[1:] {
			if !s.SomeWith(subject, po[0], po[1], graph, haltFn) {
				// No match, try next subject.
				return
			}
		}
		// Matches all po tuples, add to results list.
		p := &Projection{
			Subject:   subject,
			Graph:     graph,
			QuadStore: s,
		}
		out = append(out, p)
	})
	return out
}

// Query retrieves a list of all Projections in the graph that match the given pattern.
func (g *Graph) Query(pattern interface{}) []*Projection {
	return g.QuadStore.Query(pattern, g.Name)
}

// Map returns a map containing the predicate terms for
// the projection, mapped to their corresponding object terms.
func (p *Projection) Map() map[string][]string {
	m := make(map[string][]string)
	p.ForPredicates("*", func(predicate string) {
		m[predicate] = p.FindObjects(predicate)
	})
	return m
}

// Add a tuple to the projection.
// Returns true if the tuple was a new tuple,
// or false if the tuple already existed.
//
// If any of the given terms are "*" (an asterisk),
// then this method will panic. (The asterisk is reserved
// for wildcard operations throughout the API).
func (p *Projection) Add(predicate, object string) bool {
	// TODO(js) Would it be better if Add() added to default graph if the projection's graph is set to "*" ?
	return p.QuadStore.Add(p.Subject, predicate, object, p.Graph)
}

// Count returns a count of tuples in the projection that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (p *Projection) Count(predicate, object string) uint64 {
	return p.QuadStore.Count(p.Subject, predicate, object, p.Graph)
}

// Every tests whether all tuples in the projection pass the test
// implemented by the given function.
//
// The given callback is
// executed once for each tuple present in the projection until
// Every finds one where the callback returns false. If such
// an element is found, iteration is immediately halted and
// Every returns false. Otherwise, if the callback returns
// true for all tuples, then Every returns true.
//
// Acting like the 'for all' quantifier in maths, it should
// be noted that Every returns true for an empty store.
func (p *Projection) Every(fn TupleTestFn) bool {
	return p.QuadStore.EveryWith(p.Subject, "*", "*", p.Graph, adaptTupleTestFn(fn))
}

// EveryWith tests whether all tuples in the projection that match the
// given terms pass the test implemented by the given function.
//
// The given callback is
// executed once for each matching tuple in the projection until
// EveryWith finds one where the callback returns false. If such
// an element is found, iteration is immediately halted and
// EveryWith returns false. Otherwise, if the callback returns
// true for all tuples, then EveryWith returns true.
//
// Acting like the 'for all' quantifier in maths, it should
// be noted that EveryWith returns true for an empty projection.
// By extension, if the given parameters cause the iteration
// set to be empty, then EveryWith also returns true.
func (p *Projection) EveryWith(predicate, object string, fn TupleTestFn) bool {
	return p.QuadStore.EveryWith(p.Subject, predicate, object, p.Graph, adaptTupleTestFn(fn))
}

// FindObjects returns a list of distinct object terms for all tuples in the projection that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (p *Projection) FindObjects(predicate string) []string {
	return p.QuadStore.FindObjects(p.Subject, predicate, p.Graph)
}

// FindPredicates returns a list of distinct predicate terms for all tuples in the projection that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (p *Projection) FindPredicates(object string) []string {
	return p.QuadStore.FindPredicates(p.Subject, object, p.Graph)
}

// ForEach executes the given callback once for each tuple in the graph.
func (p *Projection) ForEach(fn TupleCallbackFn) {
	p.QuadStore.ForEachWith(p.Subject, "*", "*", p.Graph, adaptTupleCallbackFn(fn))
}

// ForEachWith executes the given callback once for each tuple in the graph
// that matches the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (p *Projection) ForEachWith(predicate, object string, fn TupleCallbackFn) {
	p.QuadStore.ForEachWith(p.Subject, predicate, object, p.Graph, adaptTupleCallbackFn(fn))
}

// ForObjects executes the given callback once for each distinct object term
// for all tuples in the graph that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (p *Projection) ForObjects(predicate string, fn StringCallbackFn) {
	p.QuadStore.ForObjects(p.Subject, predicate, p.Graph, fn)
}

// ForPredicates executes the given callback once for each distinct predicate term
// for all tuples in the graph that graph the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (p *Projection) ForPredicates(object string, fn StringCallbackFn) {
	p.QuadStore.ForPredicates(p.Subject, object, p.Graph, fn)
}

// Removes tuples from the projection. Returns true if tuples were removed,
// or false if no matching tuples exist.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (p *Projection) Remove(predicate, object string) bool {
	return p.QuadStore.Remove(p.Subject, predicate, object, p.Graph)
}

// Size returns the total count of tuples in the projection.
func (p *Projection) Size() uint64 {
	return p.QuadStore.Count(p.Subject, "*", "*", p.Graph)
}

// Some tests whether some tuple in the projection passes the test
// implemented by the given function.
//
// The given callback is
// executed once for each tuple present in the graph until
// Some finds one where the callback returns true. If such
// an element is found, iteration is immediately halted and
// Some returns true. Otherwise, if the callback returns
// false for all tuples, then Some returns false.
func (p *Projection) Some(fn TupleTestFn) bool {
	return p.QuadStore.SomeWith(p.Subject, "*", "*", p.Graph, adaptTupleTestFn(fn))
}

// SomeWith tests whether some tuple matching the given pattern
// passes the test implemented by the given function.
//
// The given callback is
// executed once for each tuple matching the given pattern until
// SomeWith finds one where the callback returns true. If such
// an element is found, iteration is immediately halted and
// SomeWith returns true. Otherwise, if the callback returns
// false for all tuples, then SomeWith returns false.
func (p *Projection) SomeWith(predicate, object string, fn TupleTestFn) bool {
	return p.QuadStore.SomeWith(p.Subject, predicate, object, p.Graph, adaptTupleTestFn(fn))
}

// String returns the contents of the projection in a human-readable format.
func (p *Projection) String() string {
	var buf bytes.Buffer
	graph := p.Graph
	if len(graph) > 0 {
		buf.WriteString(graph)
		buf.WriteByte('\n')
	}
	buf.WriteString(p.Subject)
	buf.WriteByte('\n')
	predicates := p.FindPredicates("*")
	sort.Strings(predicates)
	for _, predicate := range predicates {
		objects := p.FindObjects(predicate)
		sort.Strings(objects)
		for _, object := range objects {
			buf.WriteByte('[')
			buf.WriteString(predicate)
			buf.WriteByte(' ')
			buf.WriteString(object)
			buf.WriteByte(']')
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func adaptTupleCallbackFn(fn TupleCallbackFn) QuadCallbackFn {
	return func(s, p, o, g string) {
		fn(p, o)
	}
}

func adaptTupleTestFn(fn TupleTestFn) QuadTestFn {
	return func(s, p, o, g string) bool {
		return fn(p, o)
	}
}
