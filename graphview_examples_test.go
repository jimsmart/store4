package store4_test

import (
	"fmt"

	"github.com/jimsmart/store4"
)

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

	// Output:
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
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
	// Output:
	//
	// [s1 p1 o1]
	// [s1 p2 o2]
	// [s2 p2 o2]
	// [s2 p2 o3]
	// [s3 p3 o3]
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
