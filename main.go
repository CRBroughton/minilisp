package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func loadStdLib(env *Env) {
	// Load macros (thread macros, when, cond)
	eval(readStr(`(load "std/macro.lisp")`), env)

	// Load functions (factorial, sum)
	eval(readStr(`(load "std/functions.lisp")`), env)

	// Result type
	eval(readStr(`(load "std/result.lisp")`), env)

	eval(readStr(`(load "std/html.lisp")`), env)

}

// Read all input from stdin (for piped input)
func readAllInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}

func main() {
	// Create global environment
	env := NewEnv(nil)

	// Define built-ins
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("-", makeBuiltin(builtinSub))
	env.Define("*", makeBuiltin(builtinMul))
	env.Define("/", makeBuiltin(builtinDiv))
	env.Define("=", makeBuiltin(builtinEq))
	env.Define("<", makeBuiltin(builtinLt))
	env.Define("<=", makeBuiltin(builtinEqualOrLt))
	env.Define(">", makeBuiltin(builtinGt))
	env.Define(">=", makeBuiltin(builtinEqualOrGt))
	env.Define("pair", makeBuiltin(builtinPair))
	env.Define("list", makeBuiltin(builtinList))
	env.Define("head", makeBuiltin(builtinHead))
	env.Define("tail", makeBuiltin(builtinTail))
	env.Define("null?", makeBuiltin(builtinNullP))
	env.Define("print", makeBuiltin(builtinPrint))
	env.Define("hash", makeBuiltin(builtinHash))
	env.Define("hash-get", makeBuiltin(builtinHashGet))
	env.Define("hash-set", makeBuiltin(builtinHashSet))
	env.Define("hash-keys", makeBuiltin(builtinHashKeys))
	env.Define("fetch", makeBuiltin(builtinFetch))
	env.Define("json-stringify", makeBuiltin(builtinJsonStringify))
	env.Define("string-append", makeBuiltin(builtinStringAppend))

	env.Define("number?", makeBuiltin(builtinNumberP))
	env.Define("string?", makeBuiltin(builtinStringP))
	env.Define("symbol?", makeBuiltin(builtinSymbolP))
	env.Define("list?", makeBuiltin(builtinListP))
	env.Define("bool?", makeBuiltin(builtinBoolP))

	env.Define("@json", makeBuiltin(builtinJsonParse))
	env.Define("@string", makeBuiltin(builtinToString))
	env.Define("@number", makeBuiltin(builtinToNumber))

	env.Define("http-server", makeBuiltin(builtinHttpServer))

	env.Define("string-join", makeBuiltin(builtinStringJoin))
	env.Define("html-escape", makeBuiltin(builtinHtmlEscape))

	// Bootstrap defmacro
	defmacroCode := "(define defmacro (macro (name params body) (pair 'define (pair name (pair (pair 'macro (pair params (pair body nil))) nil)))))"
	eval(readStr(defmacroCode), env)

	// Load standard library
	loadStdLib(env)

	// Check if input is from pipe/file or interactive
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Piped input - read all at once
		input := readAllInput()
		exprs := readMultipleExprs(input)

		for _, expr := range exprs {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Error: %v\n", r)
					}
				}()
				result := eval(expr, env)
				_ = result // Don't print unless using print
			}()
		}
	} else {
		startREPL(env)
	}
}
