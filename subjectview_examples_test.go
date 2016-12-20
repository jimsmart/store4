package store4_test

import (
	"fmt"
	"sort"

	"github.com/jimsmart/store4"
)

// tupleSlice implements sort.Interface for [][2]string
// ordering by fields PO [0,1].
type tupleSlice [][2]string

func (t tupleSlice) Len() int { return len(t) }

func (t tupleSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t tupleSlice) Less(i, j int) bool {
	ti, tj := t[i], t[j]
	// Predicate.
	si, sj := ti[0], tj[0]
	if si < sj {
		return true
	}
	if si > sj {
		return false
	}
	// Object.
	oi, oj := ti[1], tj[1]
	return oi < oj
}

func sortTuples(slice [][2]string) {
	sort.Sort(tupleSlice(slice))
}

func ExampleGraphView_Query() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	g := s.GraphView("g1")

	// Query for SubjectViews of all subjects
	// that have p1=o1 and p2=o2.
	pattern := map[string]string{
		"p1": "o1",
		"p2": "o2",
	}
	results := g.Query(pattern)

	fmt.Println(len(results))
	fmt.Println(results[0])
	// Output:
	// 1
	// g1
	// s1
	// [p1 o1]
	// [p2 o2]
}

func ExampleQuadStore_Query() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Query for SubjectViews of all subjects in graph g1
	// that have p1=o1 and p2=o2.
	pattern := map[string]string{
		"p1": "o1",
		"p2": "o2",
	}
	results := s.Query(pattern, "g1")

	fmt.Println(len(results))
	fmt.Println(results[0].Subject)
	// Output:
	// 1
	// s1
}

func ExampleQuadStore_SubjectView() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Get a SubjectView of subject s1 in graph g1.
	s1 := s.SubjectView("s1", "g1")
	fmt.Println(s1)

	// Get a SubjectView of subject s2 in all graphs.
	// (Note that attempts to Add to this SubjectView will panic)
	s2 := s.SubjectView("s2", "*")
	fmt.Println(s2)

	// Create a SubjectView of subject s5 in graph g5.
	s4 := s.SubjectView("s5", "g5")
	// Add a quad.
	s4.Add("p5", "o5")
	fmt.Println(s)

	// Output:
	// g1
	// s1
	// [p1 o1]
	// [p2 o2]
	//
	// *
	// s2
	// [p2 o2]
	// [p2 o3]
	//
	// [s1 p1 o1 g1]
	// [s1 p2 o2 g1]
	// [s2 p2 o2 g1]
	// [s2 p2 o3 g2]
	// [s3 p3 o3 g2]
	// [s5 p5 o5 g5]
}

func ExampleGraphView_SubjectView() {

	g := store4.NewGraph([][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
	})
	// Get a SubjectView of subject s1.
	s1 := g.SubjectView("s1")
	fmt.Println(s1)

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})
	union := s.GraphView("*")
	// Get a SubjectView of subject s2 in all graphs.
	// (Note that attempts to Add to this SubjectView will panic)
	s2 := union.SubjectView("s2")
	fmt.Println(s2)

	g2 := store4.NewGraph()
	// Create a SubjectView of subject s5.
	s4 := g2.SubjectView("s5")
	// Add a quad.
	s4.Add("p5", "o5")
	fmt.Println(s4)

	// Output:
	// s1
	// [p1 o1]
	// [p2 o2]
	//
	// *
	// s2
	// [p2 o2]
	// [p2 o3]
	//
	// s5
	// [p5 o5]
}
