package store4

// Builder provides a fluent (chainable) convenience API for creating quads.
//
// If no quad store is provided, a new one will be created.
type Builder struct {
	QuadStore *QuadStore
	graph     string
	subject   string
}

// Graph sets the name of the current graph in which to add new statements.
func (b *Builder) Graph(g string) *Builder {
	b.graph = g
	b.subject = ""
	return b
}

// DefaultGraph sets the current graph in which to add new statements to be
// the default graph.
func (b *Builder) DefaultGraph() *Builder {
	return b.Graph("")
}

// Subject sets the current subject about which statements are to be added.
func (b *Builder) Subject(s string) *Builder {
	b.subject = s
	return b
}

// Add a statement to the underlying quad store, using the current subject and graph.
func (b *Builder) Add(p string, o interface{}) *Builder {
	if len(b.subject) == 0 {
		panic("Undefined subject term")
	}
	if len(p) == 0 {
		panic("Undefined predicate term")
	}
	if o == nil {
		panic("Undefined object term")
	}
	if b.QuadStore == nil {
		b.QuadStore = NewQuadStore()
	}
	b.QuadStore.Add(b.subject, p, o, b.graph)
	return b
}

// Build returns the underlying quad store.
func (b *Builder) Build() *QuadStore {
	if b.QuadStore == nil {
		b.QuadStore = NewQuadStore()
	}
	return b.QuadStore
}
