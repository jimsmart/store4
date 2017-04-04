package store4_test

import (
	. "github.com/jimsmart/store4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Builder", func() {

	Context("A new unused builder", func() {
		b := &Builder{}
		Describe("Build", func() {
			s := b.Build()
			It("should return an empty quad store", func() {
				Expect(s.Size()).To(BeZero())
			})
		})
	})

	Describe("Building some graphs", func() {
		b := &Builder{}
		b.Graph("g1")
		b.Subject("s1").
			Add("p1", "o1").
			Add("p2", "o2")
		b.Graph("g2").Subject("s2").Add("p3", "o3")
		b.DefaultGraph().
			Subject("s3").
			Add("p4", "o4")
		store := b.Build()
		It("should return a store containing the correct quads", func() {
			resultsList := iterResults(store)
			Expect(resultsList).To(ConsistOf([]*Quad{
				{"s1", "p1", "o1", "g1"},
				{"s1", "p2", "o2", "g1"},
				{"s2", "p3", "o3", "g2"},
				{"s3", "p4", "o4", ""},
			}))
		})
	})

	Describe("Error conditions", func() {
		Context("Add", func() {
			It("should panic if the subject is missing", func() {
				b := &Builder{}
				b.Graph("g1")
				// b.Subject("s1")
				Expect(func() { b.Add("p1", "o1") }).To(Panic())
			})
			It("should panic if the predicate is missing", func() {
				b := &Builder{}
				b.Graph("g1")
				b.Subject("s1")
				Expect(func() { b.Add("", "o1") }).To(Panic())
			})
			It("should panic if the object is missing", func() {
				b := &Builder{}
				b.Graph("g1")
				b.Subject("s1")
				Expect(func() { b.Add("p1", nil) }).To(Panic())
			})
		})
	})

})
