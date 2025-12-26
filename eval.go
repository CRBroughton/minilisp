package main

import "fmt"

// Evaluate a list of expressions
func evalList(list *Expr, env *Env) []*Expr {
	var result []*Expr
	for list != nilExpr && list.Type == Pair {
		result = append(result, eval(list.Head, env))
		list = list.Tail
	}
	return result
}

func macroexpand(e *Expr, env *Env) *Expr {
	if e == nilExpr || e.Type != Pair {
		return e
	}

	op := e.Head
	if op.Type != Symbol {
		return e
	}

	// Look up the operator
	val, ok := env.Lookup(op.Sym)
	if !ok || val.Type != Macro {
		return e
	}

	// Apply macro to unevaluated arguments
	args := listToSlice(e.Tail)
	newEnv := NewEnv(val.Env)

	// Bind parameters to unevaluated arguments
	params := val.Params
	for _, arg := range args {
		if params == nilExpr {
			panic("macro: too many arguments")
		}
		newEnv.Define(params.Head.Sym, arg)
		params = params.Tail
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
	case Pair:
		e = macroexpand(e, env)
		if e.Type != Pair {
			return eval(e, env)
		}

		op := e.Head
		args := e.Tail

		if op.Type == Symbol {
			switch op.Sym {
			case "quote":
				// (quote x) → x (unevaluated)
				return args.Head
			case "if":
				cond := eval(args.Head, env)
				if cond != nilExpr {
					return eval(args.Tail.Head, env)
				}
				return eval(args.Tail.Tail.Head, env)
			case "define":
				sym := args.Head
				val := eval(args.Tail.Head, env)
				env.Define(sym.Sym, val)
				return val
			case "macro":
				params := args.Head
				body := args.Tail.Head
				return makeLambda(params, body, env, Macro)
			case "lambda":
				params := args.Head
				body := args.Tail.Head
				return makeLambda(params, body, env, Lambda)
			case "begin":
				var result *Expr = nilExpr
				for args != nilExpr {
					result = eval(args.Head, env)
					args = args.Tail
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
				newEnv.Define(params.Head.Sym, arg)
				params = params.Tail
			}
			return eval(fn.Body, newEnv)
		}

		panic(fmt.Sprintf("not a function: %s", printExpr(fn)))
	default:
		return e
	}
}
