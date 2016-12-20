package store4_test

import (
	"fmt"

	"github.com/jimsmart/store4"
)

func ExampleGraph_Query() {

	g := store4.NewGraph([][3]string{
		{"s1", "p1", "o1"},
		{"s1", "p2", "o2"},
		{"s2", "p2", "o2"},
		{"s2", "p2", "o3"},
		{"s3", "p3", "o3"},
	})

	// Query for projections over all subjects
	// that have p1=o1 and p2=o2.
	pattern := map[string]string{
		"p1": "o1",
		"p2": "o2",
	}
	results := g.Query(pattern)

	fmt.Println(len(results))
	fmt.Println(results[0].Subject)
	// Output:
	// 1
	// s1
}

func ExampleQuadStore_Query() {

	s := store4.NewQuadStore([][4]string{
		{"s1", "p1", "o1", "g1"},
		{"s1", "p2", "o2", "g1"},
		{"s2", "p2", "o2", "g1"},
		{"s2", "p2", "o3", "g2"},
		{"s3", "p3", "o3", "g2"},
	})

	// Query for projections over all subjects in graph g1
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
