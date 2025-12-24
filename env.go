package main

type Env struct {
	bindings map[string]*Expr
	parent   *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		bindings: make(map[string]*Expr),
		parent:   parent,
	}
}

func (e *Env) Define(sym string, val *Expr) {
	e.bindings[sym] = val
}

func (e *Env) Lookup(sym string) (*Expr, bool) {
	if val, ok := e.bindings[sym]; ok {
		return val, true
	}

	// recursively check the parent
	if e.parent != nil {
		return e.parent.Lookup(sym)
	}

	// didnt find
	return nil, false
}
