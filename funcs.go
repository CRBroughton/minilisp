package main

import "fmt"

func builtinAdd(args []*Expr) *Expr {
	sum := 0
	for _, arg := range args {
		sum += arg.Num
	}
	return makeNum(sum)
}

func builtinSub(args []*Expr) *Expr {
	if len(args) == 0 {
		return makeNum(0)
	}
	result := args[0].Num
	for i := 1; i < len(args); i++ {
		result -= args[i].Num
	}
	return makeNum(result)
}

func builtinMul(args []*Expr) *Expr {
	result := 1
	for _, arg := range args {
		result *= arg.Num
	}
	return makeNum(result)
}

func builtinDiv(args []*Expr) *Expr {
	return makeNum(args[0].Num / args[1].Num)
}

func builtinEq(args []*Expr) *Expr {
	a, b := args[0], args[1]

	if a.Type != b.Type {
		return nilExpr
	}

	switch a.Type {
	case Number:
		if a.Num == b.Num {
			return trueExpr
		}
	case Symbol:
		if a.Sym == b.Sym {
			return trueExpr
		}
	case Nil:
		return trueExpr
	default:
		if a == b {
			return trueExpr
		}
	}
	return nilExpr
}

// TODO - implement Gt function, also
// maybe not nil for falsey?
func builtinLt(args []*Expr) *Expr {
	if args[0].Num < args[1].Num {
		return trueExpr
	}
	return nilExpr
}

func builtinPair(args []*Expr) *Expr {
	return pair(args[0], args[1])
}

func builtinHead(args []*Expr) *Expr {
	return args[0].Head
}

func builtinTail(args []*Expr) *Expr {
	return args[0].Tail
}

func builtinNullP(args []*Expr) *Expr {
	if args[0] == nilExpr {
		return trueExpr
	}
	return nilExpr
}

func builtinPrint(args []*Expr) *Expr {
	fmt.Println(printExpr(args[0]))
	return args[0]
}

func builtinHash(args []*Expr) *Expr {
	if len(args)%2 != 0 {
		panic("hash: expect an even number of arguments")
	}

	h := makeHash()
	for i := 0; i < len(args); i += 2 {
		key := args[i]
		value := args[i+1]

		if key.Type != String {
			panic("hash: keys must be strings")
		}
		hashSet(h, key.Str, value)
	}

	return h
}

func builtinHashGet(args []*Expr) *Expr {
	if len(args) != 2 {
		panic("hash-get: expects 2 arguments")
	}

	hash := args[0]
	key := args[1]

	if key.Type != String {
		panic("hash-get: key must be a string")
	}

	val, ok := hashGet(hash, key.Str)
	if !ok {
		return nilExpr
	}
	return val
}

func builtinHashSet(args []*Expr) *Expr {
	if len(args) != 3 {
		panic("hash-set: expects 3 arguments (hash, key, value)")
	}

	hash := args[0]
	key := args[1]
	value := args[2]

	if key.Type != String {
		panic("hash-set: key must be a string")
	}

	hashSet(hash, key.Str, value)
	return hash
}

func builtinHashKeys(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("hash-keys: expects 1 argument")
	}

	keys := hashKeys(args[0])
	result := nilExpr

	// Build list of string keys
	for i := len(keys) - 1; i >= 0; i-- {
		result = pair(makeStr(keys[i]), result)
	}

	return result
}
