package main

type ExprType string

const (
	Nil     ExprType = "Nil"
	Bool    ExprType = "Bool"
	Number  ExprType = "Number"
	String  ExprType = "String"
	Symbol  ExprType = "Symbol"
	Pair    ExprType = "Pair"
	Hash    ExprType = "Hash"
	Builtin ExprType = "Builtin"
	Lambda  ExprType = "Lambda"
	Macro   ExprType = "Macro"
)

type Expr struct {
	Type      ExprType
	Num       int
	Sym       string
	Str       string
	Head      *Expr
	Tail      *Expr
	HashTable map[string]*Expr
	Fn        func([]*Expr) *Expr
	Params    *Expr
	Body      *Expr
	Env       *Env
}

var nilExpr = &Expr{Type: Nil}
var trueExpr = &Expr{Type: Bool}
var falseExpr = &Expr{Type: Bool}

// Basic construtors for the various types
func makeNum(n int) *Expr {
	return &Expr{Type: Number, Num: n}
}

func makeStr(s string) *Expr {
	return &Expr{Type: String, Str: s}
}

func makeSym(s string) *Expr {
	if s == "nil" {
		return nilExpr
	}
	if s == "true" {
		return trueExpr
	}
	if s == "false" {
		return falseExpr
	}
	return &Expr{Type: Symbol, Sym: s}
}

func pair(head, tail *Expr) *Expr {
	return &Expr{Type: Pair, Head: head, Tail: tail}
}

func makeBuiltin(fn func([]*Expr) *Expr) *Expr {
	return &Expr{Type: Builtin, Fn: fn}
}

func makeLambda(params, body *Expr, env *Env, typ ExprType) *Expr {
	return &Expr{Type: typ, Params: params, Body: body, Env: env}
}

// some helpers Ill need for lists

func list(exprs ...*Expr) *Expr {
	result := nilExpr

	for i := len(exprs) - 1; i >= 0; i-- {
		result = pair(exprs[i], result)
	}
	return result
}

func listToSlice(e *Expr) []*Expr {
	var result []*Expr
	for e != nilExpr && e.Type == Pair {
		result = append(result, e.Head)
		e = e.Tail
	}
	return result
}

// Creates an empty hash
func makeHash() *Expr {
	return &Expr{
		Type:      Hash,
		HashTable: make(map[string]*Expr),
	}
}

// Set the hash key-values
func hashSet(hash *Expr, key string, value *Expr) {
	if hash.Type != Hash {
		panic("hashSet: not a hash")
	}
	hash.HashTable[key] = value
}

// Get hash values by key
func hashGet(hash *Expr, key string) (*Expr, bool) {
	if hash.Type != Hash {
		panic("hashGet: not a hash")
	}
	val, ok := hash.HashTable[key]
	return val, ok
}

// Get all keys from a hash
func hashKeys(hash *Expr) []string {
	if hash.Type != Hash {
		panic("hashKeys: not a hash")
	}
	keys := make([]string, 0, len(hash.HashTable))
	for k := range hash.HashTable {
		keys = append(keys, k)
	}
	return keys
}
