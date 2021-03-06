package store4

import (
	"bytes"
	"fmt"
	"sort"
)

// TupleCallbackFn is the function signature used to implement
// callback functions that receive a tuple.
//
// Used with calls to SubjectView's ForEach and ForEachWith.
type TupleCallbackFn func(p string, o interface{})

// TupleTestFn is the function signature used to implement
// callback functions performing tuple tests.
// A response of true means that the test has been passed.
//
// Used with calls to SubjectView's Every, EveryWith, Some and SomeWith.
type TupleTestFn func(p string, o interface{}) bool

// SubjectView provides a subject-centric API
// for working with predicate-object (property/value) tuples.
//
// SubjectView is a convenience façade that simply
// proxies calls to its associated QuadStore.
//
// Returned by calls to Query, SubjectView and
// SubjectViews on both QuadStore and GraphView.
type SubjectView struct {
	Subject   string
	Graph     string
	QuadStore *QuadStore
}

// SubjectView returns a SubjectView for the given subject and graph.
func (s *QuadStore) SubjectView(subject, graph string) *SubjectView {
	p := &SubjectView{
		Subject:   subject,
		Graph:     graph,
		QuadStore: s,
	}
	return p
}

// SubjectViews returns a list of SubjectViews for subjects that
// match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) SubjectViews(predicate string, object interface{}, graph string) []*SubjectView {
	var out []*SubjectView
	s.ForSubjects(predicate, object, graph, func(subject string) {
		p := &SubjectView{
			Subject:   subject,
			Graph:     graph,
			QuadStore: s,
		}
		out = append(out, p)
	})
	return out
}

// SubjectView returns a SubjectView for the given subject.
func (g *GraphView) SubjectView(subject string) *SubjectView {
	return g.QuadStore.SubjectView(subject, g.Graph)
}

// SubjectViews returns a list of SubjectViews for subjects that
// match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (g *GraphView) SubjectViews(predicate string, object interface{}) []*SubjectView {
	return g.QuadStore.SubjectViews(predicate, object, g.Graph)
}

type tuple struct {
	p string
	o interface{}
}

// Query returns a list of SubjectViews for subjects in the store
// having predicate-object terms that match the given pattern.
//
// Pattern is a collection of predicate-object tuples,
// expressed using any of the following types:
//  map[string][]interface{}
//  map[string][]string
//  map[string]interface{}
//  map[string]string
//  [][2]string
//  [2]string
func (s *QuadStore) Query(pattern interface{}, graph string) []*SubjectView {
	poList := predicateObjectList(pattern)
	// Nothing to do?
	if len(poList) == 0 {
		return nil
	}

	haltFn := func(s, p string, o interface{}, g string) bool {
		return true
	}

	var out []*SubjectView
	s.ForSubjects(poList[0].p, poList[0].o, graph, func(subject string) {
		// Got a match for first q entry,
		// check if it also satisfies all other q entries.
		for _, po := range poList[1:] {
			if !s.SomeWith(subject, po.p, po.o, graph, haltFn) {
				// No match, try next subject.
				return
			}
		}
		// Matches all po tuples, add to results list.
		p := &SubjectView{
			Subject:   subject,
			Graph:     graph,
			QuadStore: s,
		}
		out = append(out, p)
	})
	return out
}

func predicateObjectList(pattern interface{}) []tuple {
	// Convert given pattern into query list.
	var poList []tuple
	switch pattern := pattern.(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T\n", pattern))
	case map[string]string:
		for p, o := range pattern {
			poList = append(poList, tuple{p, o})
		}
	case map[string]interface{}:
		for p, o := range pattern {
			poList = append(poList, tuple{p, o})
		}
	case map[string][]string:
		for p, objects := range pattern {
			for _, o := range objects {
				poList = append(poList, tuple{p, o})
			}
		}
	case map[string][]interface{}:
		for p, objects := range pattern {
			for _, o := range objects {
				poList = append(poList, tuple{p, o})
			}
		}
	case [][2]string:
		for _, po := range pattern {
			poList = append(poList, tuple{po[0], po[1]})
		}
	case [2]string:
		poList = append(poList, tuple{pattern[0], pattern[1]})
	}
	return poList
}

// Query returns a list of SubjectViews for subjects in the graph
// having predicate-object terms that match the given pattern.
//
// Pattern is a collection of predicate-object tuples,
// expressed using any of the following types:
//  map[string][]interface{}
//  map[string][]string
//  map[string]interface{}
//  map[string]string
//  [][2]string
//  [2]string
func (g *GraphView) Query(pattern interface{}) []*SubjectView {
	return g.QuadStore.Query(pattern, g.Graph)
}

// Map returns a 'property/value' map containing the predicate terms for
// the SubjectView's subject, mapped to their corresponding object terms.
func (v *SubjectView) Map() map[string][]interface{} {
	m := make(map[string][]interface{})
	v.ForPredicates("*", func(predicate string) {
		m[predicate] = v.FindObjects(predicate)
	})
	return m
}

// Add a quad to the underlying QuadStore,
// with the given predicate and object values
// and this SubjectView's Subject and Graph values.
// Returns true if the quad was a new quad,
// or false if the quad already existed.
//
// If any of the given terms are "*" (an asterisk),
// then this method will panic. (The asterisk is reserved
// for wildcard operations throughout the API).
func (v *SubjectView) Add(predicate string, object interface{}) bool {
	return v.QuadStore.Add(v.Subject, predicate, object, v.Graph)
}

// Count returns a count of tuples in the SubjectView that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (v *SubjectView) Count(predicate string, object interface{}) uint64 {
	return v.QuadStore.Count(v.Subject, predicate, object, v.Graph)
}

// Empty returns true if the SubjectView has no contents.
func (v *SubjectView) Empty() bool {
	haltFn := func(s, p string, o interface{}, g string) bool {
		return true
	}
	return !v.QuadStore.SomeWith(v.Subject, "*", "*", v.Graph, haltFn)
}

// Every tests whether all tuples in the SubjectView pass the test
// implemented by the given function.
//
// The given callback is
// executed once for each tuple present in the SubjectView until
// Every finds one where the callback returns false. If such
// an element is found, iteration is immediately halted and
// Every returns false. Otherwise, if the callback returns
// true for all tuples, then Every returns true.
//
// If no tuples match the given terms, or the SubjectView is empty,
// then Every returns false. Note that this differs from
// the interpretation of 'every' in some other languages,
// which may return true for an empty iteration set.
func (v *SubjectView) Every(fn TupleTestFn) bool {
	return v.QuadStore.EveryWith(v.Subject, "*", "*", v.Graph, adaptTupleTestFn(fn))
}

// EveryWith tests whether all tuples in the SubjectView that match the
// given terms pass the test implemented by the given function.
//
// The given callback is
// executed once for each matching tuple in the SubjectView until
// EveryWith finds one where the callback returns false. If such
// an element is found, iteration is immediately halted and
// EveryWith returns false. Otherwise, if the callback returns
// true for all tuples, then EveryWith returns true.
//
// If no tuples match the given terms, or the SubjectView is empty,
// then EveryWith returns false. Note that this differs from
// the interpretation of 'every' in some other languages,
// which may return true for an empty iteration set.
func (v *SubjectView) EveryWith(predicate string, object interface{}, fn TupleTestFn) bool {
	return v.QuadStore.EveryWith(v.Subject, predicate, object, v.Graph, adaptTupleTestFn(fn))
}

// FindObjects returns a list of distinct object terms for all
// tuples in the SubjectView that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (v *SubjectView) FindObjects(predicate string) []interface{} {
	return v.QuadStore.FindObjects(v.Subject, predicate, v.Graph)
}

// FindPredicates returns a list of distinct predicate terms for all
// tuples in the SubjectView that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (v *SubjectView) FindPredicates(object interface{}) []string {
	return v.QuadStore.FindPredicates(v.Subject, object, v.Graph)
}

// ForEach executes the given callback once for each tuple in the SubjectView.
func (v *SubjectView) ForEach(fn TupleCallbackFn) {
	v.QuadStore.ForEachWith(v.Subject, "*", "*", v.Graph, adaptTupleCallbackFn(fn))
}

// ForEachWith executes the given callback once for each tuple in the SubjectView
// that matches the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (v *SubjectView) ForEachWith(predicate string, object interface{}, fn TupleCallbackFn) {
	v.QuadStore.ForEachWith(v.Subject, predicate, object, v.Graph, adaptTupleCallbackFn(fn))
}

// ForObjects executes the given callback once for each distinct object term
// for all tuples in the SubjectView that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (v *SubjectView) ForObjects(predicate string, fn ObjectCallbackFn) {
	v.QuadStore.ForObjects(v.Subject, predicate, v.Graph, fn)
}

// ForPredicates executes the given callback once for each distinct predicate term
// for all tuples in the SubjectView that graph the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (v *SubjectView) ForPredicates(object interface{}, fn StringCallbackFn) {
	v.QuadStore.ForPredicates(v.Subject, object, v.Graph, fn)
}

// Remove quads from the underlying QuadStore,
// with the given predicate and object values
// and this SubjectView's Subject and Graph values.
// Returns the number of quads removed.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (v *SubjectView) Remove(predicate string, object interface{}) uint64 {
	return v.QuadStore.Remove(v.Subject, predicate, object, v.Graph)
}

// Size returns the total count of tuples in the SubjectView.
func (v *SubjectView) Size() uint64 {
	return v.QuadStore.Count(v.Subject, "*", "*", v.Graph)
}

// Some tests whether some tuple in the SubjectView passes the test
// implemented by the given function.
//
// The given callback is
// executed once for each tuple present in the SubjectView until
// Some finds one where the callback returns true. If such
// an element is found, iteration is immediately halted and
// Some returns true. Otherwise, if the callback returns
// false for all tuples, then Some returns false.
func (v *SubjectView) Some(fn TupleTestFn) bool {
	return v.QuadStore.SomeWith(v.Subject, "*", "*", v.Graph, adaptTupleTestFn(fn))
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
func (v *SubjectView) SomeWith(predicate string, object interface{}, fn TupleTestFn) bool {
	return v.QuadStore.SomeWith(v.Subject, predicate, object, v.Graph, adaptTupleTestFn(fn))
}

// String returns the contents of the SubjectView in a human-readable format.
func (v *SubjectView) String() string {
	var buf bytes.Buffer
	graph := v.Graph
	if len(graph) > 0 {
		buf.WriteString(graph)
		buf.WriteByte('\n')
	}
	buf.WriteString(v.Subject)
	buf.WriteByte('\n')
	predicates := v.FindPredicates("*")
	sort.Strings(predicates)
	for _, predicate := range predicates {
		objects := v.FindObjects(predicate)
		sortObjects(objects)
		for _, object := range objects {
			buf.WriteByte('[')
			buf.WriteString(predicate)
			buf.WriteByte(' ')
			buf.WriteString(fmt.Sprint(object))
			buf.WriteByte(']')
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func adaptTupleCallbackFn(fn TupleCallbackFn) QuadCallbackFn {
	return func(s, p string, o interface{}, g string) {
		fn(p, o)
	}
}

func adaptTupleTestFn(fn TupleTestFn) QuadTestFn {
	return func(s, p string, o interface{}, g string) bool {
		return fn(p, o)
	}
}
