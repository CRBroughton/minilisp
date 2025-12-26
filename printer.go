package main

import (
	"strconv"
	"strings"
)

func printExpr(e *Expr) string {
	if e == nil {
		return "nil"
	}
	if e == nilExpr {
		return "nil"
	}

	if e == trueExpr {
		return "true"
	}
	if e == falseExpr {
		return "false"
	}

	switch e.Type {
	case Number:
		return strconv.Itoa(e.Num)
	case Symbol:
		return e.Sym
	case Builtin:
		return "<builtin>"
	case Lambda:
		return "<lambda>"
	case Macro:
		return "<macro>"
	case Cons:
		return printList(e)
	default:
		return "<unknown>"
	}

}

func printList(e *Expr) string {
	var parts []string

	// Collect all elements
	for e != nilExpr && e.Type == Cons {
		parts = append(parts, printExpr(e.Head))
		e = e.Tail
	}

	// Check if we have an improper list (1 . 2)
	if e != nilExpr {
		parts = append(parts, ".", printExpr(e))
	}

	return "(" + strings.Join(parts, " ") + ")"
}
