package minilisp

import "testing"

func TestNewEnv(t *testing.T) {
	env := NewEnv(nil)
	if env == nil {
		t.Fatal("NewEnv returned nil")
	}
	if env.bindings == nil {
		t.Fatal("NewEnv bindings map is nil")
	}
}

func TestEnvDefineAndLookup(t *testing.T) {
	env := NewEnv(nil)
	env.Define("x", makeNum(42))

	val, ok := env.Lookup("x")
	if !ok {
		t.Fatal("Lookup failed for defined variable")
	}
	if val.Type != Number || val.Num != 42 {
		t.Errorf("Lookup returned %v, want 42", val)
	}
}

func TestEnvLookupUndefined(t *testing.T) {
	env := NewEnv(nil)
	_, ok := env.Lookup("undefined")
	if ok {
		t.Error("Lookup should return false for undefined variable")
	}
}

func TestEnvParentScoping(t *testing.T) {
	// Create parent environment
	parent := NewEnv(nil)
	parent.Define("x", makeNum(10))
	parent.Define("y", makeNum(20))

	// Create child environment
	child := NewEnv(parent)
	child.Define("x", makeNum(99)) // Shadow parent's x

	// Child should see its own x
	val, ok := child.Lookup("x")
	if !ok || val.Num != 99 {
		t.Errorf("child x = %v, want 99", val)
	}

	// Child should see parent's y
	val, ok = child.Lookup("y")
	if !ok || val.Num != 20 {
		t.Errorf("child y = %v, want 20", val)
	}

	// Parent should still see original x
	val, ok = parent.Lookup("x")
	if !ok || val.Num != 10 {
		t.Errorf("parent x = %v, want 10", val)
	}
}

func TestEnvMultipleLevels(t *testing.T) {
	// Grandparent -> Parent -> Child
	grandparent := NewEnv(nil)
	grandparent.Define("a", makeNum(1))

	parent := NewEnv(grandparent)
	parent.Define("b", makeNum(2))

	child := NewEnv(parent)
	child.Define("c", makeNum(3))

	// Child should see all three
	if val, ok := child.Lookup("a"); !ok || val.Num != 1 {
		t.Error("child can't see grandparent variable")
	}
	if val, ok := child.Lookup("b"); !ok || val.Num != 2 {
		t.Error("child can't see parent variable")
	}
	if val, ok := child.Lookup("c"); !ok || val.Num != 3 {
		t.Error("child can't see own variable")
	}

	// Parent should NOT see child's variable
	if _, ok := parent.Lookup("c"); ok {
		t.Error("parent shouldn't see child variable")
	}
}

func TestEnvRedefineSameScope(t *testing.T) {
	env := NewEnv(nil)
	env.Define("x", makeNum(10))
	env.Define("x", makeNum(20)) // Redefine in same scope

	val, ok := env.Lookup("x")
	if !ok || val.Num != 20 {
		t.Errorf("x = %v, want 20", val)
	}
}
