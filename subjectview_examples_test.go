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

func ExampleSubjectView_Count() {

	g := store4.NewGraph([][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s1", "p2", "o1"},
	})

	// Get a SubjectView of subject s1.
	s1 := g.SubjectView("s1")

	// How many tuples have predicate p1?
	count := s1.Count("p1", "*")
	fmt.Println(count)

	// How many tuples have predicate p2?
	count = s1.Count("p2", "*")
	fmt.Println(count)

	// How many tuples have object o1?
	count = s1.Count("*", "o1")
	fmt.Println(count)

	// Output:
	// 1
	// 2
	// 2
}

func ExampleSubjectView_Empty() {

	g := store4.NewGraph([][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
	})

	// Get a SubjectView of subject s1.
	s1 := g.SubjectView("s1")
	// Is it empty?
	fmt.Println(s1.Empty())

	// Get a SubjectView of subject s3.
	s3 := g.SubjectView("s3")
	// Is it empty?
	fmt.Println(s3.Empty())

	// Output:
	// false
	// true
}

func ExampleSubjectView_Every() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p1", "o2", "g1"},
		{"s2", "p2", "o2", "g2"},
	})

	v1 := s.SubjectView("s1", "g1")
	v2 := s.SubjectView("s2", "g2")

	p1TestFn := func(p, o string) bool {
		return p == "p1"
	}

	// Iterate over every tuple in the SubjectView
	// while true is being returned from our
	// callback, and halt when false is returned.
	// Returns false if the callback ever returned false.
	result := v1.Every(p1TestFn)
	fmt.Println(result)
	result = v2.Every(p1TestFn)
	fmt.Println(result)

	// Note that Every will return true
	// for an empty SubjectView.
	v0 := s.SubjectView("s0", "g0")
	result = v0.Every(p1TestFn)
	fmt.Println(result)

	// Every is often used as a breakable iterator,
	// with its return value being ignored.

	// Output:
	// true
	// false
	// true
}

func ExampleSubjectView_EveryWith() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p1", "o2", "g1"},
		{"s2", "p2", "o2", "g2"},
	})

	v1 := s.SubjectView("s1", "g1")
	v2 := s.SubjectView("s2", "g2")

	o2TestFn := func(p, o string) bool {
		return p == "o2"
	}

	// Iterate over every matching tuple in the SubjectView
	// while true is being returned from our
	// callback, and halt when false is returned.
	// Returns false if the callback ever returned false.
	result := v1.EveryWith("p1", "*", o2TestFn)
	fmt.Println(result)
	result = v2.EveryWith("p2", "*", o2TestFn)
	fmt.Println(result)

	// Note that EveryWith will return true
	// for an empty SubjectView....
	v0 := s.SubjectView("s0", "g0")
	result = v0.EveryWith("*", "*", o2TestFn)
	fmt.Println(result)
	// ...or if its iteration set is empty.
	result = v1.EveryWith("p0", "o0", o2TestFn)

	// Every is often used as a breakable iterator,
	// with its return value being ignored.

	// Output:
	// false
	// false
	// true
}

func ExampleSubjectView_FindPredicates() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})

	v1 := s.SubjectView("s1", "g1")

	// Get a list of all predicates in the SubjectView.
	results := v1.FindPredicates("*")
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results)
	fmt.Println(results)

	// Find predicates in the SubjectView that have
	// object o2.
	results = v1.FindPredicates("o2")
	sort.Strings(results)
	fmt.Println(results)

	// Output:
	// [p1 p2]
	// [p2]
}

func ExampleSubjectView_FindObjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})

	v1 := s.SubjectView("s1", "g1")

	// Get a list of all objects in the SubjectView.
	results := v1.FindObjects("*")
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results)
	fmt.Println(results)

	// Find objects in the SubjectView
	// that have predicate p2.
	results = v1.FindObjects("p2")
	sort.Strings(results)
	fmt.Println(results)

	// Output:
	// [o1 o2]
	// [o2]
}

func ExampleSubjectView_ForEach() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})

	v := s.SubjectView("s2", "g1")

	// Iterate over all predicate-object tuples
	// in the SubjectView for subject s2 in graph g1.
	v.ForEach(func(p, o string) {
		fmt.Println(p, o)
	})

	// Output:
	// p2 o2
}

func ExampleSubjectView_ForEachWith() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})

	v := s.SubjectView("s1", "g1")

	// Iterate over all predicate-object tuples
	// that have predicate p1
	// for the SubjectView of subject s1 in graph g1.
	v.ForEachWith("p1", "*", func(p, o string) {
		fmt.Println(p, o)
	})

	// Output:
	// p1 o1
}

func ExampleSubjectView_Map() {

	g := store4.NewGraph([][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
		{"s2", "p2", "o3"},
		{"s3", "p3", "o3"},
	})

	// Get a SubjectView of subject s1.
	s1 := g.SubjectView("s1")
	// Get subject's values as a map,
	// having the predicate value as the key,
	// and having a slice of object values as the value.
	m := s1.Map()
	p1 := m["p1"]
	fmt.Println(p1)
	p2 := m["p2"]
	fmt.Println(p2)

	// Output:
	// [o1]
	// [o2]
}
