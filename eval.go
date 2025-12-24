package main

import "fmt"

func eval(e *Expr, env *Env) *Expr {
	switch e.Type {
	// these types are self-evaluating
	case Nil, Bool, Number:
		return e
	case Symbol:
		val, ok := env.Lookup(e.Sym)
		if !ok {
			panic(fmt.Sprintf("unbound symbol: %s", e.Sym))
		}
		return val
	case Cons:
		// TODO - doing this tomorrow, hopefully
		return nilExpr
	default:
		return e
	}
}
