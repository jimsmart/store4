package store4

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

// QuadCallbackFn is the function signature used to implement
// callback functions that receive a quad.
//
// Used with calls to QuadStore's ForEach and ForEachWith.
// type QuadCallbackFn func(s, p, o, g string)
type QuadCallbackFn func(s, p string, o interface{}, g string)

// QuadTestFn is the function signature used to implement
// callback functions performing quad tests.
// A response of true means that the test has been passed.
//
// Used with calls to QuadStore's Every, EveryWith, Some and SomeWith.
// type QuadTestFn func(s, p, o, g string) bool
type QuadTestFn func(s, p string, o interface{}, g string) bool

// StringCallbackFn is the function signature used to implement
// callback functions that receive a string.
//
// Used with calls to FindSubjects, FindPredicates and FindGraphs.
type StringCallbackFn func(s string)

// ObjectCallbackFn is the function signature used to implement
// callback functions that receive an object.
//
// Used with calls to FindObjects.
type ObjectCallbackFn func(o interface{})

// QuadStore is an in-memory string-based quad store.
//
// If provided, the OnAdd callback will be called once for every quad
// added to the store with the Add method. It is called once per quad,
// after each quad has successfully been added to the store. It is not
// called if the added quad already existed in the store. Note that
// when the callback is invoked, the store size will already have been
// incremented, and all internal indexes will be in a
// consistent state — so it is safe to add further quads (or remove
// quads) from within the callback, should one wish to do so.
//
// Likewise: if provided, the OnRemove callback will be called once for every quad
// removed from the store with the Remove method. It is called once
// per quad, after each quad has been successfully removed from the store.
// Note that when the callback is invoked, the store size will already
// have been decremented, and all internal indexes will be
// in a consistent state — so it is safe to remove further quads (or add
// quads) from within the callback, should one wish to do so.
type QuadStore struct {
	// OnAdd is called whenever a new quad is added to the store
	// (post-addition).
	OnAdd QuadCallbackFn
	// OnRemove is called whenever a quad is removed from the store
	// (post-removal).
	OnRemove QuadCallbackFn
	// size is a count of quads in the store.
	size uint64
	// graphs hold the store's graphs.
	graphs graphMap

	pool *pool
}

// graphMap is a map holding the store's graphs, keyed by graph name.
type graphMap map[string]*indexedGraph

// indexedGraph represents a graph of triples,
// held only in the indexes, which are indexed
// three ways: SPO, POS and OSP.
type indexedGraph struct {
	size     uint64
	spoIndex indexRoot
	posIndex indexRoot
	ospIndex indexRoot
}

// index is map-based index consisting of three layers.
type indexRoot map[uint64]indexMid
type indexMid map[uint64]indexLeaf
type indexLeaf map[uint64]struct{}

// NewQuadStore creates a new quad store,
// optionally initialising it with quads or triples.
//
// Initial quads or triples can be provided using any
// of the following types:
//  [][4]string
//  [][3]string
//  [4]string
//  [3]string
//
// Alternatively, initial quads or triples may be provided
// using a slice of structs (or a single instance of a struct)
// having the following fields:
//  struct {
//      Subject   string
//      Predicate string
//      Object    string // Or interface{}
//      Graph     string // Optional.
//  }
// -or-
//  struct {
//      S string
//      P string
//      O string // Or interface{}
//      G string // Optional.
//  }
//
// Furthermore, initial quads or triples can also be
// provided using a slice of types (or a single instance of
// a type) that implement either:
//  interface {
//	    Subject()   string
//	    Predicate() string
//	    Object()    interface{}
//	    Graph()     string // Optional.
//  }
// -or-
//  interface {
//      S() string
//      P() string
//      O() interface{}
//      G() string // Optional.
//  }
//
// Finally, if the type of the given args cannot be handled,
// then NewQuadStore will panic.
func NewQuadStore(args ...interface{}) *QuadStore {
	s := &QuadStore{
		graphs: make(map[string]*indexedGraph),
		pool:   newPool(),
	}
	// Initialise store with any given data.
	for _, arg := range args {
		switch arg := arg.(type) {
		default:
			if !addQuadFromInterfaces(s, arg) {
				initWithReflection(s, arg)
			}
		case [4]string:
			// Single string quad.
			s.Add(arg[0], arg[1], arg[2], arg[3])
		case [3]string:
			// Single string triple.
			s.Add(arg[0], arg[1], arg[2], "")
		case [][4]string:
			// Slice of string quads.
			for _, q := range arg {
				s.Add(q[0], q[1], q[2], q[3])
			}
		case [][3]string:
			// Slice of string triples.
			for _, q := range arg {
				s.Add(q[0], q[1], q[2], "")
			}
		}
	}
	return s
}

type tripler interface {
	Subject() string
	Predicate() string
	Object() interface{}
}

type grapher interface {
	Graph() string
}

type simplerTripler interface {
	S() string
	P() string
	O() interface{}
}

type simplerGrapher interface {
	G() string
}

func addQuadFromInterfaces(s *QuadStore, arg interface{}) bool {
	gr := ""
	if t, ok := arg.(tripler); ok {
		if g, ok := arg.(grapher); ok {
			gr = g.Graph()
		}
		s.Add(t.Subject(), t.Predicate(), t.Object(), gr)
		return true
	}
	if t, ok := arg.(simplerTripler); ok {
		if g, ok := arg.(simplerGrapher); ok {
			gr = g.G()
		}
		s.Add(t.S(), t.P(), t.O(), gr)
		return true
	}
	return false
}

func initWithReflection(s *QuadStore, arg interface{}) {
	m := make([]int, 4)
	t := reflect.TypeOf(arg)
	k := t.Kind()
	switch k {
	case reflect.Ptr:
		// Pointer to single quad- or triple-like struct.
		t = reflect.Indirect(reflect.ValueOf(arg)).Type()
		if findMappings(t, m) {
			s.Add(quadFromStruct(arg, m))
			return
		}
	case reflect.Struct:
		// Single quad-like or triple-like struct.
		if findMappings(t, m) {
			s.Add(quadFromStruct(arg, m))
			return
		}
	case reflect.Slice:
		// A slice of quad- or triple-like structs,
		// or slice of pointers to quad- or triple-like structs,
		// or, alternatively, a slice structs implementing a
		// compatible interface.
		val := reflect.ValueOf(arg)
		length := val.Len()
		if length == 0 {
			return
		}
		el := val.Index(0)
		iface := el.Interface()
		t := reflect.Indirect(el).Type()
		if addQuadFromInterfaces(s, iface) {
			for i := 1; i < length; i++ {
				el = val.Index(i)
				iface = el.Interface()
				addQuadFromInterfaces(s, iface)
			}
			return
		}
		if findMappings(t, m) {
			s.Add(quadFromStruct(iface, m))
			for i := 1; i < length; i++ {
				el = val.Index(i)
				iface = el.Interface()
				s.Add(quadFromStruct(iface, m))
			}
			return
		}
	}
	panic(fmt.Sprintf("unexpected type %T\n", arg))
}

func quadFromStruct(arg interface{}, m []int) (s, p string, o interface{}, g string) {
	// Use the mappings to pull out the required values.
	val := reflect.Indirect(reflect.ValueOf(arg))
	s = val.Field(m[0]).String()
	p = val.Field(m[1]).String()
	o = val.Field(m[2]).Interface()
	if m[3] != -1 {
		g = val.Field(m[3]).String()
	}
	return
}

func findMappings(t reflect.Type, m []int) bool {
	// Try to build a set of mappings:
	// from quad slot to field index.
	// m := [4]int{}

	matchFields := func(fields map[string]int) bool {
		m[0], m[1], m[2], m[3] = -1, -1, -1, -1
		// Look for fields with the same names as in our map.
		for i := 0; i < t.NumField(); i++ {
			vField := t.Field(i)
			name := vField.Name
			idx, ok := fields[name]
			if ok {
				k := vField.Type.Kind()
				if k == reflect.String || (k == reflect.Interface && idx == 2) {
					m[idx] = i
				}
			}
		}
		// Check that all the fields exist.
		// Graph name - m[3] - is optional.
		return m[0] != -1 && m[1] != -1 && m[2] != -1
	}

	f1 := map[string]int{
		"Subject":   0,
		"Predicate": 1,
		"Object":    2,
		"Graph":     3,
	}

	f2 := map[string]int{
		"S": 0,
		"P": 1,
		"O": 2,
		"G": 3,
	}

	if !matchFields(f1) {
		if !matchFields(f2) {
			return false
		}
	}

	return true
}

// Size returns the total count of quads in the store.
func (s *QuadStore) Size() uint64 {
	return s.size
}

// Empty returns true if the store has no contents.
func (s *QuadStore) Empty() bool {
	return s.size == 0
}

// Add a quad to the store. Returns true if the quad was a new quad,
// or false if the quad already existed.
//
// If any of the given terms are "*" (an asterisk), then this method will panic.
// (The asterisk is reserved for wildcard operations throughout the API).
func (s *QuadStore) Add(subject, predicate string, object interface{}, graph string) bool {
	// Disallow wildcard terms
	// Optimisation: we check the other params after resolvng to IDs.
	if graph == "*" {
		panic("Unexpected use of wildcard '*' for term")
	}
	// Find the graph.
	g, ok := s.graphs[graph]
	// Create the graph if it doesn't exist yet.
	if !ok {
		g = &indexedGraph{
			spoIndex: make(indexRoot),
			posIndex: make(indexRoot),
			ospIndex: make(indexRoot),
		}
		s.graphs[graph] = g
	}
	// Get internal IDs for each term.
	sid := s.pool.getOrCreateIDString(subject)
	pid := s.pool.getOrCreateIDString(predicate)
	oid := s.pool.getOrCreateIDAny(object)
	// Disallow wildcard terms.
	// Optimisation: the fast path (only path) is that terms will not be
	// the wildcard, so we avoid three extra string compares earlier in
	// this function, and instead test for wildcards with numerics here.
	if sid == 0 || pid == 0 || oid == 0 {
		panic("Unexpected use of wildcard '*' for term")
	}
	// Add triple to all indexes.
	if !addToIndex(g.spoIndex, sid, pid, oid) {
		// Already existed.
		s.pool.releaseRefString(sid)
		s.pool.releaseRefString(pid)
		s.pool.releaseRefAny(oid)
		return false
	}
	addToIndex(g.posIndex, pid, oid, sid)
	addToIndex(g.ospIndex, oid, sid, pid)
	// Update size.
	s.size++
	g.size++
	if s.OnAdd != nil {
		s.OnAdd(subject, predicate, object, graph)
	}
	return true
}

// addToIndex adds a triple to the given index,
// creating deeper index buckets as needed.
// Returns true if the entry did not exist before.
func addToIndex(index0 indexRoot, key0, key1, key2 uint64) bool {
	index1, ok := index0[key0]
	if !ok {
		index1 = make(indexMid)
		index0[key0] = index1
	}
	index2, ok := index1[key1]
	if !ok {
		index2 = make(indexLeaf)
		index1[key1] = index2
	}
	_, exists := index2[key2]
	if !exists {
		index2[key2] = struct{}{}
	}
	return !exists
}

// Remove quads from the store. Returns the number of quads removed.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) Remove(subject, predicate string, object interface{}, graph string) uint64 {
	// Find internal identifiers for terms.
	sid, sok := s.pool.stringToID(subject)
	pid, pok := s.pool.stringToID(predicate)
	oid, ook := s.pool.anyToID(object)
	// If any of the terms don't exist, then there are no matches.
	if !sok || !pok || !ook {
		return 0
	}

	removeFromIndex := func(index0 indexRoot, key0, key1, key2 uint64, fn func(key0, key1, key2 uint64)) {
		index0.forEachMatch(key0, func(key0 uint64, index1 indexMid) {
			index1.forEachMatch(key1, func(key1 uint64, index2 indexLeaf) {
				index2.forEachMatch(key2, func(key2 uint64) {
					delete(index2, key2)
					// To ensure the indexes are in a consistent state
					// if/when we call OnRemove, we do any cleanup immediately.
					if len(index2) == 0 {
						delete(index1, key1)
						if len(index1) == 0 {
							delete(index0, key0)
						}
					}
					if fn != nil {
						fn(key0, key1, key2)
					}
				})
			})
		})
		// We do not remove the root bucket, even if it is empty.
	}

	var count uint64
	s.graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {

		// This is only called while processing the SPO index.
		removeFn := func(sid, pid, oid uint64) {
			s.size--
			g.size--
			if s.OnRemove != nil {
				s.OnRemove(s.pool.idToString(sid), s.pool.idToString(pid), s.pool.idToAny(oid), graph)
			}
			s.pool.releaseRefString(sid)
			s.pool.releaseRefString(pid)
			s.pool.releaseRefAny(oid)
			count++
		}

		// Remove matching elements from all indexes.
		removeFromIndex(g.posIndex, pid, oid, sid, nil)
		removeFromIndex(g.ospIndex, oid, sid, pid, nil)
		removeFromIndex(g.spoIndex, sid, pid, oid, removeFn)
		// Cleanup empty graphs.
		if g.size == 0 {
			delete(s.graphs, graph)
		}
	})
	return count
}

// Inversion of control - the index buckets themselves
// take care of any wilcards and call back as they need to.

// Lazy helper, for less error prone / more readable code elsewhere.
func (gm graphMap) forEachMatch(query string, fn func(key string, g *indexedGraph)) {
	gm.someMatch(query, func(key string, g *indexedGraph) bool {
		fn(key, g)
		return false
	})
}

func (gm graphMap) someMatch(query string, fn func(key string, g *indexedGraph) bool) bool {
	// Either loop over all graphs, or over just one selected graph.
	if query == "*" {
		// All graphs.
		for key, g := range gm {
			if fn(key, g) {
				return true
			}
		}
	} else {
		// Single graph - if it exists.
		g, ok := gm[query]
		if ok {
			return fn(query, g)
		}
	}
	return false
}

// These three functions all operate identically,
// but differ because of the specific types at each layer.

func (idx indexRoot) forEachMatch(query uint64, fn func(key uint64, idx indexMid)) {
	// Either loop over all elements, or over just one selected element.
	if query == 0 {
		// All elements.
		for key, i := range idx {
			fn(key, i)
		}
	} else {
		// Single element - if it exists.
		i, ok := idx[query]
		if ok {
			fn(query, i)
		}
	}
}

func (idx indexMid) forEachMatch(query uint64, fn func(key uint64, idx indexLeaf)) {
	// Either loop over all elements, or over just one selected element.
	if query == 0 {
		// All elements.
		for key, i := range idx {
			fn(key, i)
		}
	} else {
		// Single element - if it exists.
		i, ok := idx[query]
		if ok {
			fn(query, i)
		}
	}
}

func (idx indexLeaf) forEachMatch(query uint64, fn func(key uint64)) {
	// Either loop over all elements, or over just one selected element.
	if query == 0 {
		// All elements.
		for key := range idx {
			fn(key)
		}
	} else {
		// Single element - if it exists.
		_, ok := idx[query]
		if ok {
			fn(query)
		}
	}
}

// Count returns a count of quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) Count(subject, predicate string, object interface{}, graph string) uint64 {
	// Find internal identifiers for terms.
	sid, sok := s.pool.stringToID(subject)
	pid, pok := s.pool.stringToID(predicate)
	oid, ook := s.pool.anyToID(object)
	// If any of the terms don't exist, then there are no matches.
	if !sok || !pok || !ook {
		return 0
	}

	var count uint64
	s.graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {
		// Choose the optimal index, based on which fields are wildcards.
		if sid != 0 {
			if oid != 0 {
				// If subject and object are given, the ospIndex will be fastest.
				count += countInIndex(g.ospIndex, oid, sid, pid)
			} else {
				// If subject and possibly predicate are given, the spoIndex will be fastest.
				count += countInIndex(g.spoIndex, sid, pid, oid)
			}
		} else {
			if pid != 0 {
				// If only predicate and possibly object are given, the posIndex will be fastest.
				count += countInIndex(g.posIndex, pid, oid, sid)
			} else if oid != 0 {
				// If only object is given, the ospIndex will be fastest.
				count += countInIndex(g.ospIndex, oid, sid, pid)
			} else {
				// If all wildcard params given, use the graph size.
				count += g.size
			}
		}
	})
	return count
}

func countInIndex(index0 indexRoot, key0, key1, key2 uint64) uint64 {
	var count uint64
	index0.forEachMatch(key0, func(key0 uint64, index1 indexMid) {
		index1.forEachMatch(key1, func(key1 uint64, index2 indexLeaf) {
			if key2 == 0 {
				// key2 is wildcard, count all entries of index2.
				count += uint64(len(index2))
			} else {
				// Count single entry of index2, if it exists.
				_, ok := index2[key2]
				if ok {
					count++
				}
			}
		})
	})
	return count
}

// ForEach executes the given callback once for each quad in the store.
func (s *QuadStore) ForEach(fn QuadCallbackFn) {
	s.ForEachWith("*", "*", "*", "*", fn)
}

// ForEachWith executes the given callback once for each quad in the store
// that matches the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) ForEachWith(subject, predicate string, object interface{}, graph string, fn QuadCallbackFn) {
	iterAllFnWrapper := func(s, p string, o interface{}, g string) bool {
		fn(s, p, o, g)
		return false
	}
	s.SomeWith(subject, predicate, object, graph, iterAllFnWrapper)
}

// Every tests whether all quads in the store pass the test
// implemented by the given function.
//
// The given callback is
// executed once for each quad present in the store until
// Every finds one where the callback returns false. If such
// an element is found, iteration is immediately halted and
// Every returns false. Otherwise, if the callback returns
// true for all quads, then Every returns true.
//
// If no quads match the given terms, or the store is empty,
// then Every returns false. Note that this differs from
// the interpretation of 'every' in some other languages,
// which may return true for an empty iteration set.
func (s *QuadStore) Every(fn QuadTestFn) bool {
	return s.EveryWith("*", "*", "*", "*", fn)
}

// EveryWith tests whether all quads in the store that match the
// given terms pass the test implemented by the given function.
//
// The given callback is
// executed once for each matching quad in the store until
// EveryWith finds one where the callback returns false. If such
// an element is found, iteration is immediately halted and
// EveryWith returns false. Otherwise, if the callback returns
// true for all quads, then EveryWith returns true.
//
// If no quads match the given terms, or the store is empty,
// then EveryWith returns false. Note that this differs from
// the interpretation of 'every' in some other languages,
// which may return true for an empty iteration set.
func (s *QuadStore) EveryWith(subject, predicate string, object interface{}, graph string, fn QuadTestFn) bool {
	some := false
	everyFn := func(s, p string, o interface{}, g string) bool {
		some = true
		return !fn(s, p, o, g)
	}
	every := !s.SomeWith(subject, predicate, object, graph, everyFn)
	// Fixup the 'for-all quantifier in maths' stuff - which
	// plainly is not useful, and violates the principal of least
	// surprise - so now, we do not return true if the iteration
	// set was empty.
	if !some {
		return false
	}
	return every
}

// Some tests whether some quad in the store passes the test
// implemented by the given function.
//
// The given callback is
// executed once for each quad present in the store until
// Some finds one where the callback returns true. If such
// an element is found, iteration is immediately halted and
// Some returns true. Otherwise, if the callback returns
// false for all quads, then Some returns false.
func (s *QuadStore) Some(fn QuadTestFn) bool {
	return s.SomeWith("*", "*", "*", "*", fn)
}

const (
	_s = iota
	_p
	_o
	_g
)

// SomeWith tests whether some quad matching the given pattern
// passes the test implemented by the given function.
//
// The given callback is
// executed once for each quad matching the given pattern until
// SomeWith finds one where the callback returns true. If such
// an element is found, iteration is immediately halted and
// SomeWith returns true. Otherwise, if the callback returns
// false for all quads, then SomeWith returns false.
func (s *QuadStore) SomeWith(subject, predicate string, object interface{}, graph string, fn QuadTestFn) bool {
	// Find internal identifiers for terms.
	sid, sok := s.pool.stringToID(subject)
	pid, pok := s.pool.stringToID(predicate)
	oid, ook := s.pool.anyToID(object)
	// If any of the terms don't exist, then there are no matches.
	if !sok || !pok || !ook {
		return false
	}

	// flags := 0
	// if sid != 0 {
	// 	flags |= 4
	// }
	// if pid != 0 {
	// 	flags |= 2
	// }
	// if oid != 0 {
	// 	flags |= 1
	// }

	// matchFn := func(graph string, g *indexedGraph) bool {

	// 	// Currently, branch prediction beats table lookup - by approx 10%.

	// 	fns := [8]func() bool{
	// 		// s = z : p = z : o = z
	// 		func() bool { return indexSomeGivenNoKeys(g.spoIndex, _s, _p, _o, graph, s, fn) },
	// 		// s = z : p = z : o = nz
	// 		func() bool { return indexSomeGivenKey0(g.ospIndex, oid, _o, _s, _p, graph, s, fn) },
	// 		// s = z : p = nz : o = z
	// 		func() bool { return indexSomeGivenKey0(g.posIndex, pid, _p, _o, _s, graph, s, fn) },
	// 		// s = z : p = nz : o = nz
	// 		func() bool { return indexSomeGivenKey0And1(g.posIndex, pid, oid, _p, _o, _s, graph, s, fn) },
	// 		// s = nz : p = z : o = z
	// 		func() bool { return indexSomeGivenKey0(g.spoIndex, sid, _s, _p, _o, graph, s, fn) },
	// 		// s = nz : p = z : o = nz
	// 		func() bool { return indexSomeGivenKey0And1(g.ospIndex, oid, sid, _o, _s, _p, graph, s, fn) },
	// 		// s = nz : p = nz : o = z
	// 		func() bool { return indexSomeGivenKey0And1(g.spoIndex, sid, pid, _s, _p, _o, graph, s, fn) },
	// 		// s = nz : p = nz : o = nz
	// 		func() bool { return indexSomeGivenAllKeys(g.spoIndex, sid, pid, oid, _s, _p, _o, graph, s, fn) },
	// 	}
	// 	return fns[flags]()
	// }

	matchFn := func(graph string, g *indexedGraph) bool {
		// Choose the optimal index, based on which fields are wildcards.
		if sid != 0 {
			if pid != 0 {
				if oid != 0 {
					// s = nz : p = nz : o = nz
					return indexSomeGivenAllKeys(g.spoIndex, sid, pid, oid, _s, _p, _o, graph, s, fn)
				}
				// s = nz : p = nz : o = z
				return indexSomeGivenKey0And1(g.spoIndex, sid, pid, _s, _p, _o, graph, s, fn)
			}
			if oid != 0 {
				// s = nz : p = z : o = nz
				return indexSomeGivenKey0And1(g.ospIndex, oid, sid, _o, _s, _p, graph, s, fn)
			}
			// s = nz : p = z : o = z
			return indexSomeGivenKey0(g.spoIndex, sid, _s, _p, _o, graph, s, fn)
		}
		if pid != 0 {
			if oid != 0 {
				// s = z : p = nz : o = nz
				return indexSomeGivenKey0And1(g.posIndex, pid, oid, _p, _o, _s, graph, s, fn)
			}
			// s = z : p = nz : o = z
			return indexSomeGivenKey0(g.posIndex, pid, _p, _o, _s, graph, s, fn)
		}
		if oid != 0 {
			// s = z : p = z : o = nz
			return indexSomeGivenKey0(g.ospIndex, oid, _o, _s, _p, graph, s, fn)
		}
		// s = z : p = z : o = z
		return indexSomeGivenNoKeys(g.spoIndex, _s, _p, _o, graph, s, fn)
	}

	// The magic numbers above (_x slot numbers) should really be properties of the index itself.
	//
	// In an ideal world, the decision as to which index to use should be function
	// that looks at given params and what indexes are present - then it would be possible
	// to add or remove indexes.

	return s.graphs.someMatch(graph, matchFn)
}

func indexSomeGivenNoKeys(index0 indexRoot, idx0, idx1, idx2 int, g string, s *QuadStore, fn QuadTestFn) bool {
	var t [3]interface{} // spo triple
	// Loop.
	for key0, index1 := range index0 {
		t[idx0] = s.pool.idToAny(key0)
		// Loop.
		for key1, index2 := range index1 {
			t[idx1] = s.pool.idToAny(key1)
			// Loop.
			for key2 := range index2 {
				t[idx2] = s.pool.idToAny(key2)
				if fn(t[0].(string), t[1].(string), t[2], g) {
					return true
				}
			}
		}
	}
	return false
}

func indexSomeGivenKey0(index0 indexRoot, key0 uint64, idx0, idx1, idx2 int, g string, s *QuadStore, fn QuadTestFn) bool {
	var t [3]interface{} // spo triple
	// Lookup.
	index1, ok := index0[key0]
	if !ok {
		return false
	}
	t[idx0] = s.pool.idToAny(key0)
	// Loop.
	for key1, index2 := range index1 {
		t[idx1] = s.pool.idToAny(key1)
		// Loop.
		for key2 := range index2 {
			t[idx2] = s.pool.idToAny(key2)
			if fn(t[0].(string), t[1].(string), t[2], g) {
				return true
			}
		}
	}
	return false
}

func indexSomeGivenKey0And1(index0 indexRoot, key0, key1 uint64, idx0, idx1, idx2 int, g string, s *QuadStore, fn QuadTestFn) bool {
	var t [3]interface{} // spo triple
	// Lookup.
	index1, ok := index0[key0]
	if !ok {
		return false
	}
	t[idx0] = s.pool.idToAny(key0)
	// Lookup.
	index2, ok := index1[key1]
	if !ok {
		return false
	}
	t[idx1] = s.pool.idToAny(key1)
	// Loop.
	for key2 := range index2 {
		t[idx2] = s.pool.idToAny(key2)
		if fn(t[0].(string), t[1].(string), t[2], g) {
			return true
		}
	}
	return false
}

func indexSomeGivenAllKeys(index0 indexRoot, key0, key1, key2 uint64, idx0, idx1, idx2 int, g string, s *QuadStore, fn QuadTestFn) bool {
	var t [3]interface{} // spo triple
	// Lookup.
	index1, ok := index0[key0]
	if !ok {
		return false
	}
	t[idx0] = s.pool.idToAny(key0)
	// Lookup.
	index2, ok := index1[key1]
	if !ok {
		return false
	}
	t[idx1] = s.pool.idToAny(key1)
	// Lookup.
	_, ok = index2[key2]
	if !ok {
		return false
	}
	t[idx2] = s.pool.idToAny(key2)
	return fn(t[0].(string), t[1].(string), t[2], g)
}

// FindGraphs returns a list of distinct graph names for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) FindGraphs(subject, predicate string, object interface{}) []string {
	var out []string
	collectResultsFn := func(g string) {
		out = append(out, g)
	}
	s.ForGraphs(subject, predicate, object, collectResultsFn)
	return out
}

// ForGraphs executes the given callback once for each distinct graph name
// for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) ForGraphs(subject, predicate string, object interface{}, fn StringCallbackFn) {
	callbackAndBreakFn := func(s, p string, o interface{}, g string) bool {
		fn(g)
		return true
	}
	for graph := range s.graphs {
		s.SomeWith(subject, predicate, object, graph, callbackAndBreakFn)
	}
}

// FindSubjects returns a list of distinct subject terms for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) FindSubjects(predicate string, object interface{}, graph string) []string {
	var out []string
	collectResultsFn := func(s string) {
		out = append(out, s)
	}
	s.ForSubjects(predicate, object, graph, collectResultsFn)
	return out
}

// ForSubjects executes the given callback once for each distinct subject term
// for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) ForSubjects(predicate string, object interface{}, graph string, fn StringCallbackFn) {
	// Find internal identifiers for terms.
	pid, pok := s.pool.stringToID(predicate)
	oid, ook := s.pool.anyToID(object)
	// If any of the terms don't exist, then there are no matches.
	if !pok || !ook {
		return
	}

	var seen = make(map[uint64]struct{})

	collectResultsFn := func(id uint64) {
		_, ok := seen[id]
		if !ok {
			seen[id] = struct{}{}
			fn(s.pool.idToString(id))
		}
	}

	s.graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {
		// We want to list all subjects.
		// The three index choices are: SPO POS OSP

		// Choose the optimal index, based on which fields are wildcards.
		if pid != 0 {
			if oid != 0 {
				// If predicate and object are given, the posIndex is best.
				// Lookup p, lookup o, loop s.
				index2KeysGivenKey0And1(g.posIndex, pid, oid, collectResultsFn)
			} else {
				// If only predicate is given, the spoIndex is best.
				// Loop s, lookup p.
				index0KeysGivenKey1(g.spoIndex, pid, collectResultsFn)
			}
		} else {
			if oid != 0 {
				// If only object is given, the ospIndex is best.
				// Lookup o, loop s.
				index1KeysGivenKey0(g.ospIndex, oid, collectResultsFn)
			} else {
				// If no params given, iterate all the subjects.
				index0Keys(g.spoIndex, collectResultsFn)
			}
		}
	})
}

// FindPredicates returns a list of distinct predicate terms for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) FindPredicates(subject string, object interface{}, graph string) []string {
	var out []string
	collectResultsFn := func(p string) {
		out = append(out, p)
	}
	s.ForPredicates(subject, object, graph, collectResultsFn)
	return out
}

// ForPredicates executes the given callback once for each distinct predicate term
// for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) ForPredicates(subject string, object interface{}, graph string, fn StringCallbackFn) {
	// Find internal identifiers for terms.
	sid, sok := s.pool.stringToID(subject)
	oid, ook := s.pool.anyToID(object)
	// If any of the terms don't exist, then there are no matches.
	if !sok || !ook {
		return
	}

	var seen = make(map[uint64]struct{})

	collectResultsFn := func(id uint64) {
		_, ok := seen[id]
		if !ok {
			seen[id] = struct{}{}
			fn(s.pool.idToString(id))
		}
	}

	s.graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {
		// We want to list all predicates.
		// The three index choices are: SPO POS OSP

		// Choose the optimal index, based on which fields are wildcards.
		if sid != 0 {
			if oid != 0 {
				// If subject and object are given, the ospIndex is best.
				// Lookup o, lookup s, loop p.
				index2KeysGivenKey0And1(g.ospIndex, oid, sid, collectResultsFn)
			} else {
				// If only subject is given, the spoIndex is best.
				// Lookup s, loop p.
				index1KeysGivenKey0(g.spoIndex, sid, collectResultsFn)
			}
		} else {
			if oid != 0 {
				// If only object is given, the posIndex is best.
				// Loop p, lookup o.
				index0KeysGivenKey1(g.posIndex, oid, collectResultsFn)
			} else {
				// If no params given, iterate all the predicates.
				index0Keys(g.posIndex, collectResultsFn)
			}
		}
	})
}

// FindObjects returns a list of distinct object terms for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) FindObjects(subject, predicate, graph string) []interface{} {
	var out []interface{}
	collectResultsFn := func(o interface{}) {
		out = append(out, o)
	}
	s.ForObjects(subject, predicate, graph, collectResultsFn)
	return out
}

// ForObjects executes the given callback once for each distinct object term
// for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) ForObjects(subject, predicate, graph string, fn ObjectCallbackFn) {
	// Find internal identifiers for terms.
	sid, sok := s.pool.stringToID(subject)
	pid, pok := s.pool.stringToID(predicate)
	// If any of the terms don't exist, then there are no matches.
	if !sok || !pok {
		return
	}

	var seen = make(map[uint64]struct{})

	collectResultsFn := func(id uint64) {
		_, ok := seen[id]
		if !ok {
			seen[id] = struct{}{}
			fn(s.pool.idToAny(id))
		}
	}

	s.graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {
		// We want to list all objects.
		// The three index choices are: SPO POS OSP

		// Choose the optimal index, based on which fields are wildcards.
		if sid != 0 {
			if pid != 0 {
				// If subject and predicate are given, the spoIndex is best.
				// Lookup s, lookup p, loop o.
				index2KeysGivenKey0And1(g.spoIndex, sid, pid, collectResultsFn)
			} else {
				// If only subject is given, the ospIndex is best.
				// Loop o, lookup s.
				index0KeysGivenKey1(g.ospIndex, sid, collectResultsFn)
			}
		} else {
			if pid != 0 {
				// If only predicate is given, the posIndex is best.
				// Lookup p, loop o.
				index1KeysGivenKey0(g.posIndex, pid, collectResultsFn)
			} else {
				// If no params given, iterate all the objects.
				index0Keys(g.ospIndex, collectResultsFn)
			}
		}
	})
}

func index2KeysGivenKey0And1(index0 indexRoot, key0, key1 uint64, fn func(key2 uint64)) {
	// Lookup.
	index1, ok := index0[key0]
	if !ok {
		return
	}
	// Lookup.
	index2, _ := index1[key1]
	// Loop.
	for key2 := range index2 {
		fn(key2)
	}
}

func index1KeysGivenKey0(index0 indexRoot, key0 uint64, fn func(key1 uint64)) {
	// Lookup.
	index1, ok := index0[key0]
	if !ok {
		return
	}
	// Loop.
	for key1 := range index1 {
		fn(key1)
	}
}

func index0KeysGivenKey1(index0 indexRoot, key1 uint64, fn func(key0 uint64)) {
	// Loop
	for key0, index1 := range index0 {
		// Lookup.
		_, ok := index1[key1]
		if ok {
			fn(key0)
		}
	}
}

func index0Keys(index0 indexRoot, fn func(key0 uint64)) {
	// Loop
	for key0 := range index0 {
		fn(key0)
	}
}

// String returns the contents of the quad store in a human-readable format.
func (s *QuadStore) String() string {
	var buf bytes.Buffer
	graphs := s.FindGraphs("*", "*", "*")
	sort.Strings(graphs)
	for _, graph := range graphs {
		subjects := s.FindSubjects("*", "*", graph)
		sort.Strings(subjects)
		for _, subject := range subjects {
			predicates := s.FindPredicates(subject, "*", graph)
			sort.Strings(predicates)
			for _, predicate := range predicates {
				objects := s.FindObjects(subject, predicate, graph)
				sortObjects(objects)
				for _, object := range objects {
					buf.WriteByte('[')
					buf.WriteString(subject)
					buf.WriteByte(' ')
					buf.WriteString(predicate)
					buf.WriteByte(' ')
					buf.WriteString(fmt.Sprint(object))
					buf.WriteByte(' ')
					buf.WriteString(graph)
					buf.WriteByte(']')
					buf.WriteByte('\n')
				}
			}
		}
	}
	return buf.String()
}
