package store4_test

import (
	. "github.com/jimsmart/store4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Triple struct {
	S string
	P string
	O interface{}
}

func (t Triple) Subject() string     { return t.S }
func (t Triple) Predicate() string   { return t.P }
func (t Triple) Object() interface{} { return t.O }

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
				var resultsList []*Triple
				graph.ForEach(func(s, p string, o interface{}) {
					resultsList = append(resultsList, &Triple{s, p, o})
				})
				Expect(resultsList).To(ConsistOf([]*Triple{
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
				var resultsList []*Triple
				graph.ForEach(func(s, p string, o interface{}) {
					resultsList = append(resultsList, &Triple{s, p, o})
				})
				Expect(resultsList).To(ConsistOf([]*Triple{
					{"s1", "p1", "o1"},
				}))
			})
		})

		Context("from an unknown type", func() {
			var g *GraphView
			shouldPanic := func() {
				g = NewGraph(true)
			}

			It("should panic", func() {
				Expect(shouldPanic).To(Panic())
			})
		})

	})
})
