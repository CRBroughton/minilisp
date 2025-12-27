package main

import (
	"fmt"
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
	case Hash:
		return printHash(e)
	case Number:
		return strconv.Itoa(e.Num)
	case String:
		return fmt.Sprintf("\"%s\"", e.Str)
	case Symbol:
		return e.Sym
	case Builtin:
		return "<builtin>"
	case Lambda:
		return "<lambda>"
	case Macro:
		return "<macro>"
	case Pair:
		return printList(e)
	default:
		return "<unknown>"
	}

}

func printHash(e *Expr) string {
	if len(e.HashTable) == 0 {
		return "{}"
	}

	parts := []string{}
	for k, v := range e.HashTable {
		parts = append(parts, fmt.Sprintf("%q: %s", k, printExpr(v)))
	}

	return "{" + strings.Join(parts, ", ") + "}"
}

func printList(e *Expr) string {
	var parts []string

	// Collect all elements
	for e != nilExpr && e.Type == Pair {
		parts = append(parts, printExpr(e.Head))
		e = e.Tail
	}

	// Check if we have an improper list (1 . 2)
	if e != nilExpr {
		parts = append(parts, ".", printExpr(e))
	}

	return "(" + strings.Join(parts, " ") + ")"
}
