package main

import (
	"fmt"
	"os"
)

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
	i := 0
	for params != nilExpr {
		// Check for &rest parameter
		if params.Head != nil && params.Head.Type == Symbol && params.Head.Sym == "&rest" {
			// Next param gets all remaining args as a list
			params = params.Tail
			if params == nilExpr {
				panic("macro: &rest requires a parameter name")
			}
			// Build a list from remaining args
			restList := nilExpr
			for j := len(args) - 1; j >= i; j-- {
				restList = pair(args[j], restList)
			}
			newEnv.Define(params.Head.Sym, restList)
			break
		}

		if i >= len(args) {
			panic("macro: not enough arguments")
		}
		newEnv.Define(params.Head.Sym, args[i])
		params = params.Tail
		i++
	}

	if i < len(args) && params == nilExpr {
		panic("macro: too many arguments")
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
				bodyExprs := args.Tail

				// If there's only one body expression, use it directly
				// Otherwise, wrap multiple expressions in an implicit begin
				var body *Expr
				if bodyExprs.Tail == nilExpr {
					// Single expression
					body = bodyExprs.Head
				} else {
					// Multiple expressions - wrap in begin
					body = pair(makeSym("begin"), bodyExprs)
				}
				return makeLambda(params, body, env, Lambda)
			case "begin":
				var result *Expr = nilExpr
				for args != nilExpr {
					result = eval(args.Head, env)
					args = args.Tail
				}
				return result
			case "load":
				// (load "filepath.lisp")
				if args == nilExpr {
					panic("load: missing filepath argument")
				}

				// Evaluate the filepath argument (could be a variable)
				filepath := eval(args.Head, env)

				if filepath.Type != String {
					panic("load: argument must be a string")
				}

				content, err := os.ReadFile(filepath.Str)
				if err != nil {
					panic(fmt.Sprintf("load: cannot read file %s: %v", filepath.Str, err))
				}

				exprs := readMultipleExprs(string(content))

				var result *Expr = nilExpr
				for _, expr := range exprs {
					result = eval(expr, env)
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
			i := 0
			for params != nilExpr {
				// Check for &rest parameter
				if params.Head != nil && params.Head.Type == Symbol && params.Head.Sym == "&rest" {
					// Next param gets all remaining args as a list
					params = params.Tail
					if params == nilExpr {
						panic("lambda: &rest requires a parameter name")
					}
					// Build a list from remaining args
					restList := nilExpr
					for j := len(evaledArgs) - 1; j >= i; j-- {
						restList = pair(evaledArgs[j], restList)
					}
					newEnv.Define(params.Head.Sym, restList)
					break
				}

				if i >= len(evaledArgs) {
					panic("not enough arguments")
				}
				newEnv.Define(params.Head.Sym, evaledArgs[i])
				params = params.Tail
				i++
			}

			if i < len(evaledArgs) && params == nilExpr {
				panic("too many arguments")
			}

			return eval(fn.Body, newEnv)
		}

		panic(fmt.Sprintf("not a function: %s", printExpr(fn)))
	default:
		return e
	}
}
