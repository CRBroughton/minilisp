package main

import (
	"fmt"
	"os"
)

const (
	colourReset  = "\033[0m"
	colourRed    = "\033[31m"
	colourGreen  = "\033[32m"
	colourYellow = "\033[33m"
	colourBlue   = "\033[34m"
	colourPurple = "\033[35m"
	colourCyan   = "\033[36m"
	colourGray   = "\033[37m"
)

var useColours = true

func init() {
	if os.Getenv("NO_COLOR") != "" {
		useColours = false
	}
}

func colourise(colour, text string) string {
	if !useColours {
		return text
	}
	return colour + text + colourReset
}

func printResult(expr *Expr) {
	output := printExprColoured(expr)
	fmt.Println("=>", output)
}

func printExprColoured(e *Expr) string {
	if e == nil || e == nilExpr {
		return colourise(colourGray, "nil")
	}
	if e == trueExpr {
		return colourise(colourGreen, "true")
	}
	if e == falseExpr {
		return colourise(colourRed, "false")
	}

	switch e.Type {
	case Number:
		return colourise(colourCyan, printExpr(e))
	case String:
		return colourise(colourYellow, printExpr(e))
	case Symbol:
		return colourise(colourPurple, printExpr(e))
	case Builtin, Lambda, Macro:
		return colourise(colourBlue, printExpr(e))
	case Hash:
		return colourise(colourGreen, printExpr(e))
	case Pair:
		// For lists, colourise the structure but not individual elements
		return printExpr(e)
	default:
		return printExpr(e)
	}
}
