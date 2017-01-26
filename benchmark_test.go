package store4_test

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"testing"

	"github.com/jimsmart/store4"
)

// $ go test -run=XXX -bench=.

// What are we timing here? I don't think much of this has any meaning.

var storeFind *store4.QuadStore
var lastDim int

// func init() {
// 	storeFind = newStoreForFind()
// }

const (
	namespace = "urn:store4:benchmark/padded-iri#"
)

func newStoreWithRandomQuads(count uint64) *store4.QuadStore {
	s := store4.NewQuadStore()
	for s.Size() < count {
		q := randomQuad("")
		s.Add(q[0], q[1], q[2], q[3])
	}
	return s
}

func randomQuad(graph string) [4]string {
	var q [4]string
	// TODO(js) Random quads do not represent the distribution of real-world data!!! :/
	q[0] = randomIRI()
	q[1] = randomIRI()
	q[2] = randomIRI()
	q[3] = graph
	return q
}

func randomIRI() string {
	return namespace + fmt.Sprintf("%16x", rand.Int63())
}

func randomQuadn(n int, graph string) [4]string {
	var q [4]string
	// TODO(js) Random quads do not represent the distribution of real-world data!!! :/
	q[0] = randomIRIn(n)
	q[1] = randomIRIn(n)
	q[2] = randomIRIn(n)
	q[3] = graph
	return q
}

func randomIRIn(n int) string {
	return namespace + fmt.Sprintf("%16x", rand.Int63n(int64(n)))
}

func benchmarkAddRemoveRandomsTo(count uint64, b *testing.B) {
	s := newStoreWithRandomQuads(count)
	q := randomQuad("")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s.Add(q[0], q[1], q[2], q[3])
		s.Remove(q[0], q[1], q[2], q[3])
	}
}

func BenchmarkAddRemoveRandomsToEmpty(b *testing.B)       { benchmarkAddRemoveRandomsTo(0, b) }
func BenchmarkAddRemoveRandomsTo10kRandoms(b *testing.B)  { benchmarkAddRemoveRandomsTo(10000, b) }
func BenchmarkAddRemoveRandomsTo100kRandoms(b *testing.B) { benchmarkAddRemoveRandomsTo(100000, b) }
func BenchmarkAddRemoveRandomsTo500kRandoms(b *testing.B) { benchmarkAddRemoveRandomsTo(500000, b) }
func BenchmarkAddRemoveRandomsTo1mRandoms(b *testing.B)   { benchmarkAddRemoveRandomsTo(1000000, b) }

// // func BenchmarkAddRemoveRandomsTo10mRandoms(b *testing.B)  { benchmarkAddRemoveRandomsTo(10000000, b) }
// // func BenchmarkAddRemoveRandomsTo50mRandoms(b *testing.B)  { benchmarkAddRemoveRandomsTo(50000000, b) }
// // func BenchmarkAddRemoveRandomsTo100mRandoms(b *testing.B) { benchmarkAddRemoveRandomsTo(100000000, b) }

func logMemStatsDifference(before, after *runtime.MemStats) {
	out := "Diff "
	out = out + fmtMemStatMiB("Alloc", after.Alloc-before.Alloc)
	out = out + fmtMemStatMiB("TotalAlloc", after.TotalAlloc-before.TotalAlloc)
	out = out + fmtMemStatMiB("Sys", after.Sys-before.Sys)
	log.Println(out)
}

func fmtMemStatMiB(name string, n uint64) string {
	return fmt.Sprintf("- %s: %.2fMiB ", name, float64(n)/1024/1024)
}

func readMemStats() *runtime.MemStats {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	return m
}

func createStoreForFind(dim int) *store4.QuadStore {
	if storeFind != nil && lastDim == dim {
		return storeFind
	}

	before := readMemStats()

	storeFind = newStoreWithDim(dim)
	log.Println("QuadStore.Size", storeFind.Size())
	lastDim = dim

	after := readMemStats()
	logMemStatsDifference(before, after)

	return storeFind
}

func newStoreWithDim(dim int) *store4.QuadStore {
	// dimSquared := dim * dim
	// dimCubed := dimSquared * dim
	// dimQuads := dimCubed * dim
	s := store4.NewQuadStore()
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			for k := 0; k < dim; k++ {
				s.Add(iri(i), iri(j), iri(k), "")
			}
		}
	}
	return s
}

func iri(i int) string {
	return namespace + fmt.Sprintf("%16x", i)
}

func BenchmarkSomeWith_NoWildcards(b *testing.B) {

	fn := func(s, p string, o interface{}, g string) bool {
		return true
	}

	dim := 256
	s := createStoreForFind(dim)
	q := randomQuadn(dim, "")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if !s.SomeWith(q[0], q[1], q[2], q[3], fn) {
			panic("bad SomeWith")
		}
	}
}

func BenchmarkSomeWith_WildcardO(b *testing.B) {

	fn := func(s, p string, o interface{}, g string) bool {
		return true
	}

	dim := 256
	s := createStoreForFind(dim)
	q := randomQuadn(dim, "")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if !s.SomeWith(q[0], q[1], "*", q[3], fn) {
			panic("bad SomeWith")
		}
	}
}

func BenchmarkSomeWith_WildcardPO(b *testing.B) {

	fn := func(s, p string, o interface{}, g string) bool {
		return true
	}

	dim := 256
	s := createStoreForFind(dim)
	q := randomQuadn(dim, "")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if !s.SomeWith(q[0], "*", "*", q[3], fn) {
			panic("bad SomeWith")
		}
	}
}

func BenchmarkSomeWith_WildcardSPO(b *testing.B) {

	fn := func(s, p string, o interface{}, g string) bool {
		return true
	}

	dim := 256
	s := createStoreForFind(dim)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if !s.SomeWith("*", "*", "*", "", fn) {
			panic("bad SomeWith")
		}
	}
}
