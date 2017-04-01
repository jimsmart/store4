package store4

type Builder struct {
	QuadStore *QuadStore
	graph     string
	subject   string
}

func (b *Builder) Graph(g string) *Builder {
	b.graph = g
	b.subject = ""
	return b
}

func (b *Builder) DefaultGraph() *Builder {
	b.Graph("")
	return b
}

func (b *Builder) Subject(s string) *Builder {
	b.subject = s
	return b
}

func (b *Builder) Add(p string, o interface{}) *Builder {
	if len(b.subject) == 0 {
		panic("subject term undefined")
	}
	if len(p) == 0 {
		panic("predicate term undefined")
	}
	if o == nil {
		panic("object term undefined")
	}
	if b.QuadStore == nil {
		b.QuadStore = NewQuadStore()
	}
	b.QuadStore.Add(b.subject, p, o, b.graph)
	return b
}

func (b *Builder) Build() *QuadStore {
	if b.QuadStore == nil {
		b.QuadStore = NewQuadStore()
	}
	return b.QuadStore
}
