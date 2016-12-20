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

func ExampleGraphView_ForEach() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g1 := s.GraphView("g1")

	// Iterate over all triples using ForEach.
	var results [][3]string
	g1.ForEach(func(s, p, o string) {
		results = append(results, [3]string{s, p, o})
	})

	// (We only sort the results before printing
	// because iteration order is not stable)
	sortTriples(results)
	for _, q := range results {
		fmt.Println(q)
	}

	// Output:
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
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
	// because iteration order is not stable)
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
