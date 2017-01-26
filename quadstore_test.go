package store4_test

import (
	. "github.com/jimsmart/store4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Note: There are currently no tests for:
// SomeWith, ForGraphs, ForSubjects, ForPredicates, ForObjects.
//
// This is because the current implementation of QuadStore makes
// internal usage of these methods - so: if the calling methods
// pass all tests, then the called methods also pass.
// (We achieve 100% coverage already anyway)
//
// This exact setup may not be the case forever.

// TODO(js) Write tests for SomeWith, ForGraphs, ForSubjects, ForPredicates, ForObjects
// - instead of relying on them being called internally.

func alwaysTrueFn(s, p string, o interface{}, g string) bool {
	return true
}
func alwaysFalseFn(s, p string, o interface{}, g string) bool {
	return false
}

type Quad struct {
	S string
	P string
	O interface{}
	G string
}

func (q Quad) Subject() string     { return q.S }
func (q Quad) Predicate() string   { return q.P }
func (q Quad) Object() interface{} { return q.O }
func (q Quad) Graph() string       { return q.G }

type boxedQuad2 struct{ s, p, o, g string }

func (b boxedQuad2) S() string      { return b.s }
func (b boxedQuad2) P() string      { return b.p }
func (b boxedQuad2) O() interface{} { return b.o }
func (b boxedQuad2) G() string      { return b.g }

var _ = Describe("QuadStore", func() {

	iterResults := func(store *QuadStore) []*Quad {
		var resultsList []*Quad
		store.ForEach(func(s, p string, o interface{}, g string) {
			resultsList = append(resultsList, &Quad{s, p, o, g})
		})
		return resultsList
	}

	Describe("Creating a new QuadStore", func() {

		Context("from [][4]string", func() {
			store := NewQuadStore([][4]string{
				{"s1", "p1", "o1", ""},
				{"s1", "p1", "o2", ""},
				{"s1", "p2", "o2", ""},
				{"s2", "p1", "o1", ""},
				{"s1", "p2", "o3", "c4"},
			})
			It("should have size 5", func() {
				Expect(store.Size()).To(Equal(uint64(5)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", "c4"},
				}))
			})
		})

		Context("from multiple [][4]string", func() {
			store := NewQuadStore(
				[][4]string{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
				},
				[][4]string{
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", "c4"},
				},
			)
			It("should have size 5", func() {
				Expect(store.Size()).To(Equal(uint64(5)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", "c4"},
				}))
			})
		})

		Context("from multiple [4]string", func() {
			store := NewQuadStore(
				[4]string{"s1", "p1", "o1", ""},
				[4]string{"s1", "p1", "o2", ""},
				[4]string{"s1", "p2", "o2", ""},
				[4]string{"s2", "p1", "o1", ""},
				[4]string{"s1", "p2", "o3", "c4"},
			)
			It("should have size 5", func() {
				Expect(store.Size()).To(Equal(uint64(5)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", "c4"},
				}))
			})
		})

		Context("from [][3]string", func() {
			store := NewQuadStore([][3]string{
				{"s1", "p1", "o1"},
				{"s1", "p1", "o2"},
				{"s1", "p2", "o2"},
				{"s2", "p1", "o1"},
				{"s1", "p2", "o3"},
			})
			It("should have size 5", func() {
				Expect(store.Size()).To(Equal(uint64(5)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", ""},
				}))
			})
		})

		Context("from multiple [3]string", func() {
			store := NewQuadStore(
				[3]string{"s1", "p1", "o1"},
				[3]string{"s1", "p1", "o2"},
				[3]string{"s1", "p2", "o2"},
				[3]string{"s2", "p1", "o1"},
				[3]string{"s1", "p2", "o3"},
			)
			It("should have size 5", func() {
				Expect(store.Size()).To(Equal(uint64(5)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", ""},
				}))
			})
		})

		Context("from a quad-like struct", func() {
			store := NewQuadStore(
				struct {
					Subject   string
					Predicate string
					Object    string
					Graph     string
				}{"s1", "p1", "o1", ""},
			)
			It("should have size 1", func() {
				Expect(store.Size()).To(Equal(uint64(1)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
				}))
			})
		})
		Context("from a quad-like struct with a non-string object type", func() {
			store := NewQuadStore(
				struct {
					Subject   string
					Predicate string
					Object    interface{}
					Graph     string
				}{"s1", "p1", "o1", ""},
			)
			It("should have size 1", func() {
				Expect(store.Size()).To(Equal(uint64(1)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
				}))
			})
		})
		Context("from a pointer to a quad-like struct", func() {
			store := NewQuadStore(
				&struct {
					Subject   string
					Predicate string
					Object    string
					Graph     string
				}{"s1", "p1", "o1", ""},
			)
			It("should have size 1", func() {
				Expect(store.Size()).To(Equal(uint64(1)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
				}))
			})
		})

		Context("from a simpler quad-like struct", func() {
			store := NewQuadStore(
				struct {
					S string
					P string
					O string
					G string
				}{"s1", "p1", "o1", ""},
			)
			It("should have size 1", func() {
				Expect(store.Size()).To(Equal(uint64(1)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
				}))
			})
		})

		Context("from a triple-like struct", func() {
			store := NewQuadStore(
				struct {
					Subject   string
					Predicate string
					Object    string
				}{"s1", "p1", "o1"},
			)
			It("should have size 1", func() {
				Expect(store.Size()).To(Equal(uint64(1)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
				}))
			})
		})

		Context("from a simpler triple-like struct", func() {
			store := NewQuadStore(
				struct {
					S string
					P string
					O string
				}{"s1", "p1", "o1"},
			)
			It("should have size 1", func() {
				Expect(store.Size()).To(Equal(uint64(1)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
				}))
			})
		})

		Context("from multiple quad-like structs", func() {
			store := NewQuadStore(
				[]struct {
					Subject   string
					Predicate string
					Object    string
					Graph     string
				}{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", ""},
				})
			It("should have size 5", func() {
				Expect(store.Size()).To(Equal(uint64(5)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", ""},
				}))
			})
		})

		Context("from a Quad struct", func() {
			store := NewQuadStore(
				Quad{"s1", "p1", "o1", ""},
			)
			It("should have size 1", func() {
				Expect(store.Size()).To(Equal(uint64(1)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
				}))
			})
		})

		Context("from a simpler boxed-quad struct", func() {
			store := NewQuadStore(
				boxedQuad2{"s1", "p1", "o1", ""},
			)
			It("should have size 1", func() {
				Expect(store.Size()).To(Equal(uint64(1)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
				}))
			})
		})

		Context("from multiple boxed-quad structs", func() {
			store := NewQuadStore(
				[]Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", ""},
				})
			It("should have size 5", func() {
				Expect(store.Size()).To(Equal(uint64(5)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", ""},
				}))
			})
		})

		Context("from multiple pointers to boxed-quad structs", func() {
			store := NewQuadStore(
				[]*Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", ""},
				})
			It("should have size 5", func() {
				Expect(store.Size()).To(Equal(uint64(5)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", "o1", ""},
					{"s1", "p1", "o2", ""},
					{"s1", "p2", "o2", ""},
					{"s2", "p1", "o1", ""},
					{"s1", "p2", "o3", ""},
				}))
			})
		})

		Context("from an empty slice of structs", func() {
			store := NewQuadStore(
				[]boxedQuad2{},
			)
			It("should have size 0", func() {
				Expect(store.Size()).To(Equal(uint64(0)))
			})
		})

		Context("from an unknown type", func() {
			var s *QuadStore
			shouldPanic := func() {
				s = NewQuadStore(true)
			}
			It("should panic", func() {
				Expect(shouldPanic).To(Panic())
			})
		})

		Context("from a slice of unknown type", func() {
			var s *QuadStore
			shouldPanic := func() {
				s = NewQuadStore([]bool{true})
			}
			It("should panic", func() {
				Expect(shouldPanic).To(Panic())
			})
		})

		Context("from a slice of structs with unknown fields", func() {
			var s *QuadStore
			shouldPanic := func() {
				s = NewQuadStore([]struct{ X string }{{X: "foo"}})
			}
			It("should panic", func() {
				Expect(shouldPanic).To(Panic())
			})
		})

	})

	Describe("An empty QuadStore", func() {
		store := NewQuadStore()

		It("should have size 0", func() {
			Expect(store.Size()).To(Equal(uint64(0)))
		})
		It("should be empty", func() {
			Expect(store.Empty()).To(BeTrue())
		})

		Describe("Add", func() {

			Context("with a wildcard subject", func() {
				It("should panic", func() {
					Expect(func() { store.Add("*", "p1", "o1", "") }).To(Panic())
				})
			})

			Context("with a wildcard predicate", func() {
				It("should panic", func() {
					Expect(func() { store.Add("s1", "*", "o1", "") }).To(Panic())
				})
			})

			Context("with a wildcard object", func() {
				It("should panic", func() {
					Expect(func() { store.Add("s1", "p1", "*", "") }).To(Panic())
				})
			})

			Context("with a wildcard graph", func() {
				It("should panic", func() {
					Expect(func() { store.Add("s1", "p1", "o1", "*") }).To(Panic())
				})
			})
		})

	})

	Describe("non-string object values", func() {
		Context("a store initialised with non-string object values", func() {
			store := NewQuadStore(
				[]struct {
					Subject   string
					Predicate string
					Object    interface{}
					Graph     string
				}{{"s1", "p1", 3, ""}, {"s1", "p2", 3, ""}, {"s1", "p2", 4, ""}},
			)
			It("should have size 3", func() {
				Expect(store.Size()).To(Equal(uint64(3)))
			})
			It("should contain the correct quads", func() {
				resultsList := iterResults(store)
				Expect(resultsList).To(ConsistOf([]*Quad{
					{"s1", "p1", 3, ""},
					{"s1", "p2", 3, ""},
					{"s1", "p2", 4, ""},
				}))
			})
			It("should delete quads correctly", func() {
				Expect(store.Remove("s1", "p2", 4, "")).To(Equal(uint64(1)))
				Expect(store.Size()).To(Equal(uint64(2)))
			})
			It("should work correctly with find methods", func() {
				Expect(store.FindSubjects("p2", 3, "")).To(Equal([]string{"s1"}))
			})
			It("should delete further quads correctly", func() {
				Expect(store.Remove("s1", "p2", 3, "")).To(Equal(uint64(1)))
				Expect(store.Size()).To(Equal(uint64(1)))
			})
		})
	})

	Describe("A QuadStore initialised with 3 elements", func() {
		store := NewQuadStore([][3]string{
			{"s1", "p1", "o1"},
			{"s1", "p1", "o2"},
			{"s1", "p1", "o3"},
		})

		It("should have size 3", func() {
			Expect(store.Size()).To(Equal(uint64(3)))
		})
		It("should not be empty", func() {
			Expect(store.Empty()).To(BeFalse())
		})

		Describe("when adding a quad that did not exist yet", func() {
			It("should return true", func() {
				Expect(store.Add("s1", "p1", "o4", "")).To(BeTrue())
			})
			It("should increase the size", func() {
				Expect(store.Size()).To(Equal(uint64(4)))
			})
		})

		Describe("when adding a quad that already exists", func() {
			It("should return false", func() {
				Expect(store.Add("s1", "p1", "o4", "")).To(BeFalse())
			})
			It("should not change the size", func() {
				Expect(store.Size()).To(Equal(uint64(4)))
			})
		})

		Describe("when removing an existing quad", func() {
			It("should return 1", func() {
				Expect(store.Remove("s1", "p1", "o4", "")).To(Equal(uint64(1)))
			})
			It("should decrease the size", func() {
				Expect(store.Size()).To(Equal(uint64(3)))
			})
		})

		Describe("when removing a non-existing quad", func() {
			It("should return zero", func() {
				Expect(store.Remove("s1", "p1", "o5", "")).To(Equal(uint64(0)))
			})
			It("should not change the size", func() {
				Expect(store.Size()).To(Equal(uint64(3)))
			})
		})

		Describe("when removing all quads using wildcards", func() {
			It("should return 3", func() {
				Expect(store.Remove("*", "*", "*", "*")).To(Equal(uint64(3)))
			})
			It("should result in an empty store", func() {
				Expect(store.Size()).To(Equal(uint64(0)))
			})
		})
	})

	Describe("A QuadStore initialised with 5 elements", func() {
		quads := []*Quad{
			{"s1", "p1", "o1", ""},
			{"s1", "p1", "o2", ""},
			{"s1", "p2", "o2", ""},
			{"s2", "p1", "o1", ""},
			{"s1", "p2", "o3", "c4"},
		}

		store := NewQuadStore(quads)

		It("should have size 5", func() {
			Expect(store.Size()).To(Equal(uint64(5)))
		})
		It("should not be empty", func() {
			Expect(store.Empty()).To(BeFalse())
		})

		Describe("ForEach", func() {

			It("should iterate over all quads in the store", func() {
				var resultsList []*Quad
				store.ForEach(func(s, p string, o interface{}, g string) {
					resultsList = append(resultsList, &Quad{s, p, o, g})
				})
				Expect(resultsList).To(ConsistOf(quads))
			})
		})

		Describe("ForEachWith", func() {

			Context("with wildcard parameters", func() {
				It("should iterate over all quads in the store", func() {
					var resultsList []*Quad
					store.ForEachWith("*", "*", "*", "*", func(s, p string, o interface{}, g string) {
						resultsList = append(resultsList, &Quad{s, p, o, g})
					})
					Expect(resultsList).To(ConsistOf(quads))
				})
			})

			Context("with an existing subject parameter", func() {
				It("should iterate over all quads with a matching subject", func() {
					var resultsList []*Quad
					store.ForEachWith("s1", "*", "*", "*", func(s, p string, o interface{}, g string) {
						resultsList = append(resultsList, &Quad{s, p, o, g})
					})
					Expect(resultsList).To(ConsistOf([]*Quad{
						{"s1", "p1", "o1", ""},
						{"s1", "p1", "o2", ""},
						{"s1", "p2", "o2", ""},
						{"s1", "p2", "o3", "c4"},
					}))
				})
			})

			Context("with an non-existing subject parameter", func() {
				It("should do nothing", func() {
					var resultsList []*Quad
					store.ForEachWith("s0", "*", "*", "*", func(s, p string, o interface{}, g string) {
						resultsList = append(resultsList, &Quad{s, p, o, g})
					})
					Expect(resultsList).To(BeEmpty())
				})
			})

			Context("with an existing predicate parameter", func() {
				It("should iterate over all quads with a matching predicate", func() {
					var resultsList []*Quad
					store.ForEachWith("*", "p2", "*", "*", func(s, p string, o interface{}, g string) {
						resultsList = append(resultsList, &Quad{s, p, o, g})
					})
					Expect(resultsList).To(ConsistOf([]*Quad{
						{"s1", "p2", "o2", ""},
						{"s1", "p2", "o3", "c4"},
					}))
				})
			})

			Context("with an non-existing predicate parameter", func() {
				It("should do nothing", func() {
					var resultsList []*Quad
					store.ForEachWith("*", "p0", "*", "*", func(s, p string, o interface{}, g string) {
						resultsList = append(resultsList, &Quad{s, p, o, g})
					})
					Expect(resultsList).To(BeEmpty())
				})
			})

			Context("with an existing object parameter", func() {
				It("should iterate over all quads with a matching object", func() {
					var resultsList []*Quad
					store.ForEachWith("*", "*", "o2", "*", func(s, p string, o interface{}, g string) {
						resultsList = append(resultsList, &Quad{s, p, o, g})
					})
					Expect(resultsList).To(ConsistOf([]*Quad{
						{"s1", "p1", "o2", ""},
						{"s1", "p2", "o2", ""},
					}))
				})
			})

			Context("with an non-existing object parameter", func() {
				It("should do nothing", func() {
					var resultsList []*Quad
					store.ForEachWith("*", "*", "o0", "*", func(s, p string, o interface{}, g string) {
						resultsList = append(resultsList, &Quad{s, p, o, g})
					})
					Expect(resultsList).To(BeEmpty())
				})
			})

			Context("with an existing graph parameter", func() {
				It("should iterate over all quads with a matching object", func() {
					var resultsList []*Quad
					store.ForEachWith("*", "*", "*", "", func(s, p string, o interface{}, g string) {
						resultsList = append(resultsList, &Quad{s, p, o, g})
					})
					Expect(resultsList).To(ConsistOf([]*Quad{
						{"s1", "p1", "o1", ""},
						{"s1", "p1", "o2", ""},
						{"s1", "p2", "o2", ""},
						{"s2", "p1", "o1", ""},
					}))
				})
			})

			Context("with an non-existing graph parameter", func() {
				It("should do nothing", func() {
					var resultsList []*Quad
					store.ForEachWith("*", "*", "*", "c0", func(s, p string, o interface{}, g string) {
						resultsList = append(resultsList, &Quad{s, p, o, g})
					})
					Expect(resultsList).To(BeEmpty())
				})
			})
		})

		Describe("FindGraphs", func() {

			Context("with wildcard parameters", func() {
				It("should return all graph names", func() {
					Expect(store.FindGraphs("*", "*", "*")).To(ConsistOf("", "c4"))
				})
			})

			// s

			Context("with an existing subject parameter", func() {
				It("should return all graphs having this subject", func() {
					Expect(store.FindGraphs("s1", "*", "*")).To(ConsistOf("", "c4"))
				})
			})

			Context("with a non-existing subject parameter", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("s0", "*", "*")).To(BeEmpty())
				})
			})

			Context("with a non-existing subject parameter that exists elsewhere", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("p1", "*", "*")).To(BeEmpty())
				})
			})

			// p

			Context("with an existing predicate parameter", func() {
				It("should return all graphs having this predicate", func() {
					Expect(store.FindGraphs("*", "p2", "*")).To(ConsistOf("", "c4"))
				})
			})

			Context("with a non-existing predicate parameter", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("*", "p0", "*")).To(BeEmpty())
				})
			})

			Context("with a non-existing predicate parameter that exists elsewhere", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("*", "s1", "*")).To(BeEmpty())
				})
			})

			// o

			Context("with an existing object parameter", func() {
				It("should return all graphs having this object", func() {
					Expect(store.FindGraphs("*", "*", "o3")).To(ConsistOf("c4"))
				})
			})

			Context("with a non-existing object parameter", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("*", "*", "o0")).To(BeEmpty())
				})
			})

			Context("with a non-existing object parameter that exists elsewhere", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("*", "*", "s1")).To(BeEmpty())
				})
			})

			// sp

			Context("with existing subject and predicate parameters", func() {
				It("should return all graphs having this subject and predicate", func() {
					Expect(store.FindGraphs("s1", "p2", "*")).To(ConsistOf("", "c4"))
				})
			})

			Context("with existing non-matching subject and predicate parameters", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("s2", "p2", "*")).To(BeEmpty())
				})
			})

			Context("with non-existing subject and predicate parameters", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("s0", "p0", "*")).To(BeEmpty())
				})
			})

			// so

			Context("with existing subject and object parameters", func() {
				It("should return all graphs having this subject and object", func() {
					Expect(store.FindGraphs("s1", "*", "o3")).To(ConsistOf("c4"))
				})
			})

			Context("with existing non-matching subject and object parameters", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("s2", "*", "o2")).To(BeEmpty())
				})
			})

			Context("with non-existing subject and object parameters", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("s0", "*", "o0")).To(BeEmpty())
				})
			})

			// po

			Context("with existing predicate and object parameters", func() {
				It("should return all graphs having this predicate and object", func() {
					Expect(store.FindGraphs("*", "p1", "o1")).To(ConsistOf(""))
				})
			})

			Context("with existing non-matching predicate and object parameters", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("*", "p1", "o3")).To(BeEmpty())
				})
			})

			Context("with non-existing predicate and object parameters", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("*", "p0", "o0")).To(BeEmpty())
				})
			})

			// spo

			Context("with existing subject, predicate and object parameters", func() {
				It("should return all graphs having this subject, predicate and object", func() {
					Expect(store.FindGraphs("s1", "p1", "o1")).To(ConsistOf(""))
				})
			})

			Context("with existing non-matching subject, predicate and object parameters", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("s1", "p1", "o3")).To(BeEmpty())
				})
			})

			Context("with non-existing subject, predicate and object parameters", func() {
				It("should return no graphs", func() {
					Expect(store.FindGraphs("s0", "p0", "o0")).To(BeEmpty())
				})
			})
		})

		Describe("FindSubjects", func() {

			Context("with wildcard parameters", func() {
				It("should return all subjects in all graphs", func() {
					Expect(store.FindSubjects("*", "*", "*")).To(ConsistOf("s1", "s2"))
				})
			})

			// p

			Context("with an existing predicate parameter", func() {
				It("should return all subjects having this predicate", func() {
					Expect(store.FindSubjects("p1", "*", "*")).To(ConsistOf("s1", "s2"))
				})
			})

			Context("with a non-existing predicate parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("p0", "*", "*")).To(BeEmpty())
				})
			})

			Context("with a non-existing predicate parameter that exists elsewhere", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("s1", "*", "*")).To(BeEmpty())
				})
			})

			// o

			Context("with an existing object parameter", func() {
				It("should return all subjects having this object", func() {
					Expect(store.FindSubjects("*", "o1", "*")).To(ConsistOf("s1", "s2"))
				})
			})

			Context("with a non-existing object parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("*", "o0", "*")).To(BeEmpty())
				})
			})

			Context("with a non-existing object parameter that exists elsewhere", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("*", "s1", "*")).To(BeEmpty())
				})
			})

			// g

			Context("with an existing graph parameter", func() {
				It("should return all subjects having this graph", func() {
					Expect(store.FindSubjects("*", "*", "")).To(ConsistOf("s1", "s2"))
				})
			})

			Context("with a non-existing graph parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("*", "*", "o0")).To(BeEmpty())
				})
			})

			Context("with a non-existing graph parameter that exists elsewhere", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("*", "*", "s1")).To(BeEmpty())
				})
			})

			// po

			Context("with existing predicate and object parameters", func() {
				It("should return all subjects having this predicate and object", func() {
					Expect(store.FindSubjects("p1", "o1", "*")).To(ConsistOf("s1", "s2"))
				})
			})

			Context("with existing non-matching predicate and object parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("p1", "o3", "*")).To(BeEmpty())
				})
			})

			Context("with non-existing predicate and object parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("p0", "o0", "*")).To(BeEmpty())
				})
			})

			// og

			Context("with existing object and graph parameters", func() {
				It("should return all subjects having this object and graph", func() {
					Expect(store.FindSubjects("*", "o3", "c4")).To(ConsistOf("s1"))
				})
			})

			Context("with existing non-matching object and graph parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("*", "o3", "")).To(BeEmpty())
				})
			})

			Context("with non-existing object and graph parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("*", "o0", "c0")).To(BeEmpty())
				})
			})

			// pg

			Context("with existing predicate and graph parameters", func() {
				It("should return all subjects having this predicate and graph", func() {
					Expect(store.FindSubjects("p1", "*", "")).To(ConsistOf("s1", "s2"))
				})
			})

			Context("with existing non-matching predicate and graph parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("p1", "*", "c4")).To(BeEmpty())
				})
			})

			Context("with non-existing predicate and graph parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("p0", "*", "c0")).To(BeEmpty())
				})
			})

			// pog

			Context("with existing predicate, object and graph parameters", func() {
				It("should return all subjects having this predicate, object and graph", func() {
					Expect(store.FindSubjects("p1", "o1", "")).To(ConsistOf("s1", "s2"))
				})
			})

			Context("with existing non-matching predicate, object and graph parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("p2", "o1", "c4")).To(BeEmpty())
				})
			})

			Context("with non-existing predicate, object and graph parameter", func() {
				It("should return no subjects", func() {
					Expect(store.FindSubjects("p0", "o0", "c0")).To(BeEmpty())
				})
			})
		})

		Describe("FindPredicates", func() {

			Context("with wildcard parameters", func() {
				It("should return all predicates in all graphs", func() {
					Expect(store.FindPredicates("*", "*", "*")).To(ConsistOf("p1", "p2"))
				})
			})

			// s

			Context("with an existing subject parameter", func() {
				It("should return all predicates having this subject", func() {
					Expect(store.FindPredicates("s1", "*", "*")).To(ConsistOf("p1", "p2"))
				})
			})

			Context("with a non-existing subject parameter", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("s0", "*", "*")).To(BeEmpty())
				})
			})

			Context("with a non-existing subject parameter that exists elsewhere", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("o2", "*", "*")).To(BeEmpty())
				})
			})

			// o

			Context("with an existing object parameter", func() {
				It("should return all predicates having this object", func() {
					Expect(store.FindPredicates("*", "o2", "*")).To(ConsistOf("p1", "p2"))
				})
			})

			Context("with a non-existing object parameter", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("*", "o0", "*")).To(BeEmpty())
				})
			})

			Context("with a non-existing object parameter that exists elsewhere", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("*", "c4", "*")).To(BeEmpty())
				})
			})

			// g

			Context("with an existing graph parameter", func() {
				It("should return all predicates having this graph", func() {
					Expect(store.FindPredicates("*", "*", "")).To(ConsistOf("p1", "p2"))
				})
			})

			Context("with a non-existing graph parameter", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("*", "*", "c0")).To(BeEmpty())
				})
			})

			Context("with a non-existing graph parameter that exists elsewhere", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("*", "*", "s1")).To(BeEmpty())
				})
			})

			// so

			Context("with existing subject and object parameters", func() {
				It("should return all predicates having this subject and object", func() {
					Expect(store.FindPredicates("s1", "o2", "*")).To(ConsistOf("p1", "p2"))
				})
			})

			Context("with existing non-matching subject and object parameters", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("s2", "o3", "*")).To(BeEmpty())
				})
			})

			Context("with non-existing subject and object parameters", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("s0", "o0", "*")).To(BeEmpty())
				})
			})

			// sg

			Context("with existing subject and graph parameters", func() {
				It("should return all predicates having this subject and graph", func() {
					Expect(store.FindPredicates("s1", "*", "")).To(ConsistOf("p1", "p2"))
				})
			})

			Context("with existing non-matching subject and graph parameter", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("s2", "*", "c4")).To(BeEmpty())
				})
			})

			Context("with non-existing subject and graph parameter", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("s0", "*", "c0")).To(BeEmpty())
				})
			})

			// og

			Context("with existing object and graph parameters", func() {
				It("should return all predicates having this object and graph", func() {
					Expect(store.FindPredicates("*", "o2", "")).To(ConsistOf("p1", "p2"))
				})
			})

			Context("with existing non-matching object and graph parameter", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("*", "o3", "")).To(BeEmpty())
				})
			})

			Context("with non-existing object and graph parameter", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("*", "o0", "c0")).To(BeEmpty())
				})
			})

			// sog

			Context("with existing subject, object and graph parameters", func() {
				It("should return all predicates having this subject, object and graph", func() {
					Expect(store.FindPredicates("s1", "o2", "")).To(ConsistOf("p1", "p2"))
				})
			})

			Context("with existing non-matching subject, object and graph parameter", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("s2", "p1", "o2")).To(BeEmpty())
				})
			})

			Context("with non-existing subject, object and graph parameter", func() {
				It("should return no predicates", func() {
					Expect(store.FindPredicates("s0", "p0", "o0")).To(BeEmpty())
				})
			})

		})

		Describe("FindObjects", func() {

			Context("with wildcard parameters", func() {
				It("should return all objects in all graphs", func() {
					Expect(store.FindObjects("*", "*", "*")).To(ConsistOf("o1", "o2", "o3"))
				})
			})

			// s

			Context("with an existing subject parameter", func() {
				It("should return all objects having this subject", func() {
					Expect(store.FindObjects("s1", "*", "*")).To(ConsistOf("o1", "o2", "o3"))
				})
			})

			Context("with a non-existing subject parameter", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("s0", "*", "*")).To(BeEmpty())
				})
			})

			Context("with a non-existing subject parameter that exists elsewhere", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("p1", "*", "*")).To(BeEmpty())
				})
			})

			// p

			Context("with an existing predicate parameter", func() {
				It("should return all objects having this predicate", func() {
					Expect(store.FindObjects("*", "p2", "*")).To(ConsistOf("o2", "o3"))
				})
			})

			Context("with a non-existing predicate parameter", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("*", "p0", "*")).To(BeEmpty())
				})
			})

			Context("with a non-existing predicate parameter that exists elsewhere", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("*", "s2", "*")).To(BeEmpty())
				})
			})

			// g

			Context("with an existing graph parameter", func() {
				It("should return all objects having this graph", func() {
					Expect(store.FindObjects("*", "*", "")).To(ConsistOf("o1", "o2"))
				})
			})

			Context("with a non-existing graph parameter", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("*", "*", "o0")).To(BeEmpty())
				})
			})

			Context("with a non-existing graph parameter that exists elsewhere", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("*", "*", "s1")).To(BeEmpty())
				})
			})

			// sp

			Context("with existing subject and predicate parameters", func() {
				It("should return all objects having this subject and predicate", func() {
					Expect(store.FindObjects("s1", "p2", "*")).To(ConsistOf("o2", "o3"))
				})
			})

			Context("with existing non-matching subject and predicate parameters", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("s2", "p2", "*")).To(BeEmpty())
				})
			})

			Context("with non-existing subject and predicate parameters", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("s0", "p0", "*")).To(BeEmpty())
				})
			})

			// sg

			Context("with existing subject and graph parameters", func() {
				It("should return all objects having this subject and graph", func() {
					Expect(store.FindObjects("s1", "*", "")).To(ConsistOf("o1", "o2"))
				})
			})

			Context("with existing non-matching subject and graph parameter", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("s2", "*", "c4")).To(BeEmpty())
				})
			})

			Context("with non-existing subject and graph parameter", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("s0", "*", "c0")).To(BeEmpty())
				})
			})

			// pg

			Context("with existing predicate and graph parameters", func() {
				It("should return all objects having this predicate and graph", func() {
					Expect(store.FindObjects("*", "p1", "")).To(ConsistOf("o1", "o2"))
				})
			})

			Context("with existing non-matching predicate and graph parameter", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("*", "p1", "c4")).To(BeEmpty())
				})
			})

			Context("with non-existing predicate and graph parameter", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("*", "p0", "c0")).To(BeEmpty())
				})
			})

			// spg

			Context("with existing subject, predicate and graph parameters", func() {
				It("should return all objects having this subject, predicate and graph", func() {
					Expect(store.FindObjects("s1", "p1", "")).To(ConsistOf("o1", "o2"))
				})
			})

			Context("with existing non-matching subject, predicate and graph parameter", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("s1", "p1", "c4")).To(BeEmpty())
				})
			})

			Context("with non-existing subject, predicate and graph parameter", func() {
				It("should return no objects", func() {
					Expect(store.FindObjects("s0", "p0", "c0")).To(BeEmpty())
				})
			})

		})

		Describe("Count", func() {

			Context("with wildcard parameters", func() {
				It("should count all quads in the store", func() {
					Expect(store.Count("*", "*", "*", "*")).To(Equal(uint64(5)))
				})
			})

			Context("with existing default graph parameter", func() {
				It("should count all quads in the default graph", func() {
					Expect(store.Count("*", "*", "*", "")).To(Equal(uint64(4)))
				})
			})

			Context("with existing non-default graph parameter", func() {
				It("should count all quads in the that graph", func() {
					Expect(store.Count("*", "*", "*", "c4")).To(Equal(uint64(1)))
				})
			})

			Context("with non-existing non-default graph parameter", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "*", "*", "c0")).To(BeZero())
				})
			})

			// s

			Context("with an existing subject parameter", func() {
				It("should count all quads having this subject", func() {
					Expect(store.Count("s1", "*", "*", "*")).To(Equal(uint64(4)))
				})
			})

			Context("with a non-existing subject parameter", func() {
				It("should be zero", func() {
					Expect(store.Count("s0", "*", "*", "*")).To(BeZero())
				})
			})

			Context("with a non-existing subject parameter that exists elsewhere", func() {
				It("should be zero", func() {
					Expect(store.Count("p1", "*", "*", "*")).To(BeZero())
				})
			})

			// p

			Context("with an existing predicate parameter", func() {
				It("should count all quads having this predicate", func() {
					Expect(store.Count("*", "p2", "*", "*")).To(Equal(uint64(2)))
				})
			})

			Context("with a non-existing predicate parameter", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "p0", "*", "*")).To(BeZero())
				})
			})

			Context("with a non-existing predicate parameter that exists elsewhere", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "s1", "*", "*")).To(BeZero())
				})
			})

			// o

			Context("with an existing object parameter", func() {
				It("should count all quads having this object", func() {
					Expect(store.Count("*", "*", "o3", "*")).To(Equal(uint64(1)))
				})
			})

			Context("with a non-existing object parameter", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "*", "o0", "*")).To(BeZero())
				})
			})

			Context("with a non-existing object parameter that exists elsewhere", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "*", "c4", "*")).To(BeZero())
				})
			})

			// g

			Context("with an existing graph parameter", func() {
				It("should count all quads having this graph", func() {
					Expect(store.Count("*", "*", "*", "")).To(Equal(uint64(4)))
				})
			})

			Context("with a non-existing graph parameter", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "*", "*", "c0")).To(BeZero())
				})
			})

			Context("with a non-existing graph parameter that exists elsewhere", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "*", "*", "p2")).To(BeZero())
				})
			})

			// sp

			Context("with existing subject and predicate parameters", func() {
				It("should count all quads having this subject and predicate", func() {
					Expect(store.Count("s1", "p2", "*", "*")).To(Equal(uint64(2)))
				})
			})

			Context("with existing non-matching subject and predicate parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s2", "p2", "*", "*")).To(BeZero())
				})
			})

			Context("with non-existing subject and predicate parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s0", "p0", "*", "*")).To(BeZero())
				})
			})

			// so

			Context("with existing subject and object parameters", func() {
				It("should count all quads having this subject and object", func() {
					Expect(store.Count("s1", "*", "o2", "*")).To(Equal(uint64(2)))
				})
			})

			Context("with existing non-matching subject and object parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s2", "*", "o2", "*")).To(BeZero())
				})
			})

			Context("with non-existing subject and object parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s0", "*", "o0", "*")).To(BeZero())
				})
			})

			// sg

			Context("with existing subject and graph parameters", func() {
				It("should count all quads having this subject and graph", func() {
					Expect(store.Count("s1", "*", "*", "c4")).To(Equal(uint64(1)))
				})
			})

			Context("with existing non-matching subject and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s2", "*", "*", "c4")).To(BeZero())
				})
			})

			Context("with non-existing subject and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s0", "*", "*", "c0")).To(BeZero())
				})
			})

			// po

			Context("with existing predicate and object parameters", func() {
				It("should count all quads having this predicate and object", func() {
					Expect(store.Count("*", "p1", "o1", "*")).To(Equal(uint64(2)))
				})
			})

			Context("with existing non-matching predicate and object parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "p1", "o3", "*")).To(BeZero())
				})
			})

			Context("with non-existing predicate and object parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "p0", "o0", "*")).To(BeZero())
				})
			})

			// pg

			Context("with existing predicate and graph parameters", func() {
				It("should count all quads having this predicate and graph", func() {
					Expect(store.Count("*", "p1", "*", "")).To(Equal(uint64(3)))
				})
			})

			Context("with existing non-matching predicate and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "p1", "*", "c4")).To(BeZero())
				})
			})

			Context("with non-existing predicate and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "p0", "*", "")).To(BeZero())
				})
			})

			// og

			Context("with existing object and graph parameters", func() {
				It("should count all quads having this object and graph", func() {
					Expect(store.Count("*", "*", "o2", "")).To(Equal(uint64(2)))
				})
			})

			Context("with existing non-matching object and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "*", "o2", "c4")).To(BeZero())
				})
			})

			Context("with non-existing object and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "*", "o0", "c0")).To(BeZero())
				})
			})

			// spo

			Context("with existing subject, predicate and object parameters", func() {
				It("should count all quads having this subject, predicate and object", func() {
					Expect(store.Count("s1", "p1", "o1", "*")).To(Equal(uint64(1)))
				})
			})

			Context("with existing non-matching subject, predicate and object parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s2", "p1", "o3", "*")).To(BeZero())
				})
			})

			Context("with non-existing subject, predicate and object parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s0", "p0", "o0", "*")).To(BeZero())
				})
			})

			// spg

			Context("with existing subject, predicate and graph parameters", func() {
				It("should count all quads having this subject, predicate and graph", func() {
					Expect(store.Count("s1", "p1", "*", "")).To(Equal(uint64(2)))
				})
			})

			Context("with existing non-matching subject, predicate and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s1", "p1", "*", "c4")).To(BeZero())
				})
			})

			Context("with non-existing subject, predicate and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s0", "p0", "*", "c0")).To(BeZero())
				})
			})

			// sog

			Context("with existing subject, object and graph parameters", func() {
				It("should count all quads having this subject, object and graph", func() {
					Expect(store.Count("s1", "*", "o2", "")).To(Equal(uint64(2)))
				})
			})

			Context("with existing non-matching subject, object and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s2", "*", "o3", "c4")).To(BeZero())
				})
			})

			Context("with non-existing subject, object and graph parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s0", "*", "o0", "c0")).To(BeZero())
				})
			})

			// pog

			Context("with existing predicate, object and graph parameters", func() {
				It("should count all quads having this predicate, object and graph", func() {
					Expect(store.Count("*", "p1", "o2", "")).To(Equal(uint64(1)))
				})
			})

			Context("with existing non-matching predicate, object and graph parameter", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "p2", "o1", "c4")).To(BeZero())
				})
			})

			Context("with non-existing predicate, object and graph parameter", func() {
				It("should be zero", func() {
					Expect(store.Count("*", "p0", "o0", "c0")).To(BeZero())
				})
			})

			// spog

			Context("with existing subject, predicate, object and graph parameters", func() {
				It("should count all quads having this subject, predicate, object and graph", func() {
					Expect(store.Count("s1", "p1", "o2", "")).To(Equal(uint64(1)))
				})
			})

			Context("with existing subject, predicate, object and graph but non-matching subject parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s2", "p2", "o2", "")).To(BeZero())
				})
			})

			Context("with existing subject, predicate, object and graph but non-matching predicate parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s1", "p2", "o1", "")).To(BeZero())
				})
			})

			Context("with existing subject, predicate, object and graph but non-matching object parameters", func() {
				It("should be zero", func() {
					Expect(store.Count("s1", "p1", "o3", "")).To(BeZero())
				})
			})

			Context("with non-existing subject, predicate, object and graph parameter", func() {
				It("should be zero", func() {
					Expect(store.Count("s0", "p0", "o0", "")).To(BeZero())
				})
			})
		})

		Describe("Remove", func() {

			Context("when trying to remove a quad with a non-existing subject", func() {
				It("should return zero", func() {
					Expect(store.Remove("s0", "p1", "o1", "")).To(Equal(uint64(0)))
				})
				It("should still have size 5", func() {
					Expect(store.Size()).To(Equal(uint64(5)))
				})
			})

			Context("when trying to remove a quad with a non-existing predicate", func() {
				It("should return zero", func() {
					Expect(store.Remove("s1", "p0", "o1", "")).To(Equal(uint64(0)))
				})
				It("should still have size 5", func() {
					Expect(store.Size()).To(Equal(uint64(5)))
				})
			})

			Context("when trying to remove a quad with a non-existing object", func() {
				It("should return zero", func() {
					Expect(store.Remove("s1", "p1", "o0", "")).To(Equal(uint64(0)))
				})
				It("should still have size 5", func() {
					Expect(store.Size()).To(Equal(uint64(5)))
				})
			})

			Context("when trying to remove a quad with a non-existing graph", func() {
				It("should return zero", func() {
					Expect(store.Remove("s1", "p1", "o1", "c0")).To(Equal(uint64(0)))
				})
				It("should still have size 5", func() {
					Expect(store.Size()).To(Equal(uint64(5)))
				})
			})

			Context("when trying to remove quads with a non-existing subject and wildcards", func() {
				It("should return zero", func() {
					Expect(store.Remove("s0", "*", "*", "*")).To(Equal(uint64(0)))
				})
				It("should still have size 5", func() {
					Expect(store.Size()).To(Equal(uint64(5)))
				})
			})

			Context("when trying to remove quads with a non-existing predicate and wildcards", func() {
				It("should return zero", func() {
					Expect(store.Remove("*", "p0", "*", "*")).To(Equal(uint64(0)))
				})
				It("should still have size 5", func() {
					Expect(store.Size()).To(Equal(uint64(5)))
				})
			})

			Context("when trying to remove quads with a non-existing object and wildcards", func() {
				It("should return zero", func() {
					Expect(store.Remove("*", "*", "o0", "*")).To(Equal(uint64(0)))
				})
				It("should still have size 5", func() {
					Expect(store.Size()).To(Equal(uint64(5)))
				})
			})

			Context("when trying to remove quads with a non-existing graph and wildcards", func() {
				It("should return zero", func() {
					Expect(store.Remove("*", "*", "*", "c0")).To(Equal(uint64(0)))
				})
				It("should still have size 5", func() {
					Expect(store.Size()).To(Equal(uint64(5)))
				})
			})

			Context("when trying to remove a quad with existing non-matching parameters", func() {
				It("should return zero", func() {
					Expect(store.Remove("s2", "p2", "o1", "")).To(Equal(uint64(0)))
				})
				It("should still have size 5", func() {
					Expect(store.Size()).To(Equal(uint64(5)))
				})
			})

			Context("when removing an existing quad from the default graph", func() {
				It("should return 1", func() {
					Expect(store.Remove("s1", "p1", "o1", "")).To(Equal(uint64(1)))
				})
				It("should have size 4", func() {
					Expect(store.Size()).To(Equal(uint64(4)))
				})
			})

			Context("when removing an existing quad from the non-default graph", func() {
				It("should return 1", func() {
					Expect(store.Remove("s1", "p2", "o3", "c4")).To(Equal(uint64(1)))
				})
				It("should have size 3", func() {
					Expect(store.Size()).To(Equal(uint64(3)))
				})
			})

			Context("when adding and removing a quad", func() {
				It("should have unchanged size", func() {
					Expect(store.Add("s5", "p5", "o5", "c5")).To(BeTrue())
					Expect(store.Remove("s5", "p5", "o5", "c5")).To(Equal(uint64(1)))
					Expect(store.Size()).To(Equal(uint64(3)))
				})
			})

			// })

			Describe("with wildcards", func() {

				store := NewQuadStore()

				BeforeEach(func() {
					store.Add("s1", "p1", "o1", "")
					store.Add("s1", "p1", "o2", "")
					store.Add("s1", "p2", "o2", "")
					store.Add("s2", "p1", "o1", "")
					store.Add("s1", "p2", "o3", "c4")
				})

				Context("and non-existing graph parameter", func() {
					It("should have unchanged size", func() {
						Expect(store.Remove("*", "*", "*", "c0")).To(Equal(uint64(0)))
						Expect(store.Size()).To(Equal(uint64(5)))
					})
				})

				Context("and non-existing subject parameter", func() {
					It("should have unchanged size", func() {
						Expect(store.Remove("s0", "*", "*", "*")).To(Equal(uint64(0)))
						Expect(store.Size()).To(Equal(uint64(5)))
					})
				})

				Context("and non-existing predicate parameter", func() {
					It("should have unchanged size", func() {
						Expect(store.Remove("*", "p0", "*", "*")).To(Equal(uint64(0)))
						Expect(store.Size()).To(Equal(uint64(5)))
					})
				})

				Context("and non-existing object parameter", func() {
					It("should have unchanged size", func() {
						Expect(store.Remove("*", "*", "o0", "*")).To(Equal(uint64(0)))
						Expect(store.Size()).To(Equal(uint64(5)))
					})
				})

				//

				Context("and an existing graph parameter", func() {
					It("should have size 1", func() {
						Expect(store.Remove("*", "*", "*", "")).To(Equal(uint64(4)))
						Expect(store.Size()).To(Equal(uint64(1)))
					})
				})

				Context("and an existing subject parameter", func() {
					It("should have size 1", func() {
						Expect(store.Remove("s1", "*", "*", "*")).To(Equal(uint64(4)))
						Expect(store.Size()).To(Equal(uint64(1)))
					})
				})

				Context("and an existing predicate parameter", func() {
					It("should have size 3", func() {
						Expect(store.Remove("*", "p2", "*", "*")).To(Equal(uint64(2)))
						Expect(store.Size()).To(Equal(uint64(3)))
					})
				})

				Context("and an existing object parameter", func() {
					It("should have size 3", func() {
						Expect(store.Remove("*", "*", "o2", "*")).To(Equal(uint64(2)))
						Expect(store.Size()).To(Equal(uint64(3)))
					})
				})
			})
		})
	})

	Describe("Every", func() {

		It("should return false when the store is empty", func() {
			store := NewQuadStore()
			Expect(store.Size()).To(BeZero())
			Expect(store.Every(alwaysTrueFn)).To(BeFalse())
			Expect(store.Every(alwaysFalseFn)).To(BeFalse())
		})

		Context("with a store initialised with 5 items", func() {
			store := NewQuadStore([]*Quad{
				{"s1", "p1", "o1", ""},
				{"s1", "p1", "o2", ""},
				{"s1", "p1", "o3", ""},
				{"s1", "p1", "o4", ""},
				{"s1", "p1", "o5", ""},
			})

			It("should return true when the callback returns true for every item", func() {
				Expect(store.Every(alwaysTrueFn)).To(BeTrue())
			})
			It("should return false when the callback returns false for every item", func() {
				Expect(store.Every(alwaysFalseFn)).To(BeFalse())
			})

			callCount := 0
			trueTwiceFn := func(s, p string, o interface{}, g string) bool {
				callCount++
				if callCount > 1 {
					return false
				}
				return true
			}

			It("should return false immediately when the callback returns false for any item", func() {
				callCount = 0
				Expect(store.Every(trueTwiceFn)).To(BeFalse())
				Expect(callCount).To(Equal(2))
			})
		})
	})

	Describe("EveryWith", func() {

		Context("with an empty store", func() {
			store := NewQuadStore()

			It("should return false when called with all wildcards", func() {
				Expect(store.Size()).To(BeZero())
				Expect(store.EveryWith("*", "*", "*", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "*", "*", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing subject", func() {
				Expect(store.Size()).To(BeZero())
				Expect(store.EveryWith("s0", "*", "*", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("s0", "*", "*", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing predicate", func() {
				Expect(store.Size()).To(BeZero())
				Expect(store.EveryWith("*", "p0", "*", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "p0", "*", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing object", func() {
				Expect(store.Size()).To(BeZero())
				Expect(store.EveryWith("*", "*", "o0", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "*", "o0", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing graph", func() {
				Expect(store.Size()).To(BeZero())
				Expect(store.EveryWith("*", "*", "*", "c0", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "*", "*", "c0", alwaysFalseFn)).To(BeFalse())
			})
		})

		Context("with a store initialised with 5 items", func() {
			store := NewQuadStore([]*Quad{
				{"s1", "p1", "o1", ""},
				{"s1", "p1", "o2", ""},
				{"s1", "p1", "o3", ""},
				{"s1", "p1", "o4", ""},
				{"s1", "p1", "o5", ""},
			})

			It("should return true when the callback returns true for every item", func() {
				Expect(store.EveryWith("*", "*", "*", "*", alwaysTrueFn)).To(BeTrue())
			})
			It("should return false when the callback returns false for every item", func() {
				Expect(store.EveryWith("*", "*", "*", "*", alwaysFalseFn)).To(BeFalse())
			})

			callCount := 0
			trueTwiceFn := func(s, p string, o interface{}, g string) bool {
				callCount++
				if callCount > 1 {
					return false
				}
				return true
			}

			It("should return false immediately when the callback returns false for any item", func() {
				callCount = 0
				Expect(store.EveryWith("*", "*", "*", "*", trueTwiceFn)).To(BeFalse())
				Expect(callCount).To(Equal(2))
			})
			It("should return false when called with a non-existing subject", func() {
				Expect(store.EveryWith("s0", "*", "*", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("s0", "*", "*", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing subject that exists elsewhere", func() {
				Expect(store.EveryWith("p1", "*", "*", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("p1", "*", "*", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing predicate", func() {
				Expect(store.EveryWith("*", "p0", "*", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "p0", "*", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing predicate that exists elsewhere", func() {
				Expect(store.EveryWith("*", "s1", "*", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "s1", "*", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing object", func() {
				Expect(store.EveryWith("*", "*", "o0", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "*", "o0", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing object that exists elsewhere", func() {
				Expect(store.EveryWith("*", "*", "p1", "*", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "*", "p1", "*", alwaysFalseFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing graph", func() {
				Expect(store.EveryWith("*", "*", "*", "c0", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "*", "*", "c0", alwaysTrueFn)).To(BeFalse())
			})
			It("should return false when called with a non-existing graph that exists elsewhere", func() {
				Expect(store.EveryWith("*", "*", "*", "s1", alwaysTrueFn)).To(BeFalse())
				Expect(store.EveryWith("*", "*", "*", "s1", alwaysTrueFn)).To(BeFalse())
			})
		})
	})

	Describe("Some", func() {

		It("should return false when the store is empty", func() {
			store := NewQuadStore()
			Expect(store.Size()).To(BeZero())
			Expect(store.Some(alwaysTrueFn)).To(BeFalse())
			Expect(store.Some(alwaysFalseFn)).To(BeFalse())
		})

		Context("with a store initialised with 5 items", func() {
			store := NewQuadStore([]*Quad{
				{"s1", "p1", "o1", ""},
				{"s1", "p1", "o2", ""},
				{"s1", "p1", "o3", ""},
				{"s1", "p1", "o4", ""},
				{"s1", "p1", "o5", ""},
			})

			It("should return false when the callback returns false for every item", func() {
				Expect(store.Some(alwaysFalseFn)).To(BeFalse())
			})

			callCount := 0
			falseTwiceFn := func(s, p string, o interface{}, g string) bool {
				callCount++
				if callCount > 1 {
					return true
				}
				return false
			}

			It("should return true immediately when the callback returns true for any item", func() {
				callCount = 0
				Expect(store.Some(falseTwiceFn)).To(BeTrue())
				Expect(callCount).To(Equal(2))
			})
		})
	})

	Describe("Query", func() {

		store := NewQuadStore([]*Quad{
			{"s1", "p1", "o1", ""},
			{"s1", "p1", "o2", ""},
			{"s1", "p2", "o2", ""},
			{"s2", "p1", "o1", ""},
			{"s1", "p2", "o3", "c4"},
		})

		Context("when called with an unknown type", func() {
			It("should panic", func() {
				shouldPanicFn := func() {
					store.Query(false, "*")
				}
				Expect(shouldPanicFn).To(Panic())
			})
		})

		Context("when called with an empty pattern", func() {
			pattern := [][2]string{}
			results := store.Query(pattern, "*")
			It("should return nothing", func() {
				Expect(len(results)).To(Equal(0))
			})
		})

		Context("when called with non-existing pattern", func() {
			pattern := [][2]string{{"p0", "o0"}}
			results := store.Query(pattern, "*")
			It("should return nothing", func() {
				Expect(len(results)).To(Equal(0))
			})
		})

		Context("when called with a matching pattern as [][2]string", func() {
			pattern := [][2]string{{"p1", "o1"}, {"p1", "o2"}}
			results := store.Query(pattern, "*")
			It("should return the correct results", func() {
				Expect(len(results)).To(Equal(1))
				Expect(results[0].Subject).To(Equal("s1"))
			})
		})

		Context("when called with a matching pattern as [2]string", func() {
			pattern := [2]string{"p2", "o2"}
			results := store.Query(pattern, "*")
			It("should return the correct results", func() {
				Expect(len(results)).To(Equal(1))
				Expect(results[0].Subject).To(Equal("s1"))
			})
		})

		Context("when called with a matching pattern as map[string]string", func() {
			pattern := map[string]string{"p1": "o1", "p2": "o3"}
			results := store.Query(pattern, "*")
			It("should return the correct results", func() {
				Expect(len(results)).To(Equal(1))
				Expect(results[0].Subject).To(Equal("s1"))
			})
		})

		Context("when called with a matching pattern as map[string][]string", func() {
			pattern := map[string][]string{"p1": {"o1", "o2"}, "p2": {"o3"}}
			results := store.Query(pattern, "*")
			It("should return the correct results", func() {
				Expect(len(results)).To(Equal(1))
				Expect(results[0].Subject).To(Equal("s1"))
			})
		})
	})

	Describe("OnAdd", func() {

		store := NewQuadStore()

		It("should be called when adding a quad", func() {
			var subject, predicate, graph string
			var object interface{}
			store.OnAdd = func(s, p string, o interface{}, g string) {
				subject = s
				predicate = p
				object = o
				graph = g
			}
			store.Add("s1", "p1", "o1", "g1")
			Expect(subject).To(Equal("s1"))
			Expect(predicate).To(Equal("p1"))
			Expect(object).To(Equal("o1"))
			Expect(graph).To(Equal("g1"))
		})
		It("should not be called when adding an already existing quad", func() {
			count := 0
			store.OnAdd = func(s, p string, o interface{}, g string) {
				count++
			}
			store.Add("s1", "p1", "o1", "g1")
			Expect(count).To(Equal(0))
		})
	})

	Describe("OnRemove", func() {

		store := NewQuadStore([]*Quad{
			{"s1", "p1", "o1", ""},
			{"s1", "p1", "o2", ""},
			{"s1", "p2", "o2", ""},
			{"s2", "p1", "o1", ""},
			{"s1", "p2", "o3", "c4"},
		})

		It("should be called when removing a quad", func() {
			var subject, predicate, graph string
			var object interface{}
			store.OnRemove = func(s, p string, o interface{}, g string) {
				subject = s
				predicate = p
				object = o
				graph = g
			}
			store.Remove("s1", "p1", "o1", "")
			Expect(subject).To(Equal("s1"))
			Expect(predicate).To(Equal("p1"))
			Expect(object).To(Equal("o1"))
			Expect(graph).To(Equal(""))
		})
		It("should be called when removing multiple quads", func() {
			count := 0
			store.OnRemove = func(s, p string, o interface{}, g string) {
				count++
			}
			store.Remove("s1", "*", "*", "")
			Expect(count).To(Equal(2))
		})
	})
})
