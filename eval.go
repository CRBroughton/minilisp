package main

import "fmt"

// Evaluate a list of expressions
func evalList(list *Expr, env *Env) []*Expr {
	var result []*Expr
	for list != nilExpr && list.Type == Cons {
		result = append(result, eval(list.Car, env))
		list = list.Cdr
	}
	return result
}

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
	// Head | tail support
	case Cons:
		op := e.Car
		args := e.Cdr

		fn := eval(op, env)
		evaledArgs := evalList(args, env)

		if fn.Type == Builtin {
			return fn.Fn(evaledArgs)
		}
		// TODO - lambda support
		return nilExpr
	default:
		return e
	}
}
