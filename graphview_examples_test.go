package store4_test

import (
	"fmt"
	"sort"

	"github.com/jimsmart/store4"
)

// tripleSlice implements sort.Interface for [][3]string
// ordering by fields SPO [0,1,2].
type tripleSlice [][3]string

func (t tripleSlice) Len() int { return len(t) }

func (t tripleSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t tripleSlice) Less(i, j int) bool {
	ti, tj := t[i], t[j]
	// Subject.
	si, sj := ti[0], tj[0]
	if si < sj {
		return true
	}
	if si > sj {
		return false
	}
	// Predicate.
	pi, pj := ti[1], tj[1]
	if pi < pj {
		return true
	}
	if pi > pj {
		return false
	}
	// Object.
	oi, oj := ti[2], tj[2]
	return oi < oj
}

func sortTriples(slice [][3]string) {
	sort.Sort(tripleSlice(slice))
}

func ExampleNewGraph() {

	// A new empty graph.
	g1 := store4.NewGraph()
	fmt.Println(g1.String())

	// A new graph initialised with a slice of triples.
	g2 := store4.NewGraph([][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
	})
	fmt.Println(g2)
	fmt.Println(g2.QuadStore)

	// Output:
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
	//
	// [s1 p1 o1 ]
	// [s1 p2 o2 ]
	// [s2 p2 o2 ]
}

func ExampleQuadStore_GraphView() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Get GraphView for graph g1.
	g1 := s.GraphView("g1")
	fmt.Println(g1)

	// Get GraphView for graph g2.
	g2 := s.GraphView("g2")
	fmt.Println(g2)

	// Get GraphView for graph g3 (which is empty).
	g3 := s.GraphView("g3")
	fmt.Println(g3)
	fmt.Println(g3.Size())

	// Use wildcard to get GraphView for
	// union of all graphs. (Note that attempts
	// to Add to this GraphView will panic)
	g4 := s.GraphView("*")
	fmt.Println(g4)

	// Output:
	// g1
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
	//
	// g2
	// [s2 p2 o3]
	// [s3 p3 o3]
	//
	// g3
	//
	// 0
	// *
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
	// [s2 p2 o3]
	// [s3 p3 o3]
}

func ExampleQuadStore_GraphViews() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Get GraphViews for all graphs in the store.
	results1 := s.GraphViews("*", "*", "*")
	fmt.Println(len(results1))
	// (Exact order may vary)
	//fmt.Println(results1[0])
	//fmt.Println(results1[1])

	// Get GraphViews for all graphs featuring
	// quads with subject s1 and predicate p2.
	results2 := s.GraphViews("s1", "p2", "*")
	fmt.Println(len(results2))
	fmt.Println(results2[0])

	// Output:
	// 2
	// 1
	// g1
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
}

func ExampleGraphView_Add() {

	triples := [][4]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
		{"s2", "p2", "o3"},
		{"s3", "p3", "o3"},
	}

	g := store4.NewGraph()

	for _, t := range triples {
		// Add triples to the graph, testing for duplicates.
		ok := g.Add(t[0], t[1], t[2])
		if !ok {
			panic("duplicate triple")
		}
	}

	fmt.Println(g)
	fmt.Println(g.QuadStore)
	// Output:
	//
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
	// [s2 p2 o3]
	// [s3 p3 o3]
	//
	// [s1 p1 o1 ]
	// [s1 p2 o2 ]
	// [s2 p2 o2 ]
	// [s2 p2 o3 ]
	// [s3 p3 o3 ]
}

func ExampleGraphView_Count() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	// Count triples that have subject s1.
	count := g1.Count("s1", "*", "*")
	fmt.Println(count)

	// Count triples that have predicate p2.
	count = g1.Count("*", "p2", "*")
	fmt.Println(count)

	// Output:
	// 2
	// 2
}

func ExampleGraphView_Empty() {
	g := store4.NewGraph()
	fmt.Println(g.Empty())

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})
	g1 := s.GraphView("g1")
	fmt.Println(g1.Empty())
	// Output:
	// true
	// false
}

func ExampleGraphView_Every() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g2"},
	})

	g1 := s.GraphView("g1")
	g2 := s.GraphView("g2")

	s1TestFn := func(s, p, o string) bool {
		return s == "s1"
	}

	// Iterate over every triple in the graph
	// while true is being returned from our
	// callback, and halt when false is returned.
	// Returns false if the callback ever returned false.
	result := g1.Every(s1TestFn)
	fmt.Println(result)
	result = g2.Every(s1TestFn)
	fmt.Println(result)

	// Note that Every will return true
	// for an empty graph.
	g0 := store4.NewGraph()
	result = g0.Every(s1TestFn)
	fmt.Println(result)

	// Every is often used as a breakable iterator,
	// with its return value being ignored.

	// Output:
	// true
	// false
	// true
}

func ExampleGraphView_EveryWith() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s2", "p2", "o1", "g1"},
		{"s2", "p3", "o1", "g1"},
		{"s2", "p3", "o2", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	o1TestFn := func(s, p, o string) bool {
		return o == "o1"
	}

	g1 := s.GraphView("g1")

	// Iterate over every quad in the graph having subject s2
	// while true is being returned from our
	// callback, and halt when false is returned.
	// Returns false if the callback ever returned false.
	result := g1.EveryWith("s2", "*", "*", o1TestFn)
	fmt.Println(result)

	// Note that EveryWith will return true
	// for an empty graph...
	g0 := store4.NewGraph()
	result = g0.EveryWith("*", "*", "*", o1TestFn)
	fmt.Println(result)
	// ...or if its iteration set is empty.
	result = g1.EveryWith("s0", "*", "*", o1TestFn)
	fmt.Println(result)

	// EveryWith is often used as a breakable iterator,
	// with its return value being ignored.

	// Output:
	// true
	// true
	// true
}

func ExampleGraphView_FindSubjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	// Get a list of all subjects in the graph.
	results := g1.FindSubjects("*", "*")
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results)
	fmt.Println(results)

	// Find subjects for triples in the graph
	// that have both predicate p2 and object o2.
	results = g1.FindSubjects("p2", "o2")
	sort.Strings(results)
	fmt.Println(results)

	// Output:
	// [s1 s2]
	// [s1 s2]
}

func ExampleGraphView_FindPredicates() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	// Get a list of all predicates in graph g1.
	results := g1.FindPredicates("*", "*")
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results)
	fmt.Println(results)

	// Find predicates for quads in graph g1 that have
	// both subject s1 and object o2.
	results = g1.FindPredicates("s1", "o2")
	sort.Strings(results)
	fmt.Println(results)

	// Output:
	// [p1 p2]
	// [p2]
}

func ExampleGraphView_FindObjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	// Get a list of all objects in graph g1.
	results := g1.FindObjects("*", "*")
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results)
	fmt.Println(results)

	union := s.GraphView("*")

	// Find objects for triples in any graph that have
	// both subject s1 and predicate p2.
	results = union.FindObjects("s2", "p2")
	sort.Strings(results)
	fmt.Println(results)

	// Output:
	// [o1 o2]
	// [o2 o3]
}

func ExampleGraphView_ForEach() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	// Iterate over all triples.
	var results [][3]string
	g1.ForEach(func(s, p, o string) {
		results = append(results, [3]string{s, p, o})
	})

	// (We only sort the results before printing
	// because iteration order is unstable)
	sortTriples(results)
	for _, q := range results {
		fmt.Println(q)
	}

	// Output:
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
}

func ExampleGraphView_ForEachWith() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	// Iterate over all triples with predicate p2.
	var results [][3]string
	g1.ForEachWith("*", "p2", "*", func(s, p, o string) {
		results = append(results, [3]string{s, p, o})
	})

	// (We only sort the results before printing
	// because iteration order is unstable)
	sortTriples(results)
	for _, q := range results {
		fmt.Println(q)
	}

	// Output:
	// [s1 p2 o2]
	// [s2 p2 o2]
}

func ExampleGraphView_ForSubjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	var results1 []string
	// Iterate over all subjects in graph g1.
	g1.ForSubjects("*", "*", func(s string) {
		results1 = append(results1, s)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results1)
	fmt.Println(results1)

	g2 := s.GraphView("g2")

	var results2 []string
	// Iterate over all subjects in graph g2.
	g2.ForSubjects("*", "*", func(s string) {
		results2 = append(results2, s)
	})
	sort.Strings(results2)
	fmt.Println(results2)

	union := s.GraphView("*")

	var results3 []string
	// Iterate over subjects for triples in any graph that have
	// both predicate p2 and object o2.
	union.ForSubjects("p2", "o2", func(s string) {
		results3 = append(results3, s)
	})
	sort.Strings(results3)
	fmt.Println(results3)

	// Output:
	// [s1 s2]
	// [s2 s3]
	// [s1 s2]
}

func ExampleGraphView_ForPredicates() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	var results1 []string
	// Iterate over all predicates in graph g1.
	g1.ForPredicates("*", "*", func(p string) {
		results1 = append(results1, p)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results1)
	fmt.Println(results1)

	g2 := s.GraphView("g2")

	var results2 []string
	// Iterate over all predicates in graph g2.
	g2.ForPredicates("*", "*", func(p string) {
		results2 = append(results2, p)
	})
	sort.Strings(results2)
	fmt.Println(results2)

	var results3 []string
	// Iterate over predicates for triples in graph g1 that have
	// both subject s1 and object o2.
	g1.ForPredicates("s1", "o2", func(p string) {
		results3 = append(results3, p)
	})
	sort.Strings(results3)
	fmt.Println(results3)

	// Output:
	// [p1 p2]
	// [p2 p3]
	// [p2]
}

func ExampleGraphView_ForObjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	var results1 []string
	// Iterate over all objects in graph g1.
	g1.ForObjects("*", "*", func(o string) {
		results1 = append(results1, o)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results1)
	fmt.Println(results1)

	g2 := s.GraphView("g2")

	var results2 []string
	// Iterate over all objects in graph g2.
	g2.ForObjects("*", "*", func(o string) {
		results2 = append(results2, o)
	})
	sort.Strings(results2)
	fmt.Println(results2)

	var results3 []string
	// Iterate over objects for triples in graph g1 that have
	// both subject s1 and predicate p2.
	g1.ForObjects("s2", "p2", func(o string) {
		results3 = append(results3, o)
	})
	sort.Strings(results3)
	fmt.Println(results3)

	// Output:
	// [o1 o2]
	// [o3]
	// [o2]
}

func ExampleGraphView_Remove() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")
	fmt.Println(g1.Size())

	// Remove a specific triple from the graph.
	g1.Remove("s2", "p2", "o2")
	fmt.Println(g1.Size())

	// Remove all triples that have predicate p2.
	g1.Remove("*", "p2", "*")
	fmt.Println(g1.Size())

	g2 := s.GraphView("g2")
	fmt.Println(g2.Size())

	// Remove all triples from graph g2.
	g2.Remove("*", "*", "*")
	fmt.Println(g2.Size())

	fmt.Println(s)

	// Output:
	// 3
	// 2
	// 1
	// 2
	// 0
	// [s1 p1 o1 g1]
}

func ExampleGraphView_Size() {

	g := store4.NewGraph([][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
		{"s2", "p2", "o3"},
		{"s3", "p3", "o3"},
	})

	// How many triples are in the graph?
	count := g.Size()
	fmt.Println(count)

	// Output: 5
}

func ExampleGraphView_String() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")
	// (Println calls String)
	fmt.Println(g1)

	g2 := s.GraphView("g2")
	fmt.Println(g2)

	// Output:
	// g1
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
	//
	// g2
	// [s2 p2 o3]
	// [s3 p3 o3]
}

func ExampleGraphView_Some() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")
	g2 := s.GraphView("g2")

	s1TestFn := func(s, p, o string) bool {
		return s == "s1"
	}

	// Is there some quad in graph g1 with subject s1?
	result := g1.Some(s1TestFn)
	fmt.Println(result)

	// Is there some quad in graph g2 with subject s1?
	result = g2.Some(s1TestFn)
	fmt.Println(result)

	// Output:
	// true
	// false
}

func ExampleGraphView_SomeWith() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	alwaysTrueFn := func(s, p, o string) bool {
		return true
	}

	g1 := s.GraphView("g1")
	g2 := s.GraphView("g2")

	// Is there some quad in graph g1 with object o1?
	result := g1.SomeWith("*", "*", "o1", alwaysTrueFn)
	fmt.Println(result)

	// Is there some quad in graph g2 with object o1?
	result = g2.SomeWith("*", "*", "o1", alwaysTrueFn)
	fmt.Println(result)

	// Output:
	// true
	// false
}
