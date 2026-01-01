package main

import (
	"strconv"
	"unicode"
)

type Reader struct {
	input string
	pos   int
}

func (r *Reader) skipWhitespace() {
	for r.pos < len(r.input) {
		ch := r.input[r.pos]

		// Skip whitespace
		if unicode.IsSpace(rune(ch)) {
			r.pos++
			continue
		}

		// Skip comments (semicolon to end of line)
		if ch == ';' {
			for r.pos < len(r.input) && r.input[r.pos] != '\n' {
				r.pos++
			}
			continue
		}

		break
	}
}

func (r *Reader) peek() byte {
	if r.pos >= len(r.input) {
		return 0
	}
	return r.input[r.pos]
}

func (r *Reader) next() byte {
	if r.pos >= len(r.input) {
		return 0
	}
	ch := r.input[r.pos]
	r.pos++
	return ch
}

func (r *Reader) readExpr() *Expr {
	r.skipWhitespace()
	ch := r.peek()

	if ch == 0 {
		return nilExpr
	}

	if ch == '"' {
		return r.readStr()
	}

	// Quote sugar: 'x → (quote x)
	if ch == '\'' {
		r.next()
		return list(makeSym("quote"), r.readExpr())
	}

	// List: (...)
	if ch == '(' {
		r.next()
		return r.readList()
	}

	// Error: unexpected closing paren
	if ch == ')' {
		panic("unexpected )")
	}

	// Number: 42, -10
	if unicode.IsDigit(rune(ch)) || (ch == '-' && r.pos+1 < len(r.input) && unicode.IsDigit(rune(r.input[r.pos+1]))) {
		return r.readNumber()
	}

	// Symbol: x, +, factorial
	return r.readSymbol()
}

func (r *Reader) readStr() *Expr {
	r.next() // skip the opening quote
	start := r.pos

	for {
		ch := r.peek()
		if ch == 0 {
			panic("unterminated string")
		}
		if ch == '"' {
			break
		}
		// TODO - Handle escape sequences
		r.next()
	}
	str := r.input[start:r.pos]
	r.next()
	return makeStr(str)

}

func (r *Reader) readList() *Expr {
	r.skipWhitespace()

	// Empty list: ()
	if r.peek() == ')' {
		r.next()
		return nilExpr
	}

	// Read elements until )
	head := r.readExpr()
	tail := r.readList()
	return pair(head, tail)
}

func (r *Reader) readNumber() *Expr {
	start := r.pos

	// Handle negative sign
	if r.peek() == '-' {
		r.next()
	}

	// Read digits
	for unicode.IsDigit(rune(r.peek())) {
		r.next()
	}

	numStr := r.input[start:r.pos]
	n, _ := strconv.Atoi(numStr)
	return makeNum(n)
}

func (r *Reader) readSymbol() *Expr {
	start := r.pos

	// Read until whitespace or special character
	for {
		ch := r.peek()
		if ch == 0 || unicode.IsSpace(rune(ch)) || ch == '(' || ch == ')' {
			break
		}
		r.next()
	}

	sym := r.input[start:r.pos]
	return makeSym(sym)
}

// Helper function to read from a string
func readStr(s string) *Expr {
	r := &Reader{input: s}
	return r.readExpr()
}

// Read multiple expressions from input (used for loading files and piped input)
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
