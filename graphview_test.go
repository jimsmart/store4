package store4_test

import (
	. "github.com/jimsmart/store4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GraphView", func() {

	Describe("Creating a new GraphView", func() {

		Context("from [][3]string", func() {
			graph := NewGraph([][3]string{
				{"s1", "p1", "o1"},
				{"s1", "p1", "o2"},
				{"s1", "p2", "o2"},
				{"s2", "p1", "o1"},
				{"s1", "p2", "o3"},
			})

			It("should have size 5", func() {
				Expect(graph.Size()).To(Equal(uint64(5)))
			})

			It("should contain the correct triples", func() {
				var resultsList [][3]string
				graph.ForEach(func(s, p, o string) {
					resultsList = append(resultsList, [3]string{s, p, o})
				})
				Expect(resultsList).To(ConsistOf([][3]string{
					{"s1", "p1", "o1"},
					{"s1", "p1", "o2"},
					{"s1", "p2", "o2"},
					{"s2", "p1", "o1"},
					{"s1", "p2", "o3"},
				}))
			})
		})

		Context("from a single [3]string", func() {
			graph := NewGraph([3]string{"s1", "p1", "o1"})

			It("should have size 1", func() {
				Expect(graph.Size()).To(Equal(uint64(1)))
			})

			It("should contain the correct triple", func() {
				var resultsList [][3]string
				graph.ForEach(func(s, p, o string) {
					resultsList = append(resultsList, [3]string{s, p, o})
				})
				Expect(resultsList).To(ConsistOf([][3]string{
					{"s1", "p1", "o1"},
				}))
			})
		})

		Context("from an unknown type", func() {
			shouldPanic := func() {
				NewGraph(true)
			}

			It("should panic", func() {
				Expect(shouldPanic).To(Panic())
			})
		})

	})
})
