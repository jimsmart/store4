package store4_test

import (
	"fmt"
	"sort"

	"github.com/jimsmart/store4"
)

// quadSlice implements sort.Interface for [][4]string
// ordering by fields GSPO [3,0,1,2].
type quadSlice [][4]string

func (q quadSlice) Len() int { return len(q) }

func (q quadSlice) Swap(i, j int) { q[i], q[j] = q[j], q[i] }

func (q quadSlice) Less(i, j int) bool {
	qi, qj := q[i], q[j]
	// Graph.
	gi, gj := qi[3], qj[3]
	if gi < gj {
		return true
	}
	if gi > gj {
		return false
	}
	// Subject.
	si, sj := qi[0], qj[0]
	if si < sj {
		return true
	}
	if si > sj {
		return false
	}
	// Predicate.
	pi, pj := qi[1], qj[1]
	if pi < pj {
		return true
	}
	if pi > pj {
		return false
	}
	// Object.
	oi, oj := qi[2], qj[2]
	return oi < oj
}

func sortQuads(slice [][4]string) {
	sort.Sort(quadSlice(slice))
}

func sortStrings(slice []string) {
	sort.Strings(slice)
}

// func Example_titleHere() {}

// func ExampleNew() {}

func ExampleQuadStore_Add() {

	quads := [][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	}

	s := store4.NewQuadStore()

	for _, q := range quads {
		// Add quads to the store, testing for duplicates.
		ok := s.Add(q[0], q[1], q[2], q[3])
		if !ok {
			panic("duplicate quad")
		}
	}

	fmt.Println(s.String())
	// Output:
	// [s1 p1 o1 g1]
	// [s1 p2 o2 g1]
	// [s2 p2 o2 g1]
	// [s2 p2 o3 g2]
	// [s3 p3 o3 g2]
}

func ExampleQuadStore_Size() {

	quads := [][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	}

	s := store4.NewQuadStore()
	for _, q := range quads {
		s.Add(q[0], q[1], q[2], q[3])
	}

	// How many quads are in the store?
	count := s.Size()
	fmt.Println(count)

	// Output: 5
}

func ExampleQuadStore_Remove() {

	quads := [][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	}

	s := store4.NewQuadStore()
	for _, q := range quads {
		s.Add(q[0], q[1], q[2], q[3])
	}

	// Remove a specific quad from the store.
	s.Remove("s3", "p3", "o3", "g2")
	fmt.Println(s.Size())

	// Remove all quads that have object o3.
	s.Remove("*", "*", "o3", "*")
	fmt.Println(s.Size())

	// Remove all quads from the store.
	s.Remove("*", "*", "*", "*")
	fmt.Println(s.Size())

	// Output:
	// 4
	// 3
	// 0
}

func ExampleQuadStore_Count() {

	quads := [][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	}

	s := store4.NewQuadStore()
	for _, q := range quads {
		s.Add(q[0], q[1], q[2], q[3])
	}

	// Count quads that have subject s1.
	count := s.Count("s1", "*", "*", "*")
	fmt.Println(count)

	// Count quads that have predicate p2.
	count = s.Count("*", "p2", "*", "*")
	fmt.Println(count)

	// Count quads that have predicate p2 in graph g2.
	count = s.Count("*", "p2", "*", "g1")
	fmt.Println(count)

	// Output:
	// 2
	// 3
	// 2
}

func ExampleQuadStore_ForEach() {

	quads := [][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	}

	s := store4.NewQuadStore()
	for _, q := range quads {
		s.Add(q[0], q[1], q[2], q[3])
	}

	var results [][4]string
	s.ForEach(func(s, p, o, g string) {
		results = append(results, [4]string{s, p, o, g})
	})

	// (We only sort the results before printing
	// because iteration order is not stable)
	sortQuads(results)
	for _, q := range results {
		fmt.Println(q)
	}

	// Output:
	// [s1 p1 o1 g1]
	// [s1 p2 o2 g1]
	// [s2 p2 o2 g1]
	// [s2 p2 o3 g2]
	// [s3 p3 o3 g2]
}

func ExampleQuadStore_FindGraphs() {

	quads := [][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	}

	s := store4.NewQuadStore()
	for _, q := range quads {
		s.Add(q[0], q[1], q[2], q[3])
	}

	// Get a list of all graphs.
	results := s.FindGraphs("*", "*", "*")
	// (We only sort the results before printing
	// because iteration order is not stable)
	sortStrings(results)
	fmt.Println(results)

	// Find graphs containing quads that have subject s1.
	results = s.FindGraphs("s1", "*", "*")
	sortStrings(results)
	fmt.Println(results)

	// Find graphs containing quads that have
	// both subject s2 and predicate p2.
	results = s.FindGraphs("s2", "p2", "*")
	sortStrings(results)
	fmt.Println(results)

	// Find graphs containging quads that have object o3.
	results = s.FindGraphs("*", "*", "o3")
	sortStrings(results)
	fmt.Println(results)

	// Output:
	// [g1 g2]
	// [g1]
	// [g1 g2]
	// [g2]
}

func ExampleQuadStore_FindSubjects() {

	quads := [][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	}

	s := store4.NewQuadStore()
	for _, q := range quads {
		s.Add(q[0], q[1], q[2], q[3])
	}

	// Get a list of all subjects in the store.
	results := s.FindSubjects("*", "*", "*")
	// (We only sort the results before printing
	// because iteration order is not stable)
	sortStrings(results)
	fmt.Println(results)

	// Find all subjects in graph g2.
	results = s.FindSubjects("*", "*", "g2")
	sortStrings(results)
	fmt.Println(results)

	// Find subjects for quads that have
	// both predicate p2 and object o2.
	results = s.FindSubjects("p2", "o2", "*")
	sortStrings(results)
	fmt.Println(results)

	// Output:
	// [s1 s2 s3]
	// [s2 s3]
	// [s1 s2]
}
