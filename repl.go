package main

import (
	"bufio"
	"os"
	"strings"
)

func readMultipleExprs(input string) []*Expr {
	r := &Reader{input: input}
	var exprs []*Expr

	for {
		r.skipWhitespace()
		if r.pos >= len(r.input) {
			break
		}

		expr := r.readExpr()
		if expr != nilExpr || r.pos < len(r.input) {
			exprs = append(exprs, expr)
		}
	}
	return exprs
}
func readAllInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}
