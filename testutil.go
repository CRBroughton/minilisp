package main

// setupTestEnv creates a test environment with all necessary built-ins
// This is used across test files to provide a consistent testing environment
func setupTestEnv() *Env {
	env := NewEnv(nil)

	// Arithmetic
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("-", makeBuiltin(builtinSub))
	env.Define("*", makeBuiltin(builtinMul))
	env.Define("/", makeBuiltin(builtinDiv))

	// Comparison
	env.Define("=", makeBuiltin(builtinEq))
	env.Define("!=", makeBuiltin(builtinNotEq))
	env.Define("<", makeBuiltin(builtinLt))
	env.Define("<=", makeBuiltin(builtinEqualOrLt))
	env.Define(">", makeBuiltin(builtinGt))
	env.Define(">=", makeBuiltin(builtinEqualOrGt))

	// List operations
	env.Define("pair", makeBuiltin(builtinPair))
	env.Define("head", makeBuiltin(builtinHead))
	env.Define("tail", makeBuiltin(builtinTail))
	env.Define("null?", makeBuiltin(builtinNullP))
	env.Define("list", makeBuiltin(builtinList))

	// Hash operations
	env.Define("hash", makeBuiltin(builtinHash))
	env.Define("hash-get", makeBuiltin(builtinHashGet))
	env.Define("hash-set", makeBuiltin(builtinHashSet))
	env.Define("hash-keys", makeBuiltin(builtinHashKeys))

	// Type predicates
	env.Define("number?", makeBuiltin(builtinNumberP))
	env.Define("string?", makeBuiltin(builtinStringP))
	env.Define("symbol?", makeBuiltin(builtinSymbolP))
	env.Define("list?", makeBuiltin(builtinListP))
	env.Define("bool?", makeBuiltin(builtinBoolP))

	// Type converters
	env.Define("@string", makeBuiltin(builtinToString))
	env.Define("@number", makeBuiltin(builtinToNumber))
	env.Define("@json", makeBuiltin(builtinJsonParse))

	// String operations
	env.Define("string-append", makeBuiltin(builtinStringAppend))
	env.Define("string-join", makeBuiltin(builtinStringJoin))
	env.Define("html-escape", makeBuiltin(builtinHtmlEscape))

	// HTTP
	env.Define("fetch", makeBuiltin(builtinFetch))
	env.Define("http-server", makeBuiltin(builtinHttpServer))

	// JSON
	env.Define("json-stringify", makeBuiltin(builtinJsonStringify))

	// I/O
	env.Define("print", makeBuiltin(builtinPrint))

	// Result type helpers
	env.Define("ok", makeBuiltin(builtinOk))
	env.Define("err", makeBuiltin(builtinErr))
	env.Define("ok?", makeBuiltin(builtinOkP))
	env.Define("err?", makeBuiltin(builtinErrP))
	env.Define("unwrap", makeBuiltin(builtinUnwrap))
	env.Define("unwrap-err", makeBuiltin(builtinUnwrapErr))
	env.Define("unwrap-or", makeBuiltin(builtinUnwrapOr))

	return env
}
