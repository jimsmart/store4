package store4

import (
	"bytes"
	"fmt"
	"sort"
)

// QuadCallbackFn is the function signature used to implement
// callback functions that receive a quad.
//
// Used with calls to QuadStore's ForEach and ForEachWith.
type QuadCallbackFn func(s, p, o, g string)

// QuadTestFn is the function signature used to implement
// callback functions performing quad tests.
// A response of true means that the test has been passed.
//
// Used with calls to QuadStore's Every, EveryWith, Some and SomeWith.
type QuadTestFn func(s, p, o, g string) bool

// StringCallbackFn is the function signature used to implement
// callback functions that receive a string.
//
// Used with calls to FindSubjects, FindPredicates, FindObjects and FindGraphs.
type StringCallbackFn func(s string)

// QuadStore is an in-memory string-based quad store.
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
	// termToID maps terms to internal IDs.
	termToID map[string]uint64
	// idToTermInfo maps internal IDs to term info.
	idToTermInfo map[uint64]*termInfo
	// nextID holds the next term ID to issue.
	nextID uint64
}

// idToTerm returns the term for a given ID.
// The given ID must exist.
// This function replaces the simple map lookup
// after refactoring for reference counting.
func (s *QuadStore) idToTerm(id uint64) string {
	return s.idToTermInfo[id].term
}

// termInfo holds details for each term.
type termInfo struct {
	term     string // The term itself.
	refCount uint64 // Reference count.
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
func NewQuadStore(args ...interface{}) *QuadStore {
	s := &QuadStore{
		graphs: make(map[string]*indexedGraph),
	}
	s.termToID = make(map[string]uint64)
	s.idToTermInfo = make(map[uint64]*termInfo)
	// We use "*" as a wildcard, so we give it ID 0
	// to make things easy elsewhere.
	s.termToID["*"] = 0
	// Start IDs from 1.
	s.nextID = 1
	// Initialise store with any given data.
	for _, arg := range args {
		switch arg := arg.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T\n", arg))
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

// Size returns the total count of quads in the store.
func (s *QuadStore) Size() uint64 {
	return s.size
}

// Empty returns true is the store has no contents.
func (s *QuadStore) Empty() bool {
	return s.size == 0
}

// Add a quad to the store. Returns true if the quad was a new quad,
// or false if the quad already existed.
//
// If any of the given terms are "*" (an asterisk), then this method will panic.
// (The asterisk is reserved for wildcard operations throughout the API).
func (s *QuadStore) Add(subject, predicate, object, graph string) bool {
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
	sid := s.getOrCreateID(subject)
	pid := s.getOrCreateID(predicate)
	oid := s.getOrCreateID(object)
	// Disallow wildcard terms.
	// Optimisation: the fast path (only path) is that terms will not be
	// the wildcard, so we avoid three extra string compares earlier in
	// this function, and instead test for wildcards with numerics here.
	if sid == 0 || pid == 0 || oid == 0 {
		panic("Unexpected use of wildcard '*' for term")
	}
	// Add triple to all indexes.
	if !addToIndex(g.spoIndex, sid, pid, oid) {
		// Already exists.
		s.releaseRef(sid)
		s.releaseRef(pid)
		s.releaseRef(oid)
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

// getOrCreateID returns an internal ID for a given term.
// If no existing ID is present, it creates a new one.
// For any existing term, it also increments the reference count.
func (s *QuadStore) getOrCreateID(term string) uint64 {
	id, ok := s.termToID[term]
	if ok {
		if id != 0 {
			s.idToTermInfo[id].refCount++
		}
	} else {
		id = s.nextID
		s.nextID++
		s.termToID[term] = id
		s.idToTermInfo[id] = &termInfo{
			term:     term,
			refCount: 1,
		}
	}
	return id
}

// releaseRef decrements a term's reference count.
// When a term is no longer referenced, it is removed
// from all maps.
// The given id must exist or releaseRef will aspolde.
func (s *QuadStore) releaseRef(id uint64) {
	info := s.idToTermInfo[id]
	c := info.refCount
	c--
	if c == 0 {
		delete(s.termToID, info.term)
		delete(s.idToTermInfo, id)
		return
	}
	info.refCount = c
}

// Removes quads from the store. Returns the number of quads removed.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) Remove(subject, predicate, object, graph string) uint64 {
	termToID := s.termToID
	idToTerm := s.idToTerm
	graphs := s.graphs
	// Find internal identifiers for terms.
	sid, sok := termToID[subject]
	pid, pok := termToID[predicate]
	oid, ook := termToID[object]
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
	graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {

		// This is only called while processing the SPO index.
		removeFn := func(sid, pid, oid uint64) {
			s.size--
			g.size--
			if s.OnRemove != nil {
				s.OnRemove(idToTerm(sid), idToTerm(pid), idToTerm(oid), graph)
			}
			s.releaseRef(sid)
			s.releaseRef(pid)
			s.releaseRef(oid)
			count++
		}

		// Remove matching elements from all indexes.
		removeFromIndex(g.posIndex, pid, oid, sid, nil)
		removeFromIndex(g.ospIndex, oid, sid, pid, nil)
		removeFromIndex(g.spoIndex, sid, pid, oid, removeFn)
		// Cleanup empty graphs.
		if g.size == 0 {
			delete(graphs, graph)
		}
	})
	return count
}

// Count returns a count of quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) Count(subject, predicate, object, graph string) uint64 {
	termToId := s.termToID
	graphs := s.graphs
	// Find internal identifiers for terms.
	sid, sok := termToId[subject]
	pid, pok := termToId[predicate]
	oid, ook := termToId[object]
	// If any of the terms don't exist, then there are no matches.
	if !sok || !pok || !ook {
		return 0
	}

	var count uint64
	graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {
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
func (s *QuadStore) ForEachWith(subject, predicate, object, graph string, fn QuadCallbackFn) {
	iterAllFnWrapper := func(s, p, o, g string) bool {
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
// Acting like the 'for all' quantifier in maths, it should
// be noted that Every returns true for an empty store.
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
// Acting like the 'for all' quantifier in maths, it should
// be noted that EveryWith returns true for an empty store.
// By extension, if the given parameters cause the iteration
// set to be empty, then EveryWith also returns true.
func (s *QuadStore) EveryWith(subject, predicate, object, graph string, fn QuadTestFn) bool {
	everyFn := func(s, p, o, g string) bool {
		return !fn(s, p, o, g)
	}
	return !s.SomeWith(subject, predicate, object, graph, everyFn)
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

// SomeWith tests whether some quad matching the given pattern
// passes the test implemented by the given function.
//
// The given callback is
// executed once for each quad matching the given pattern until
// SomeWith finds one where the callback returns true. If such
// an element is found, iteration is immediately halted and
// SomeWith returns true. Otherwise, if the callback returns
// false for all quads, then SomeWith returns false.
func (s *QuadStore) SomeWith(subject, predicate, object, graph string, fn QuadTestFn) bool {
	termToID := s.termToID
	idToTerm := s.idToTerm
	graphs := s.graphs
	// Find internal identifiers for terms.
	sid, sok := termToID[subject]
	pid, pok := termToID[predicate]
	oid, ook := termToID[object]
	// If any of the terms don't exist, then there are no matches.
	if !sok || !pok || !ook {
		return false
	}

	var q [4]string // spog quad

	// Inlined, so it has easy access to term map and q results.
	indexSomeWith := func(index0 indexRoot, key0, key1, key2 uint64, idx0, idx1, idx2 int) bool {
		return index0.someMatch(key0, func(key0 uint64, index1 indexMid) bool {
			q[idx0] = idToTerm(key0)
			return index1.someMatch(key1, func(key1 uint64, index2 indexLeaf) bool {
				q[idx1] = idToTerm(key1)
				return index2.someMatch(key2, func(key2 uint64) bool {
					q[idx2] = idToTerm(key2)
					return fn(q[0], q[1], q[2], q[3])
				})
			})
		})
	}

	return graphs.someMatch(graph, func(graph string, g *indexedGraph) bool {
		q[3] = graph
		// Choose the optimal index, based on which fields are wildcards.
		if sid != 0 {
			if oid != 0 {
				// If subject and object are given, the ospIndex will be fastest.
				return indexSomeWith(g.ospIndex, oid, sid, pid, 2, 0, 1)
			} else {
				// If subject and possibly predicate are given, the spoIndex will be fastest.
				return indexSomeWith(g.spoIndex, sid, pid, oid, 0, 1, 2)
			}
		} else {
			if pid != 0 {
				// If only predicate and possibly object are given, the posIndex will be fastest.
				return indexSomeWith(g.posIndex, pid, oid, sid, 1, 2, 0)
			} else if oid != 0 {
				// If only object is given, the ospIndex will be fastest.
				return indexSomeWith(g.ospIndex, oid, sid, pid, 2, 0, 1)
			} else {
				// If no params given, iterate subject and predicates first.
				return indexSomeWith(g.spoIndex, sid, pid, oid, 0, 1, 2)
			}
			// The magic numbers (slot numbers) above should really be properties of the index itself.
			//
			// In an ideal world, the decision as to which index to use should be function
			// that looks at given params and what indexes are present - then it would be possible
			// to add or remove indexes.
		}
	})
}

// Inversion of control - the index buckets themselves
// take care of any wilcards and call back as they need to.

// These four functions all operate identically,
// but differ because of the specific types at each layer.

func (gm graphMap) someMatch(query string, fn func(key string, g *indexedGraph) bool) bool {
	// Either loop over all graphs, or over just one selected graph.
	if query == "*" {
		// All graphs.
		for key, g := range gm {
			if fn(key, g) {
				// break
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

func (idx indexRoot) someMatch(query uint64, fn func(key uint64, idx indexMid) bool) bool {
	// Either loop over all elements, or over just one selected element.
	if query == 0 {
		// All elements.
		for key, i := range idx {
			if fn(key, i) {
				return true // break
			}
		}
	} else {
		// Single element - if it exists.
		i, ok := idx[query]
		if ok {
			return fn(query, i)
		}
	}
	return false
}

func (idx indexMid) someMatch(query uint64, fn func(key uint64, idx indexLeaf) bool) bool {
	// Either loop over all elements, or over just one selected element.
	if query == 0 {
		// All elements.
		for key, i := range idx {
			if fn(key, i) {
				return true // break
			}
		}
	} else {
		// Single element - if it exists.
		i, ok := idx[query]
		if ok {
			return fn(query, i)
		}
	}
	return false
}

func (idx indexLeaf) someMatch(query uint64, fn func(key uint64) bool) bool {
	// Either loop over all elements, or over just one selected element.
	if query == 0 {
		// All elements.
		for key, _ := range idx {
			if fn(key) {
				return true // break
			}
		}
	} else {
		// Single element - if it exists.
		_, ok := idx[query]
		if ok {
			return fn(query)
		}
	}
	return false
}

// Lazy helpers, for less error prone / more readable code elsewhere.

func (gm graphMap) forEachMatch(query string, fn func(key string, g *indexedGraph)) {
	gm.someMatch(query, func(key string, g *indexedGraph) bool {
		fn(key, g)
		return false
	})
}

func (idx indexRoot) forEachMatch(query uint64, fn func(key uint64, idx indexMid)) {
	idx.someMatch(query, func(key uint64, i indexMid) bool {
		fn(key, i)
		return false
	})
}

func (idx indexMid) forEachMatch(query uint64, fn func(key uint64, idx indexLeaf)) {
	idx.someMatch(query, func(key uint64, i indexLeaf) bool {
		fn(key, i)
		return false
	})
}

func (idx indexLeaf) forEachMatch(query uint64, fn func(key uint64)) {
	idx.someMatch(query, func(key uint64) bool {
		fn(key)
		return false
	})
}

// FindGraphs returns a list of distinct graph names for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) FindGraphs(subject, predicate, object string) []string {
	var out []string
	collectResultsFn := func(s string) {
		out = append(out, s)
	}
	s.ForGraphs(subject, predicate, object, collectResultsFn)
	return out
}

// ForGraphs executes the given callback once for each distinct graph name
// for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) ForGraphs(subject, predicate, object string, fn StringCallbackFn) {
	callbackAndBreakFn := func(s, p, o, g string) bool {
		fn(g)
		return true
	}
	for graph, _ := range s.graphs {
		s.SomeWith(subject, predicate, object, graph, callbackAndBreakFn)
	}
}

// FindSubjects returns a list of distinct subject terms for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) FindSubjects(predicate, object, graph string) []string {
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
func (s *QuadStore) ForSubjects(predicate, object, graph string, fn StringCallbackFn) {
	termToID := s.termToID
	idToTerm := s.idToTerm
	graphs := s.graphs
	// Find internal identifiers for terms.
	pid, pok := termToID[predicate]
	oid, ook := termToID[object]
	// If any of the terms don't exist, then there are no matches.
	if !pok || !ook {
		return
	}

	var seen = make(map[uint64]struct{})

	collectResultsFn := func(id uint64) {
		_, ok := seen[id]
		if !ok {
			seen[id] = struct{}{}
			fn(idToTerm(id))
		}
	}

	graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {
		// We want to list all subjects.
		// The three index choices are: SPO POS OSP

		// Choose the optimal index, based on which fields are wildcards.
		if pid != 0 {
			if oid != 0 {
				// If predicate and object are given, the posIndex is best.
				// Lookup p, lookup o, loop s.
				lookupIndex0LookupIndex1LoopIndex2(g.posIndex, pid, oid, collectResultsFn)
			} else {
				// If only predicate is given, the spoIndex is best.
				// Loop s, lookup p.
				loopIndex0LookupIndex1(g.spoIndex, pid, collectResultsFn)
			}
		} else {
			if oid != 0 {
				// If only object is given, the ospIndex is best.
				// Lookup o, loop s.
				lookupIndex0LoopIndex1(g.ospIndex, oid, collectResultsFn)
			} else {
				// If no params given, iterate all the subjects.
				loopIndex0(g.spoIndex, collectResultsFn)
			}
		}
	})
}

// FindPredicates returns a list of distinct predicate terms for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) FindPredicates(subject, object, graph string) []string {
	var out []string
	collectResultsFn := func(s string) {
		out = append(out, s)
	}
	s.ForPredicates(subject, object, graph, collectResultsFn)
	return out
}

// ForPredicates executes the given callback once for each distinct predicate term
// for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) ForPredicates(subject, object, graph string, fn StringCallbackFn) {
	termToID := s.termToID
	idToTerm := s.idToTerm
	graphs := s.graphs
	// Find internal identifiers for terms.
	sid, sok := termToID[subject]
	oid, ook := termToID[object]
	// If any of the terms don't exist, then there are no matches.
	if !sok || !ook {
		return
	}

	var seen = make(map[uint64]struct{})

	collectResultsFn := func(id uint64) {
		_, ok := seen[id]
		if !ok {
			seen[id] = struct{}{}
			fn(idToTerm(id))
		}
	}

	graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {
		// We want to list all predicates.
		// The three index choices are: SPO POS OSP

		// Choose the optimal index, based on which fields are wildcards.
		if sid != 0 {
			if oid != 0 {
				// If subject and object are given, the ospIndex is best.
				// Lookup o, lookup s, loop p.
				lookupIndex0LookupIndex1LoopIndex2(g.ospIndex, oid, sid, collectResultsFn)
			} else {
				// If only subject is given, the spoIndex is best.
				// Lookup s, loop p.
				lookupIndex0LoopIndex1(g.spoIndex, sid, collectResultsFn)
			}
		} else {
			if oid != 0 {
				// If only object is given, the posIndex is best.
				// Loop p, lookup o.
				loopIndex0LookupIndex1(g.posIndex, oid, collectResultsFn)
			} else {
				// If no params given, iterate all the predicates.
				loopIndex0(g.posIndex, collectResultsFn)
			}
		}
	})
}

// FindObjects returns a list of distinct object terms for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) FindObjects(subject, predicate, graph string) []string {
	var out []string
	collectResultsFn := func(s string) {
		out = append(out, s)
	}
	s.ForObjects(subject, predicate, graph, collectResultsFn)
	return out
}

// ForObjects executes the given callback once for each distinct object term
// for all quads in the store that match the given pattern.
//
// Passing "*" (an asterisk) for any parameter acts as a
// match-everything wildcard for that term.
func (s *QuadStore) ForObjects(subject, predicate, graph string, fn StringCallbackFn) {
	termToID := s.termToID
	idToTerm := s.idToTerm
	graphs := s.graphs
	// Find internal identifiers for terms.
	sid, sok := termToID[subject]
	pid, pok := termToID[predicate]
	// If any of the terms don't exist, then there are no matches.
	if !sok || !pok {
		return
	}

	var seen = make(map[uint64]struct{})

	collectResultsFn := func(id uint64) {
		_, ok := seen[id]
		if !ok {
			seen[id] = struct{}{}
			fn(idToTerm(id))
		}
	}

	graphs.forEachMatch(graph, func(graph string, g *indexedGraph) {
		// We want to list all objects.
		// The three index choices are: SPO POS OSP

		// Choose the optimal index, based on which fields are wildcards.
		if sid != 0 {
			if pid != 0 {
				// If subject and predicate are given, the spoIndex is best.
				// Lookup s, lookup p, loop o.
				lookupIndex0LookupIndex1LoopIndex2(g.spoIndex, sid, pid, collectResultsFn)
			} else {
				// If only subject is given, the ospIndex is best.
				// Loop o, lookup s.
				loopIndex0LookupIndex1(g.ospIndex, sid, collectResultsFn)
			}
		} else {
			if pid != 0 {
				// If only predicate is given, the posIndex is best.
				// Lookup p, loop o.
				lookupIndex0LoopIndex1(g.posIndex, pid, collectResultsFn)
			} else {
				// If no params given, iterate all the objects.
				loopIndex0(g.ospIndex, collectResultsFn)
			}
		}
	})
}

// The names for these are terrible. :/

func lookupIndex0LookupIndex1LoopIndex2(index0 indexRoot, key0, key1 uint64, fn func(key2 uint64)) {
	// Lookup.
	index1, ok := index0[key0]
	if !ok {
		return
	}
	// Lookup.
	index2, _ := index1[key1]
	// Loop.
	for key2, _ := range index2 {
		fn(key2)
	}
}

func lookupIndex0LoopIndex1(index0 indexRoot, key0 uint64, fn func(key1 uint64)) {
	// Lookup.
	index1, ok := index0[key0]
	if !ok {
		return
	}
	// Loop.
	for key1, _ := range index1 {
		fn(key1)
	}
}

func loopIndex0LookupIndex1(index0 indexRoot, key1 uint64, fn func(key0 uint64)) {
	// Loop
	for key0, index1 := range index0 {
		// Lookup.
		_, ok := index1[key1]
		if ok {
			fn(key0)
		}
	}
}

func loopIndex0(index0 indexRoot, fn func(key0 uint64)) {
	// Loop
	for key0, _ := range index0 {
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
				sort.Strings(objects)
				for _, object := range objects {
					buf.WriteByte('[')
					buf.WriteString(subject)
					buf.WriteByte(' ')
					buf.WriteString(predicate)
					buf.WriteByte(' ')
					buf.WriteString(object)
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
