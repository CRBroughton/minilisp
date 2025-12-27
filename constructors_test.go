package main

import "testing"

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
