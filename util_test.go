package store4

import

// . "github.com/jimsmart/store4"

. "github.com/onsi/ginkgo"

// . "github.com/onsi/gomega"

var _ = Describe("util functions", func() {

	// Describe("SortQuads", func() {

	// 	quads := []*Quad{
	// 		{"s3", "p3", "o5", "g3"},
	// 		{"s2", "p2", "o4", "g2"},
	// 		{"s2", "p2", "o3", "g1"},
	// 		{"s1", "p1", "o2", "g1"},
	// 		{"s1", "p0", "o1", "g1"},
	// 		{"s1", "p1", "o1", "g1"},
	// 	}

	// 	It("should not prevent 100% coverage", func() {
	// 		SortQuads(quads)
	// 	})
	// })

	// Describe("SortTriples", func() {

	// 	triples := []*Triple{
	// 		{"s3", "p3", "o5"},
	// 		{"s2", "p2", "o4"},
	// 		{"s2", "p2", "o3"},
	// 		{"s1", "p1", "o2"},
	// 		{"s1", "p0", "o1"},
	// 		{"s1", "p1", "o1"},
	// 	}

	// 	It("should not prevent 100% coverage", func() {
	// 		SortTriples(triples)
	// 	})
	// })

	// Describe("SortTuples", func() {

	// 	tuples := []*Tuple{
	// 		{"p3", "o5"},
	// 		{"p2", "o4"},
	// 		{"p2", "o3"},
	// 		{"p1", "o2"},
	// 		{"p0", "o1"},
	// 		{"p1", "o1"},
	// 	}

	// 	It("should not prevent 100% coverage", func() {
	// 		SortTuples(tuples)
	// 	})
	// })

	Describe("sortObjects", func() {

		objects := []interface{}{
			"o5",
			"o4",
			"o3",
			"o2",
			"o1",
			"o1",
		}

		It("should not prevent 100% coverage", func() {
			sortObjects(objects)
		})
	})

})
