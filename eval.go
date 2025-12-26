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

		if op.Type == Symbol {
			switch op.Sym {
			case "quote":
				// (quote x) → x (unevaluated)
				return args.Car
			case "if":
				cond := eval(args.Car, env)
				if cond != nilExpr {
					return eval(args.Cdr.Car, env)
				}
				return eval(args.Cdr.Cdr.Car, env)
			case "define":
				sym := args.Car
				val := eval(args.Cdr.Car, env)
				env.Define(sym.Sym, val)
				return val
			case "begin":
				var result *Expr = nilExpr
				for args != nilExpr {
					result = eval(args.Car, env)
					args = args.Cdr
				}
				return result
			}
		}

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
