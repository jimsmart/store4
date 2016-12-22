package store4_test

import (
	"fmt"
	"sort"

	"github.com/jimsmart/store4"
)

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

func ExampleQuadStore_SubjectViews() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Get SubjectViews for all subjects in the store.
	results1 := s.SubjectViews("*", "*", "*")
	fmt.Println(len(results1))
	// (Exact order may vary)
	//fmt.Println(results1[0])
	//fmt.Println(results1[1])

	// Get SubjectViews for all subjects featuring
	// predicate p2 in graph g2.
	results2 := s.SubjectViews("p2", "*", "g2")
	fmt.Println(len(results2))
	fmt.Println(results2[0])

	// Output:
	// 3
	// 1
	// g2
	// s2
	// [p2 o3]
}

func ExampleGraphView_SubjectViews() {

	g := store4.NewGraph([][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
		{"s2", "p2", "o3"},
		{"s3", "p3", "o3"},
	})

	// Get SubjectViews for all subjects in the graph.
	results1 := g.SubjectViews("*", "*")
	fmt.Println(len(results1))
	// (Exact order may vary)
	//fmt.Println(results1[0])
	//fmt.Println(results1[1])

	// Get SubjectViews for all subjects in the graph
	// featuring predicate p3.
	results2 := g.SubjectViews("p3", "*")
	fmt.Println(len(results2))
	fmt.Println(results2[0])

	// Output:
	// 3
	// 1
	// s3
	// [p3 o3]
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

func ExampleSubjectView_ForPredicates() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	v := s.SubjectView("s1", "g1")

	var results1 []string
	// Iterate over all predicates for subject s1 in graph g1.
	v.ForPredicates("*", func(p string) {
		results1 = append(results1, p)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results1)
	fmt.Println(results1)

	var results2 []string
	// Iterate over all predicates for subject s1 in graph g1
	// that have object o1.
	v.ForPredicates("o1", func(p string) {
		results2 = append(results2, p)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results2)
	fmt.Println(results2)

	// Output:
	// [p1 p2]
	// [p1]
}

func ExampleSubjectView_ForObjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	v := s.SubjectView("s1", "g1")

	var results1 []string
	// Iterate over all objects for subject s1 in graph g1.
	v.ForObjects("*", func(o string) {
		results1 = append(results1, o)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results1)
	fmt.Println(results1)

	var results2 []string
	// Iterate over objects for subject s1 in graph g1
	// that have predicate p2.
	v.ForObjects("p2", func(o string) {
		results2 = append(results2, o)
	})
	sort.Strings(results2)
	fmt.Println(results2)

	// Output:
	// [o1 o2]
	// [o2]
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

func ExampleSubjectView_Remove() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	v := s.SubjectView("s1", "g1")
	fmt.Println(v.Size())

	// Remove a specific tuple from the view.
	v.Remove("p2", "o2")
	fmt.Println(v.Size())

	// Remove all tuples that have predicate p1.
	v.Remove("p1", "*")
	fmt.Println(v.Size())

	fmt.Println(s.Size())

	// Output:
	// 2
	// 1
	// 0
	// 3
}

func ExampleSubjectView_Some() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	v1 := s.SubjectView("s1", "g1")
	v2 := s.SubjectView("s2", "*")

	o1TestFn := func(p, o string) bool {
		return o == "o1"
	}

	// Is there some tuple in the view with object o1?
	result := v1.Some(o1TestFn)
	fmt.Println(result)

	// Is there some tuple in the other view with object o1?
	result = v2.Some(o1TestFn)
	fmt.Println(result)

	// Output:
	// true
	// false
}

func ExampleSubjectView_SomeWith() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	alwaysTrueFn := func(p, o string) bool {
		return true
	}

	v1 := s.SubjectView("s1", "g1")
	v2 := s.SubjectView("s2", "*")

	// Is there some tuple in the view with object o1?
	result := v1.SomeWith("*", "o1", alwaysTrueFn)
	fmt.Println(result)

	// Is there some tuple in the other view with predicate p2?
	result = v2.SomeWith("p2", "*", alwaysTrueFn)
	fmt.Println(result)

	// Output:
	// true
	// true
}
