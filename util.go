package store4

import (
	"fmt"
	"sort"
)

// objectSlice implements sort.Interface for []interface{}
// ordering by fields by fmt.Sprintf value.
type objectSlice []interface{}

func (o objectSlice) Len() int { return len(o) }

func (o objectSlice) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (o objectSlice) Less(i, j int) bool {
	// Object.
	return fmt.Sprint(o[i]) < fmt.Sprint(o[j])
}

// sortObjects sorts a slice of object values,
// by fmt.Sprintf value.
func sortObjects(slice []interface{}) {
	sort.Sort(objectSlice(slice))
}

// // quadSlice implements sort.Interface for []*Quad
// type quadSlice []*Quad

// func (q quadSlice) Len() int { return len(q) }

// func (q quadSlice) Swap(i, j int) { q[i], q[j] = q[j], q[i] }

// func (q quadSlice) Less(i, j int) bool {
// 	qi, qj := q[i], q[j]
// 	// Graph.
// 	gi, gj := qi.G, qj.G
// 	if gi < gj {
// 		return true
// 	}
// 	if gi > gj {
// 		return false
// 	}
// 	// Subject.
// 	si, sj := qi.S, qj.S
// 	if si < sj {
// 		return true
// 	}
// 	if si > sj {
// 		return false
// 	}
// 	// Predicate.
// 	pi, pj := qi.P, qj.P
// 	if pi < pj {
// 		return true
// 	}
// 	if pi > pj {
// 		return false
// 	}
// 	// Object.
// 	oi, oj := qi.O, qj.O
// 	return fmt.Sprint(oi) < fmt.Sprint(oj)
// }

// // SortQuads sorts a slice of quads,
// // by graph then subject then predicate then object.
// func SortQuads(slice []*Quad) {
// 	sort.Sort(quadSlice(slice))
// }

// // tripleSlice implements sort.Interface for []*Triple
// // ordering by fields SPO.
// type tripleSlice []*Triple

// func (t tripleSlice) Len() int { return len(t) }

// func (t tripleSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// func (t tripleSlice) Less(i, j int) bool {
// 	ti, tj := t[i], t[j]
// 	// Subject.
// 	si, sj := ti.S, tj.S
// 	if si < sj {
// 		return true
// 	}
// 	if si > sj {
// 		return false
// 	}
// 	// Predicate.
// 	pi, pj := ti.P, tj.P
// 	if pi < pj {
// 		return true
// 	}
// 	if pi > pj {
// 		return false
// 	}
// 	// Object.
// 	oi, oj := ti.O, tj.O
// 	return fmt.Sprint(oi) < fmt.Sprint(oj)
// }

// // SortTriples sorts a slice of subject-predicate-object triples,
// // by subject then predicate then object.
// func SortTriples(slice []*Triple) {
// 	sort.Sort(tripleSlice(slice))
// }

// // tupleSlice implements sort.Interface for []*Tuple
// // ordering by fields PO.
// type tupleSlice []*Tuple

// func (t tupleSlice) Len() int { return len(t) }

// func (t tupleSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// func (t tupleSlice) Less(i, j int) bool {
// 	ti, tj := t[i], t[j]
// 	// Predicate.
// 	si, sj := ti.P, tj.P
// 	if si < sj {
// 		return true
// 	}
// 	if si > sj {
// 		return false
// 	}
// 	// Object.
// 	oi, oj := ti.O, tj.O
// 	return fmt.Sprint(oi) < fmt.Sprint(oj)
// }

// // SortTuples sorts a slice of predicate-object tuples,
// // by predicate then object.
// func SortTuples(slice []*Tuple) {
// 	sort.Sort(tupleSlice(slice))
// }
