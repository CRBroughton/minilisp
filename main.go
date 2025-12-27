package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
	env.Define("json-parse", makeBuiltin(builtinJsonParse))
	env.Define("json-stringify", makeBuiltin(builtinJsonStringify))
	env.Define("string-append", makeBuiltin(builtinStringAppend))

	// Bootstrap defmacro
	defmacroCode := "(define defmacro (macro (name params body) (pair 'define (pair name (pair (pair 'macro (pair params (pair body nil))) nil)))))"
	eval(readStr(defmacroCode), env)

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
		// Interactive REPL
		fmt.Println("MiniLisp - Type expressions (Ctrl+D to exit)")
		scanner := bufio.NewScanner(os.Stdin)

		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}

			line := scanner.Text()
			line = strings.TrimSpace(line)

			// Skip empty lines and comments
			if line == "" || strings.HasPrefix(line, ";") {
				continue
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Error: %v\n", r)
					}
				}()

				expr := readStr(line)
				if expr == nilExpr {
					return
				}
				result := eval(expr, env)
				fmt.Println(printExpr(result))
			}()
		}
	}
}
