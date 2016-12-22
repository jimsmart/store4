package store4

import "sort"

// quadSlice implements sort.Interface for [][4]string
// ordering by fields GSPO [3,0,1,2].
type quadSlice [][4]string

func (q quadSlice) Len() int { return len(q) }

func (q quadSlice) Swap(i, j int) { q[i], q[j] = q[j], q[i] }

func (q quadSlice) Less(i, j int) bool {
	qi, qj := q[i], q[j]
	// Graph.
	gi, gj := qi[3], qj[3]
	if gi < gj {
		return true
	}
	if gi > gj {
		return false
	}
	// Subject.
	si, sj := qi[0], qj[0]
	if si < sj {
		return true
	}
	if si > sj {
		return false
	}
	// Predicate.
	pi, pj := qi[1], qj[1]
	if pi < pj {
		return true
	}
	if pi > pj {
		return false
	}
	// Object.
	oi, oj := qi[2], qj[2]
	return oi < oj
}

// SortQuads sorts a slice of subject-predicate-object-graph
// quads represented as [4]string,
// by graph then subject then predicate then object.
func SortQuads(slice [][4]string) {
	sort.Sort(quadSlice(slice))
}

// tripleSlice implements sort.Interface for [][3]string
// ordering by fields SPO [0,1,2].
type tripleSlice [][3]string

func (t tripleSlice) Len() int { return len(t) }

func (t tripleSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t tripleSlice) Less(i, j int) bool {
	ti, tj := t[i], t[j]
	// Subject.
	si, sj := ti[0], tj[0]
	if si < sj {
		return true
	}
	if si > sj {
		return false
	}
	// Predicate.
	pi, pj := ti[1], tj[1]
	if pi < pj {
		return true
	}
	if pi > pj {
		return false
	}
	// Object.
	oi, oj := ti[2], tj[2]
	return oi < oj
}

// SortTriples sorts a slice of subject-predicate-object
// triples represented as [3]string,
// by subject then predicate then object.
func SortTriples(slice [][3]string) {
	sort.Sort(tripleSlice(slice))
}

// tupleSlice implements sort.Interface for [][2]string
// ordering by fields PO [0,1].
type tupleSlice [][2]string

func (t tupleSlice) Len() int { return len(t) }

func (t tupleSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t tupleSlice) Less(i, j int) bool {
	ti, tj := t[i], t[j]
	// Predicate.
	si, sj := ti[0], tj[0]
	if si < sj {
		return true
	}
	if si > sj {
		return false
	}
	// Object.
	oi, oj := ti[1], tj[1]
	return oi < oj
}

// SortTuples sorts a slice of predicate-object
// tuples represented as [2]string,
// by predicate then object.
func SortTuples(slice [][2]string) {
	sort.Sort(tupleSlice(slice))
}
