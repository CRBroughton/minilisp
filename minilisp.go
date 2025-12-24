package minilisp

type ExprType int

const (
	Nil ExprType = iota
	Bool
	Number
	Symbol
	Cons
	Builtin
	Lambda
	Macro
)

type Env struct{}
type Expr struct {
	Type   ExprType
	Num    int
	Sym    string
	Car    *Expr
	Cdr    *Expr
	Fn     func([]*Expr) *Expr
	Params *Expr
	Body   *Expr
	Env    *Env
}

var nilExpr = &Expr{Type: Nil}
var trueExpr = &Expr{Type: Bool}

// Basic construtors for the various types
func makeNum(n int) *Expr {
	return &Expr{Type: Number, Num: n}
}

func makeSym(s string) *Expr {
	if s == "nil" {
		return nilExpr
	}
	if s == "true" {
		return trueExpr
	}
	return &Expr{Type: Symbol, Sym: s}
}

func cons(car, cdr *Expr) *Expr {
	return &Expr{Type: Cons, Car: car, Cdr: cdr}
}
