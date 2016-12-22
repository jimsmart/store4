package store4_test

import (
	. "github.com/jimsmart/store4"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("util functions", func() {

	Describe("SortQuads", func() {

		quads := [][4]string{
			{"s3", "p3", "o5", "g3"},
			{"s2", "p2", "o4", "g2"},
			{"s2", "p2", "o3", "g1"},
			{"s1", "p1", "o2", "g1"},
			{"s1", "p0", "o1", "g1"},
			{"s1", "p1", "o1", "g1"},
		}

		It("should not prevent 100% coverage", func() {
			SortQuads(quads)
		})
	})

	Describe("SortTriples", func() {

		triples := [][3]string{
			{"s3", "p3", "o5"},
			{"s2", "p2", "o4"},
			{"s2", "p2", "o3"},
			{"s1", "p1", "o2"},
			{"s1", "p0", "o1"},
			{"s1", "p1", "o1"},
		}

		It("should not prevent 100% coverage", func() {
			SortTriples(triples)
		})
	})

	Describe("SortTuples", func() {

		tuples := [][2]string{
			{"p3", "o5"},
			{"p2", "o4"},
			{"p2", "o3"},
			{"p1", "o2"},
			{"p0", "o1"},
			{"p1", "o1"},
		}

		It("should not prevent 100% coverage", func() {
			SortTuples(tuples)
		})
	})

})
