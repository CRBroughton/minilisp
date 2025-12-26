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

func macroexpand(e *Expr, env *Env) *Expr {
	if e == nilExpr || e.Type != Cons {
		return e
	}

	op := e.Car
	if op.Type != Symbol {
		return e
	}

	// Look up the operator
	val, ok := env.Lookup(op.Sym)
	if !ok || val.Type != Macro {
		return e
	}

	// Apply macro to unevaluated arguments
	args := listToSlice(e.Cdr)
	newEnv := NewEnv(val.Env)

	// Bind parameters to unevaluated arguments
	params := val.Params
	for _, arg := range args {
		if params == nilExpr {
			panic("macro: too many arguments")
		}
		newEnv.Define(params.Car.Sym, arg)
		params = params.Cdr
	}

	// Evaluate macro body to get new code
	expanded := eval(val.Body, newEnv)

	// Recursively expand the result
	return macroexpand(expanded, env)
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
		e = macroexpand(e, env)
		if e.Type != Cons {
			return eval(e, env)
		}

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
			case "macro":
				params := args.Car
				body := args.Cdr.Car
				return makeLambda(params, body, env, Macro)
			case "lambda":
				params := args.Car
				body := args.Cdr.Car
				return makeLambda(params, body, env, Lambda)
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

		if fn.Type == Lambda {
			newEnv := NewEnv(fn.Env)

			params := fn.Params
			for _, arg := range evaledArgs {
				if params == nilExpr {
					panic("too many arguments")
				}
				newEnv.Define(params.Car.Sym, arg)
				params = params.Cdr
			}
			return eval(fn.Body, newEnv)
		}

		panic(fmt.Sprintf("not a function: %s", printExpr(fn)))
	default:
		return e
	}
}
