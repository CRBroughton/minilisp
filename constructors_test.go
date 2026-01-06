package main

import (
	"errors"
	"testing"
)

func TestMakeHash(t *testing.T) {
	hash := makeHash()

	if hash.Type != Hash {
		t.Errorf("makeHash() type = %v, want Hash", hash.Type)
	}

	if hash.HashTable == nil {
		t.Error("makeHash() should initialise HashTable")
	}
}

func TestHashSetGet(t *testing.T) {
	hash := makeHash()

	hashSet(hash, "name", makeStr("Alice"))

	val, ok := hashGet(hash, "name")
	if !ok {
		t.Error("hashGet should find 'name'")
	}
	if val.Type != String || val.Str != "Alice" {
		t.Errorf("hashGet('name') = %v, want 'Alice'", val.Str)
	}
}

func TestHashGetMissing(t *testing.T) {
	hash := makeHash()

	_, ok := hashGet(hash, "missing")
	if ok {
		t.Error("hashGet should return false for missing key")
	}
}

func TestHashMultipleKeys(t *testing.T) {
	hash := makeHash()

	hashSet(hash, "name", makeStr("Alice"))
	hashSet(hash, "age", makeNum(30))
	hashSet(hash, "active", trueExpr)

	// Check all values
	name, _ := hashGet(hash, "name")
	age, _ := hashGet(hash, "age")
	active, _ := hashGet(hash, "active")

	if name.Str != "Alice" {
		t.Errorf("name = %v, want Alice", name.Str)
	}
	if age.Num != 30 {
		t.Errorf("age = %v, want 30", age.Num)
	}
	if active != trueExpr {
		t.Error("active should be true")
	}
}

func TestHashOverwrite(t *testing.T) {
	hash := makeHash()

	hashSet(hash, "x", makeNum(10))
	hashSet(hash, "x", makeNum(20)) // Overwrite

	val, _ := hashGet(hash, "x")
	if val.Num != 20 {
		t.Errorf("x = %v, want 20", val.Num)
	}
}

func TestResultToExprOk(t *testing.T) {
	// Create an Ok result with a string value
	value := makeStr("test value")
	result := Ok(value)

	// Convert to Expr
	expr := resultToExpr(result)

	// Verify it's a hash
	if expr.Type != Hash {
		t.Errorf("resultToExpr(Ok) type = %v, want Hash", expr.Type)
	}

	// Verify type field is "ok"
	typeVal, ok := hashGet(expr, "type")
	if !ok {
		t.Error("resultToExpr(Ok) should have 'type' field")
	}
	if typeVal.Type != String || typeVal.Str != "ok" {
		t.Errorf("type field = %v, want 'ok'", typeVal.Str)
	}

	// Verify value field contains the original value
	valueVal, ok := hashGet(expr, "value")
	if !ok {
		t.Error("resultToExpr(Ok) should have 'value' field")
	}
	if valueVal.Type != String || valueVal.Str != "test value" {
		t.Errorf("value field = %v, want 'test value'", valueVal.Str)
	}

	// Verify no error field
	_, hasError := hashGet(expr, "error")
	if hasError {
		t.Error("resultToExpr(Ok) should not have 'error' field")
	}
}

func TestResultToExprOkWithNumber(t *testing.T) {
	// Create an Ok result with a number value
	value := makeNum(42)
	result := Ok(value)

	// Convert to Expr
	expr := resultToExpr(result)

	// Verify value field contains the number
	valueVal, ok := hashGet(expr, "value")
	if !ok {
		t.Error("resultToExpr(Ok) should have 'value' field")
	}
	if valueVal.Type != Number || valueVal.Num != 42 {
		t.Errorf("value field = %v, want 42", valueVal.Num)
	}
}

func TestResultToExprOkWithHash(t *testing.T) {
	// Create an Ok result with a hash value
	hash := makeHash()
	hashSet(hash, "name", makeStr("Alice"))
	result := Ok(hash)

	// Convert to Expr
	expr := resultToExpr(result)

	// Verify value field contains the hash
	valueVal, ok := hashGet(expr, "value")
	if !ok {
		t.Error("resultToExpr(Ok) should have 'value' field")
	}
	if valueVal.Type != Hash {
		t.Errorf("value field type = %v, want Hash", valueVal.Type)
	}

	// Verify the hash contents
	name, ok := hashGet(valueVal, "name")
	if !ok || name.Str != "Alice" {
		t.Error("value hash should contain 'name' = 'Alice'")
	}
}

func TestResultToExprErr(t *testing.T) {
	// Create an Err result
	result := Err[*Expr](errors.New("test error"))

	// Convert to Expr
	expr := resultToExpr(result)

	// Verify it's a hash
	if expr.Type != Hash {
		t.Errorf("resultToExpr(Err) type = %v, want Hash", expr.Type)
	}

	// Verify type field is "err"
	typeVal, ok := hashGet(expr, "type")
	if !ok {
		t.Error("resultToExpr(Err) should have 'type' field")
	}
	if typeVal.Type != String || typeVal.Str != "err" {
		t.Errorf("type field = %v, want 'err'", typeVal.Str)
	}

	// Verify error field contains the error message
	errorVal, ok := hashGet(expr, "error")
	if !ok {
		t.Error("resultToExpr(Err) should have 'error' field")
	}
	if errorVal.Type != String || errorVal.Str != "test error" {
		t.Errorf("error field = %v, want 'test error'", errorVal.Str)
	}

	// Verify no value field
	_, hasValue := hashGet(expr, "value")
	if hasValue {
		t.Error("resultToExpr(Err) should not have 'value' field")
	}
}

func TestResultToExprErrWithComplexError(t *testing.T) {
	// Create an Err result with a more complex error message
	result := Err[*Expr](errors.New("HTTP 404: Not Found"))

	// Convert to Expr
	expr := resultToExpr(result)

	// Verify error field contains the full error message
	errorVal, ok := hashGet(expr, "error")
	if !ok {
		t.Error("resultToExpr(Err) should have 'error' field")
	}
	if errorVal.Str != "HTTP 404: Not Found" {
		t.Errorf("error field = %v, want 'HTTP 404: Not Found'", errorVal.Str)
	}
}
