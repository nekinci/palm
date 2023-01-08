package parse

type Scope struct {
	variables map[string]any
	outer     *Scope
}

func NewScope(outer *Scope) *Scope {
	return &Scope{
		variables: make(map[string]any),
		outer:     outer,
	}
}

func (s *Scope) Resolve(name string) (any, bool) {
	val, ok := s.variables[name]
	if !ok && s.outer != nil {
		return s.outer.Resolve(name)
	}
	return val, ok
}

func (s *Scope) ResolveLocal(name string) (any, bool) {
	val, ok := s.variables[name]
	return val, ok
}

func (s *Scope) Define(name string, val any) {
	s.variables[name] = val
}

func (s *Scope) Parent() *Scope {
	return s.outer
}
