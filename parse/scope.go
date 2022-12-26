package parse

type Scope struct {
	variables map[string]Node
	outer     *Scope
}

func NewScope(outer *Scope) *Scope {
	return &Scope{
		variables: make(map[string]Node),
		outer:     outer,
	}
}

func (s *Scope) Resolve(name string) (Node, bool) {
	node, ok := s.variables[name]
	if !ok && s.outer != nil {
		return s.outer.Resolve(name)
	}
	return node, ok
}

func (s *Scope) Define(name string, node Node) {
	s.variables[name] = node
}

func (s *Scope) Parent() *Scope {
	return s.outer
}
