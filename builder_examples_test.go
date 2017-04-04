package store4_test

import (
	"fmt"

	"github.com/jimsmart/store4"
)

func ExampleBuilder() {

	b := store4.Builder{}

	// Set the subject to 's1', and add two statements.
	b.Subject("s1").
		Add("p1", "o1").
		Add("p2", "o2")

	// Change the subject to 's2', and add two statements.
	b.Subject("s2").
		Add("p1", "o3").
		Add("p2", "o4")

	// Change to graph 'g1'.
	b.Graph("g1")
	// Set the subject to 's3' and add a statement.
	b.Subject("s3").Add("p3", "o5")

	s := b.Build()
	fmt.Println(s)

	// Output:
	// [s1 p1 o1 ]
	// [s1 p2 o2 ]
	// [s2 p1 o3 ]
	// [s2 p2 o4 ]
	// [s3 p3 o5 g1]
}
