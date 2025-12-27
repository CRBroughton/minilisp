package main

import "testing"

func TestBuiltinAdd(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))

	tests := []struct {
		input string
		want  int
	}{
		{"(+ 1 2)", 3},
		{"(+ 1 2 3)", 6},
		{"(+ 1 2 3 4 5)", 15},
		{"(+ 0)", 0},
		{"(+ -5 5)", 0},
		{"(+ 100 -50)", 50},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		if result.Type != Number {
			t.Errorf("%s: type = %v, want Number", tt.input, result.Type)
		}
		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestBuiltinSub(t *testing.T) {
	env := NewEnv(nil)
	env.Define("-", makeBuiltin(builtinSub))

	tests := []struct {
		input string
		want  int
	}{
		{"(- 10 5)", 5},
		{"(- 100 50 25)", 25},
		{"(- 0 10)", -10},
		{"(- 5 5)", 0},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestBuiltinMul(t *testing.T) {
	env := NewEnv(nil)
	env.Define("*", makeBuiltin(builtinMul))

	tests := []struct {
		input string
		want  int
	}{
		{"(* 2 3)", 6},
		{"(* 2 3 4)", 24},
		{"(* 5)", 5},
		{"(* 10 0)", 0},
		{"(* -2 3)", -6},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)
		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestBuiltinDiv(t *testing.T) {
	env := NewEnv(nil)
	env.Define("/", makeBuiltin(builtinDiv))

	tests := []struct {
		input string
		want  int
	}{
		{"(/ 10 2)", 5},
		{"(/ 100 10)", 10},
		{"(/ 7 2)", 3}, // Integer division
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)
		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestBuiltinEq(t *testing.T) {
	env := NewEnv(nil)
	env.Define("=", makeBuiltin(builtinEq))

	tests := []struct {
		input    string
		wantTrue bool
	}{
		{"(= 5 5)", true},
		{"(= 5 6)", false},
		{"(= 0 0)", true},
		{"(= -5 -5)", true},
		{"(= 10 5)", false},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		isTrue := result == trueExpr
		if isTrue != tt.wantTrue {
			t.Errorf("%s = %v, want %v", tt.input, isTrue, tt.wantTrue)
		}
	}
}

func TestBuiltinLt(t *testing.T) {
	env := NewEnv(nil)
	env.Define("<", makeBuiltin(builtinLt))

	tests := []struct {
		input    string
		wantTrue bool
	}{
		{"(< 3 5)", true},
		{"(< 5 3)", false},
		{"(< 5 5)", false},
		{"(< -10 0)", true},
		{"(< 0 -10)", false},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		isTrue := result == trueExpr
		if isTrue != tt.wantTrue {
			t.Errorf("%s = %v, want %v", tt.input, isTrue, tt.wantTrue)
		}
	}
}

func TestBuiltinPairs(t *testing.T) {
	env := NewEnv(nil)
	env.Define("pair", makeBuiltin(builtinPair))

	expr := readStr("(pair 1 2)")
	result := eval(expr, env)

	if result.Type != Pair {
		t.Fatalf("pair result type = %v, want Pair", result.Type)
	}
	if result.Head.Num != 1 {
		t.Errorf("head = %d, want 1", result.Head.Num)
	}
	if result.Tail.Num != 2 {
		t.Errorf("tail = %d, want 2", result.Tail.Num)
	}
}

func TestBuiltinHead(t *testing.T) {
	env := NewEnv(nil)
	env.Define("head", makeBuiltin(builtinHead))
	env.Define("pair", makeBuiltin(builtinPair))

	expr := readStr("(head (pair 1 2))")
	result := eval(expr, env)

	if result.Num != 1 {
		t.Errorf("head = %d, want 1", result.Num)
	}
}

func TestBuiltinTail(t *testing.T) {
	env := NewEnv(nil)
	env.Define("tail", makeBuiltin(builtinTail))
	env.Define("pair", makeBuiltin(builtinPair))

	expr := readStr("(tail (pair 1 2))")
	result := eval(expr, env)

	if result.Num != 2 {
		t.Errorf("tail = %d, want 2", result.Num)
	}
}

func TestBuiltinNullP(t *testing.T) {
	env := NewEnv(nil)
	env.Define("null?", makeBuiltin(builtinNullP))

	tests := []struct {
		input    string
		wantTrue bool
	}{
		{"(null? nil)", true},
		{"(null? 42)", false},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		isTrue := result == trueExpr
		if isTrue != tt.wantTrue {
			t.Errorf("%s = %v, want %v", tt.input, isTrue, tt.wantTrue)
		}
	}
}

func TestEvalList(t *testing.T) {
	env := NewEnv(nil)
	env.Define("x", makeNum(10))
	env.Define("y", makeNum(20))

	// Evaluate (x y) should give [10, 20]
	lst := list(makeSym("x"), makeSym("y"))
	results := evalList(lst, env)

	if len(results) != 2 {
		t.Fatalf("evalList length = %d, want 2", len(results))
	}
	if results[0].Num != 10 || results[1].Num != 20 {
		t.Errorf("evalList = [%d, %d], want [10, 20]", results[0].Num, results[1].Num)
	}
}

func TestComplexArithmetic(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))
	env.Define("-", makeBuiltin(builtinSub))

	tests := []struct {
		input string
		want  int
	}{
		{"(+ (* 2 3) (* 4 5))", 26}, // (+ 6 20) = 26
		{"(+ 1 (+ 2 (+ 3 4)))", 10}, // 1 + 2 + 3 + 4
		{"(- (+ 10 5) (* 2 3))", 9}, // 15 - 6
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)
		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestBuiltinHash(t *testing.T) {
	// Empty hash
	result := builtinHash([]*Expr{})
	if result.Type != Hash {
		t.Errorf("builtinHash() type = %v, want Hash", result.Type)
	}
	if len(result.HashTable) != 0 {
		t.Errorf("empty hash should have 0 entries")
	}

	// Hash with key-value pairs
	result = builtinHash([]*Expr{
		makeStr("name"), makeStr("Alice"),
		makeStr("age"), makeNum(30),
	})

	if result.Type != Hash {
		t.Fatalf("type = %v, want Hash", result.Type)
	}

	name, ok := hashGet(result, "name")
	if !ok || name.Str != "Alice" {
		t.Errorf("name should be 'Alice'")
	}

	age, ok := hashGet(result, "age")
	if !ok || age.Num != 30 {
		t.Errorf("age should be 30")
	}
}

func TestBuiltinHashOddArgs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("builtinHash with odd number of args should panic")
		}
	}()

	builtinHash([]*Expr{makeStr("key")})
}

func TestBuiltinHashNonStringKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("builtinHash with non-string key should panic")
		}
	}()

	builtinHash([]*Expr{makeNum(42), makeStr("value")})
}

func TestBuiltinHashGet(t *testing.T) {
	hash := makeHash()
	hashSet(hash, "name", makeStr("Alice"))
	hashSet(hash, "age", makeNum(30))

	// Get existing key
	result := builtinHashGet([]*Expr{hash, makeStr("name")})
	if result.Type != String || result.Str != "Alice" {
		t.Errorf("hash-get name = %v, want 'Alice'", result.Str)
	}

	result = builtinHashGet([]*Expr{hash, makeStr("age")})
	if result.Type != Number || result.Num != 30 {
		t.Errorf("hash-get age = %v, want 30", result.Num)
	}

	// Get missing key
	result = builtinHashGet([]*Expr{hash, makeStr("missing")})
	if result != nilExpr {
		t.Error("hash-get with missing key should return nil")
	}
}

func TestBuiltinHashGetWrongArgs(t *testing.T) {
	tests := []struct {
		name string
		args []*Expr
	}{
		{"no args", []*Expr{}},
		{"one arg", []*Expr{makeHash()}},
		{"three args", []*Expr{makeHash(), makeStr("key"), makeStr("extra")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("hash-get with %s should panic", tt.name)
				}
			}()
			builtinHashGet(tt.args)
		})
	}
}

func TestBuiltinHashGetNonStringKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("hash-get with non-string key should panic")
		}
	}()

	hash := makeHash()
	builtinHashGet([]*Expr{hash, makeNum(42)})
}

func TestBuiltinHashSet(t *testing.T) {
	hash := makeHash()

	// Set new value
	result := builtinHashSet([]*Expr{hash, makeStr("name"), makeStr("Alice")})

	// Should return the hash
	if result != hash {
		t.Error("hash-set should return the same hash")
	}

	// Check value was set
	val, ok := hashGet(hash, "name")
	if !ok || val.Str != "Alice" {
		t.Error("value should be set")
	}

	// Overwrite existing value
	builtinHashSet([]*Expr{hash, makeStr("name"), makeStr("Bob")})
	val, _ = hashGet(hash, "name")
	if val.Str != "Bob" {
		t.Errorf("name should be overwritten to 'Bob', got %v", val.Str)
	}
}

func TestBuiltinHashSetWrongArgs(t *testing.T) {
	tests := []struct {
		name string
		args []*Expr
	}{
		{"no args", []*Expr{}},
		{"one arg", []*Expr{makeHash()}},
		{"two args", []*Expr{makeHash(), makeStr("key")}},
		{"four args", []*Expr{makeHash(), makeStr("k"), makeStr("v"), makeStr("extra")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("hash-set with %s should panic", tt.name)
				}
			}()
			builtinHashSet(tt.args)
		})
	}
}

func TestBuiltinHashSetNonStringKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("hash-set with non-string key should panic")
		}
	}()

	hash := makeHash()
	builtinHashSet([]*Expr{hash, makeNum(42), makeStr("value")})
}

func TestBuiltinHashKeys(t *testing.T) {
	hash := makeHash()
	hashSet(hash, "name", makeStr("Alice"))
	hashSet(hash, "age", makeNum(30))
	hashSet(hash, "active", trueExpr)

	result := builtinHashKeys([]*Expr{hash})

	// Should return a list
	if result.Type != Pair && result != nilExpr {
		t.Fatalf("hash-keys should return a list, got %v", result.Type)
	}

	// Convert to slice and check
	keys := listToSlice(result)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}

	// Check all keys are strings
	keyStrings := make(map[string]bool)
	for _, key := range keys {
		if key.Type != String {
			t.Errorf("key should be String, got %v", key.Type)
		}
		keyStrings[key.Str] = true
	}

	// Check expected keys are present
	expectedKeys := []string{"name", "age", "active"}
	for _, expected := range expectedKeys {
		if !keyStrings[expected] {
			t.Errorf("expected key %q not found", expected)
		}
	}
}

func TestBuiltinHashKeysEmpty(t *testing.T) {
	hash := makeHash()
	result := builtinHashKeys([]*Expr{hash})

	if result != nilExpr {
		t.Errorf("hash-keys on empty hash should return nil, got %v", printExpr(result))
	}
}

func TestBuiltinHashKeysWrongArgs(t *testing.T) {
	tests := []struct {
		name string
		args []*Expr
	}{
		{"no args", []*Expr{}},
		{"two args", []*Expr{makeHash(), makeHash()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("hash-keys with %s should panic", tt.name)
				}
			}()
			builtinHashKeys(tt.args)
		})
	}
}

func TestHashWithDifferentTypes(t *testing.T) {
	hash := builtinHash([]*Expr{
		makeStr("string"), makeStr("hello"),
		makeStr("number"), makeNum(42),
		makeStr("bool"), trueExpr,
		makeStr("nil"), nilExpr,
		makeStr("list"), list(makeNum(1), makeNum(2)),
		makeStr("nested"), makeHash(),
	})

	// Verify all types stored correctly
	str, _ := hashGet(hash, "string")
	if str.Type != String || str.Str != "hello" {
		t.Error("string value incorrect")
	}

	num, _ := hashGet(hash, "number")
	if num.Type != Number || num.Num != 42 {
		t.Error("number value incorrect")
	}

	boolean, _ := hashGet(hash, "bool")
	if boolean != trueExpr {
		t.Error("bool value incorrect")
	}

	nilVal, _ := hashGet(hash, "nil")
	if nilVal != nilExpr {
		t.Error("nil value incorrect")
	}

	lst, _ := hashGet(hash, "list")
	if lst.Type != Pair {
		t.Error("list value incorrect")
	}

	nested, _ := hashGet(hash, "nested")
	if nested.Type != Hash {
		t.Error("nested hash incorrect")
	}
}

func TestHashMutation(t *testing.T) {
	// Create hash
	hash := builtinHash([]*Expr{
		makeStr("count"), makeNum(0),
	})

	// Mutate multiple times
	for i := 1; i <= 5; i++ {
		builtinHashSet([]*Expr{hash, makeStr("count"), makeNum(i)})
	}

	// Check final value
	count, _ := hashGet(hash, "count")
	if count.Num != 5 {
		t.Errorf("count should be 5, got %d", count.Num)
	}
}

func TestHashFromLisp(t *testing.T) {
	env := NewEnv(nil)
	env.Define("hash", makeBuiltin(builtinHash))
	env.Define("hash-get", makeBuiltin(builtinHashGet))
	env.Define("hash-set", makeBuiltin(builtinHashSet))
	env.Define("hash-keys", makeBuiltin(builtinHashKeys))

	// Create hash from Lisp
	code := `(hash "name" "Alice" "age" 30)`
	expr := readStr(code)
	hash := eval(expr, env)

	if hash.Type != Hash {
		t.Fatalf("eval should produce Hash, got %v", hash.Type)
	}

	// Get value from Lisp
	getCode := `(hash-get (hash "x" 10) "x")`
	result := eval(readStr(getCode), env)
	if result.Num != 10 {
		t.Errorf("hash-get should return 10, got %d", result.Num)
	}

	// Set value from Lisp - test hash-set directly
	h := makeHash()
	env.Define("h", h)
	eval(readStr(`(hash-set h "b" 2)`), env)

	val, _ := hashGet(h, "b")
	if val.Num != 2 {
		t.Error("hash-set from Lisp should work")
	}
}
