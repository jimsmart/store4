package store4_test

import (
	. "github.com/jimsmart/store4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubjectView", func() {

	Describe("Creating a single SubjectView from a QuadStore", func() {

		Context("with an existing subject and wildcard graph", func() {
			store := NewQuadStore([][4]string{
				{"s1", "p1", "o1", ""},
				{"s1", "p1", "o2", ""},
				{"s1", "p2", "o2", ""},
				{"s2", "p1", "o1", ""},
				{"s1", "p2", "o3", "c4"},
			})

			view := store.SubjectView("s1", "*")

			It("should have size 4", func() {
				Expect(view.Size()).To(Equal(uint64(4)))
			})

			It("should contain the correct predicates", func() {
				Expect(view.FindPredicates("*")).To(ConsistOf([]string{"p1", "p2"}))
			})

			It("should contain the correct object values for predicate p1", func() {
				Expect(view.FindObjects("p1")).To(ConsistOf([]string{"o1", "o2"}))
			})

			It("should contain the correct object values for predicate p2", func() {
				Expect(view.FindObjects("p2")).To(ConsistOf([]string{"o2", "o3"}))
			})
		})

	})

})
