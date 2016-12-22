package store4_test

import (
	"fmt"
	"sort"

	"github.com/jimsmart/store4"
)

// TODO(js) Write an example for QuadStore.SubjectViews.

func ExampleNewQuadStore() {

	// A new empty store.
	s1 := store4.NewQuadStore()
	fmt.Println(s1)

	// A new store initialised with a slice of quads.
	s2 := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})
	fmt.Println(s2)

	// A new store initialised with a slice of triples.
	s3 := store4.NewQuadStore([][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
	})
	fmt.Println(s3)

	triples := [][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
	}
	quads := [][4]string{
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	}
	// A new store initialised with some quads and some triples.
	s4 := store4.NewQuadStore(triples, quads)
	fmt.Println(s4)

	// Output:
	// [s1 p1 o1 g1]
	// [s1 p2 o2 g1]
	// [s2 p2 o2 g1]
	//
	// [s1 p1 o1 ]
	// [s1 p2 o2 ]
	// [s2 p2 o2 ]
	//
	// [s1 p1 o1 ]
	// [s1 p2 o2 ]
	// [s2 p2 o2 ]
	// [s2 p2 o3 g2]
	// [s3 p3 o3 g2]
}

func ExampleQuadStore_Add() {

	s := store4.NewQuadStore()

	quads := [][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	}

	for _, q := range quads {
		// Add quads to the store, testing for duplicates.
		ok := s.Add(q[0], q[1], q[2], q[3])
		if !ok {
			panic("duplicate quad")
		}
	}

	fmt.Println(s)
	// Output:
	// [s1 p1 o1 g1]
	// [s1 p2 o2 g1]
	// [s2 p2 o2 g1]
	// [s2 p2 o3 g2]
	// [s3 p3 o3 g2]
}

func ExampleQuadStore_Count() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

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

func ExampleQuadStore_Empty() {
	s1 := store4.NewQuadStore()
	fmt.Println(s1.Empty())

	s2 := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})
	fmt.Println(s2.Empty())
	// Output:
	// true
	// false
}

func ExampleQuadStore_Every() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})

	g1TestFn := func(s, p, o, g string) bool {
		return g == "g1"
	}

	// Iterate over every quad in the store
	// while true is being returned from our
	// callback, and halt when false is returned.
	// Returns false if the callback ever returned false.
	result := s.Every(g1TestFn)
	fmt.Println(result)

	// Note that Every will return true
	// for an empty store.
	s0 := store4.NewQuadStore()
	result = s0.Every(g1TestFn)
	fmt.Println(result)

	// Every is often simply used as a breakable iterator,
	// with its return value being ignored.

	// Output:
	// true
	// true
}

func ExampleQuadStore_EveryWith() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})

	g1TestFn := func(s, p, o, g string) bool {
		return g == "g1"
	}

	// Iterate over every quad having subject s1
	// while true is being returned from our
	// callback, and halt when false is returned.
	// Returns false if the callback ever returned false.
	result := s.EveryWith("s1", "*", "*", "*", g1TestFn)
	fmt.Println(result)

	// Note that EveryWith will return true
	// for an empty store...
	s0 := store4.NewQuadStore()
	result = s0.EveryWith("*", "*", "*", "*", g1TestFn)
	fmt.Println(result)
	// ...or if its iteration set is empty.
	result = s.EveryWith("s0", "*", "*", "*", g1TestFn)
	fmt.Println(result)

	// EveryWith is often simply used as a breakable iterator,
	// with its return value being ignored.

	// Output:
	// true
	// true
	// true
}

func ExampleQuadStore_FindGraphs() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Get a list of all graphs.
	results := s.FindGraphs("*", "*", "*")
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results)
	fmt.Println(results)

	// Find graphs containing quads that have subject s1.
	results = s.FindGraphs("s1", "*", "*")
	sort.Strings(results)
	fmt.Println(results)

	// Find graphs containing quads that have
	// both subject s2 and predicate p2.
	results = s.FindGraphs("s2", "p2", "*")
	sort.Strings(results)
	fmt.Println(results)

	// Find graphs containging quads that have object o3.
	results = s.FindGraphs("*", "*", "o3")
	sort.Strings(results)
	fmt.Println(results)

	// Output:
	// [g1 g2]
	// [g1]
	// [g1 g2]
	// [g2]
}

func ExampleQuadStore_FindSubjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Get a list of all subjects in the store.
	results := s.FindSubjects("*", "*", "*")
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results)
	fmt.Println(results)

	// Find all subjects in graph g2.
	results = s.FindSubjects("*", "*", "g2")
	sort.Strings(results)
	fmt.Println(results)

	// Find subjects for quads that have
	// both predicate p2 and object o2.
	results = s.FindSubjects("p2", "o2", "*")
	sort.Strings(results)
	fmt.Println(results)

	// Output:
	// [s1 s2 s3]
	// [s2 s3]
	// [s1 s2]
}

func ExampleQuadStore_FindPredicates() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Get a list of all predicates in the store.
	results := s.FindPredicates("*", "*", "*")
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results)
	fmt.Println(results)

	// Find all predicates in graph g1.
	results = s.FindPredicates("*", "*", "g1")
	sort.Strings(results)
	fmt.Println(results)

	// Find predicates for quads that have
	// both subject s1 and object o2.
	results = s.FindPredicates("s1", "o2", "*")
	sort.Strings(results)
	fmt.Println(results)

	// Output:
	// [p1 p2 p3]
	// [p1 p2]
	// [p2]
}

func ExampleQuadStore_FindObjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Get a list of all objects in the store.
	results := s.FindObjects("*", "*", "*")
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results)
	fmt.Println(results)

	// Find all objects in graph g1.
	results = s.FindObjects("*", "*", "g1")
	sort.Strings(results)
	fmt.Println(results)

	// Find objects for quads that have
	// both subject s1 and predicate p2.
	results = s.FindObjects("s2", "p2", "*")
	sort.Strings(results)
	fmt.Println(results)

	// Output:
	// [o1 o2 o3]
	// [o1 o2]
	// [o2 o3]
}

func ExampleQuadStore_ForEach() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Iterate over all quads.
	var results [][4]string
	s.ForEach(func(s, p, o, g string) {
		results = append(results, [4]string{s, p, o, g})
	})

	// (We only sort the results before printing
	// because iteration order is unstable)
	store4.SortQuads(results)
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

func ExampleQuadStore_ForEachWith() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Iterate over all quads with predicate p2.
	var results [][4]string
	s.ForEachWith("*", "p2", "*", "*", func(s, p, o, g string) {
		results = append(results, [4]string{s, p, o, g})
	})

	// (We only sort the results before printing
	// because iteration order is unstable)
	store4.SortQuads(results)
	for _, q := range results {
		fmt.Println(q)
	}

	// Output:
	// [s1 p2 o2 g1]
	// [s2 p2 o2 g1]
	// [s2 p2 o3 g2]
}

func ExampleQuadStore_ForGraphs() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	var results1 []string
	// Iterate over all graphs.
	s.ForGraphs("*", "*", "*", func(g string) {
		results1 = append(results1, g)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results1)
	fmt.Println(results1)

	var results2 []string
	// Iterate over graphs containing quads that
	// have subject s1.
	s.ForGraphs("s1", "*", "*", func(g string) {
		results2 = append(results2, g)
	})
	sort.Strings(results2)
	fmt.Println(results2)

	var results3 []string
	// Iterate over graphs containing quads that
	// have both subject s2 and predicate p2.
	s.ForGraphs("s2", "p2", "*", func(g string) {
		results3 = append(results3, g)

	})
	sort.Strings(results3)
	fmt.Println(results3)

	var results4 []string
	// Iterate over graphs containging quads that have object o3.
	s.ForGraphs("*", "*", "o3", func(g string) {
		results4 = append(results4, g)
	})
	sort.Strings(results4)
	fmt.Println(results4)

	// Output:
	// [g1 g2]
	// [g1]
	// [g1 g2]
	// [g2]
}

func ExampleQuadStore_ForSubjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	var results1 []string
	// Iterate over all subjects in the store.
	s.ForSubjects("*", "*", "*", func(s string) {
		results1 = append(results1, s)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results1)
	fmt.Println(results1)

	var results2 []string
	// Iterate over all subjects in graph g2.
	s.ForSubjects("*", "*", "g2", func(s string) {
		results2 = append(results2, s)
	})
	sort.Strings(results2)
	fmt.Println(results2)

	var results3 []string
	// Iterate over subjects for quads that have
	// both predicate p2 and object o2.
	s.ForSubjects("p2", "o2", "*", func(s string) {
		results3 = append(results3, s)
	})
	sort.Strings(results3)
	fmt.Println(results3)

	// Output:
	// [s1 s2 s3]
	// [s2 s3]
	// [s1 s2]
}

func ExampleQuadStore_ForPredicates() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	var results1 []string
	// Iterate over all predicates in the store.
	s.ForPredicates("*", "*", "*", func(p string) {
		results1 = append(results1, p)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results1)
	fmt.Println(results1)

	var results2 []string
	// Iterate over all predicates in graph g1.
	s.ForPredicates("*", "*", "g1", func(p string) {
		results2 = append(results2, p)
	})
	sort.Strings(results2)
	fmt.Println(results2)

	var results3 []string
	// Iterate over predicates for quads that have
	// both subject s1 and object o2.
	s.ForPredicates("s1", "o2", "*", func(p string) {
		results3 = append(results3, p)
	})
	sort.Strings(results3)
	fmt.Println(results3)

	// Output:
	// [p1 p2 p3]
	// [p1 p2]
	// [p2]
}

func ExampleQuadStore_ForObjects() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	var results1 []string
	// Iterate over all objects in the store.
	s.ForObjects("*", "*", "*", func(o string) {
		results1 = append(results1, o)
	})
	// (We only sort the results before printing
	// because iteration order is unstable)
	sort.Strings(results1)
	fmt.Println(results1)

	var results2 []string
	// Iterate over all objects in graph g1.
	s.ForObjects("*", "*", "g1", func(o string) {
		results2 = append(results2, o)
	})
	sort.Strings(results2)
	fmt.Println(results2)

	var results3 []string
	// Iterate over objects for quads that have
	// both subject s1 and predicate p2.
	s.ForObjects("s2", "p2", "*", func(o string) {
		results3 = append(results3, o)
	})
	sort.Strings(results3)
	fmt.Println(results3)

	// Output:
	// [o1 o2 o3]
	// [o1 o2]
	// [o2 o3]
}

func ExampleQuadStore_Remove() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	fmt.Println(s.Size())

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
	// 5
	// 4
	// 3
	// 0
}

func ExampleQuadStore_Size() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// How many quads are in the store?
	count := s.Size()
	fmt.Println(count)

	// Output: 5
}

func ExampleQuadStore_String() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// (Println calls String)
	fmt.Println(s)

	// Output:
	// [s1 p1 o1 g1]
	// [s1 p2 o2 g1]
	// [s2 p2 o2 g1]
	// [s2 p2 o3 g2]
	// [s3 p3 o3 g2]
}

func ExampleQuadStore_Some() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})

	s2TestFn := func(s, p, o, g string) bool {
		return s == "s2"
	}

	g3TestFn := func(s, p, o, g string) bool {
		return g == "g3"
	}

	// Is there some quad with subject s2?
	result := s.Some(s2TestFn)
	fmt.Println(result)

	// Is there some quad in graph g3?
	result = s.Some(g3TestFn)
	fmt.Println(result)

	// Output:
	// true
	// false
}

func ExampleQuadStore_SomeWith() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
	})

	alwaysTrueFn := func(s, p, o, g string) bool {
		return true
	}

	// Is there some quad with subject s1?
	result := s.SomeWith("s1", "*", "*", "*", alwaysTrueFn)
	fmt.Println(result)

	// Output:
	// true
}
