package main

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"
)

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
	case String:
		if a.Str == b.Str {
			return trueExpr
		}
	case Nil:
		return trueExpr
	case Pair:
		// Structural equality for lists
		if structuralEq(a, b) {
			return trueExpr
		}
	default:
		if a == b {
			return trueExpr
		}
	}
	return nilExpr
}

func structuralEq(a, b *Expr) bool {
	// Both nil
	if a == nilExpr && b == nilExpr {
		return true
	}

	// One nil, one not
	if a == nilExpr || b == nilExpr {
		return false
	}

	// Different types
	if a.Type != b.Type {
		return false
	}

	// For pairs, recursively check head and tail
	if a.Type == Pair {
		return structuralEq(a.Head, b.Head) && structuralEq(a.Tail, b.Tail)
	}

	// For atoms, check value equality
	switch a.Type {
	case Number:
		return a.Num == b.Num
	case Symbol:
		return a.Sym == b.Sym
	case String:
		return a.Str == b.Str
	case Bool:
		return a == b // trueExpr is a singleton
	default:
		return a == b // Pointer equality for other types
	}
}

// TODO -  maybe not nil for falsey?
func builtinLt(args []*Expr) *Expr {
	if args[0].Num < args[1].Num {
		return trueExpr
	}
	return nilExpr
}

func builtinGt(args []*Expr) *Expr {
	if args[0].Num > args[1].Num {
		return trueExpr
	}
	return nilExpr
}

func builtinPair(args []*Expr) *Expr {
	return pair(args[0], args[1])
}

func builtinList(args []*Expr) *Expr {
	result := nilExpr
	for i := len(args) - 1; i >= 0; i-- {
		result = pair(args[i], result)
	}
	return result
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

func builtinStringAppend(args []*Expr) *Expr {
	var result string
	for _, arg := range args {
		if arg.Type != String {
			panic("string-append: all arguments must be strings")
		}
		result += arg.Str
	}
	return makeStr(result)
}

func builtinJsonParse(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("@json: expects 1 argument")
	}

	if args[0].Type != String {
		panic("@json: argument must be a string")
	}

	var data interface{}
	err := json.Unmarshal([]byte(args[0].Str), &data)
	if err != nil {
		panic(fmt.Sprintf("@json: %v", err))
	}

	return jsonToExpr(data)
}

func jsonToExpr(data interface{}) *Expr {
	switch v := data.(type) {
	case nil:
		return nilExpr

	case bool:
		if v {
			return trueExpr
		}
		return nilExpr

	case float64:
		return makeNum(int(v))

	case string:
		return makeStr(v)

	case []interface{}:
		// JSON array → Lisp list
		result := nilExpr
		for i := len(v) - 1; i >= 0; i-- {
			result = pair(jsonToExpr(v[i]), result)
		}
		return result

	case map[string]interface{}:
		// JSON object → Hash
		hash := makeHash()
		for key, val := range v {
			hashSet(hash, key, jsonToExpr(val))
		}
		return hash

	default:
		panic(fmt.Sprintf("@json: unsupported type %T", v))
	}
}

func builtinJsonStringify(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("json-stringify: expects 1 argument")
	}

	data := exprToJson(args[0])

	bytes, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("json-stringify: %v", err))
	}

	return makeStr(string(bytes))
}

func exprToJson(e *Expr) interface{} {
	switch e.Type {
	case Nil:
		return nil

	case Bool:
		return e == trueExpr

	case Number:
		return e.Num

	case String:
		return e.Str

	case Pair:
		// Lisp list → JSON array
		items := listToSlice(e)
		result := make([]interface{}, len(items))
		for i, item := range items {
			result[i] = exprToJson(item)
		}
		return result

	case Hash:
		// Hash → JSON object
		result := make(map[string]interface{})
		for key, val := range e.HashTable {
			result[key] = exprToJson(val)
		}
		return result

	default:
		panic(fmt.Sprintf("json-stringify: cannot convert type %v", e.Type))
	}
}

// Type-checkers, might split out into own file later, maybe

func builtinNumberP(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("number?: expect 1 argument")
	}
	if args[0].Type == Number {
		return trueExpr
	}
	return nilExpr
}

func builtinStringP(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("string?: expect 1 argument")
	}
	if args[0].Type == String {
		return trueExpr
	}
	return nilExpr
}

func builtinSymbolP(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("symbol?: expect 1 argument")
	}
	if args[0].Type == Symbol {
		return trueExpr
	}
	return nilExpr
}

func builtinListP(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("list?: expect 1 argument")
	}
	if args[0].Type == Pair {
		return trueExpr
	}
	return nilExpr
}

func builtinBoolP(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("bool?: expect 1 argument")
	}
	// In MiniLisp, true is a special value, nil is also considered bool
	if args[0] == trueExpr || args[0] == nilExpr {
		return trueExpr
	}
	return nilExpr
}

func builtinToString(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("@string: expect 1 argument")
	}

	val := args[0]

	switch val.Type {
	case Number:
		return makeStr(fmt.Sprintf("%d", val.Num))

	case String:
		return val

	case Symbol:
		return makeStr(val.Sym)

	case Nil:
		return makeStr("nil")

	case Bool:
		if val == trueExpr {
			return makeStr("true")
		}
		return makeStr("false")

	case Pair:
		return makeStr(printExpr(val))

	default:
		return makeStr(printExpr(val))
	}
}

func builtinToNumber(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("@number: expect 1 argument")
	}

	val := args[0]

	switch val.Type {
	case Number:
		return val

	case String:
		num := 0
		_, err := fmt.Sscanf(val.Str, "%d", &num)
		if err != nil {
			panic(fmt.Sprintf("@number: cannot parse '%s' as number", val.Str))
		}
		return makeNum(num)

	default:
		panic(fmt.Sprintf("@number: cannot convert %s to number", val.Type))
	}
}

func builtinStringJoin(args []*Expr) *Expr {
	if len(args) != 2 {
		panic("string-join: expects 2 arguments (list, separator)")
	}

	items := listToSlice(args[0])
	sep := ""
	if args[1].Type == String {
		sep = args[1].Str
	}

	parts := make([]string, len(items))
	for i, item := range items {
		if item.Type == String {
			parts[i] = item.Str
		} else {
			parts[i] = printExpr(item)
		}
	}

	return makeStr(strings.Join(parts, sep))
}

func builtinHtmlEscape(args []*Expr) *Expr {
	if len(args) != 1 {
		panic("html-escape: expects 1 argument")
	}

	if args[0].Type != String {
		panic("html-escape: argument must be a string")
	}

	escaped := html.EscapeString(args[0].Str)
	return makeStr(escaped)
}
