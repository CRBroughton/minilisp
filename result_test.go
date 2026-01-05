package main

import (
	"errors"
	"testing"
)

func TestOk(t *testing.T) {
	r := Ok(42)

	if !r.IsOk() {
		t.Error("Ok result should be ok")
	}
	if r.IsErr() {
		t.Error("Ok result should not be err")
	}
	if r.Unwrap() != 42 {
		t.Errorf("Unwrap() = %d, want 42", r.Unwrap())
	}
}

func TestErr(t *testing.T) {
	r := Err[int](errors.New("test error"))

	if r.IsOk() {
		t.Error("Err result should not be ok")
	}
	if !r.IsErr() {
		t.Error("Err result should be err")
	}
	if r.Error().Error() != "test error" {
		t.Errorf("Error() = %v, want 'test error'", r.Error())
	}
}

func TestUnwrapPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Unwrap on Err should panic")
		}
	}()

	r := Err[int](errors.New("test"))
	r.Unwrap()
}

func TestUnwrapOr(t *testing.T) {
	okResult := Ok(42)
	if okResult.UnwrapOr(99) != 42 {
		t.Error("UnwrapOr on Ok should return value")
	}

	errResult := Err[int](errors.New("test"))
	if errResult.UnwrapOr(99) != 99 {
		t.Error("UnwrapOr on Err should return default")
	}
}

func TestUnwrapOrElse(t *testing.T) {
	okResult := Ok(42)
	val := okResult.UnwrapOrElse(func(err error) int {
		return 99
	})
	if val != 42 {
		t.Error("UnwrapOrElse on Ok should return value")
	}

	errResult := Err[int](errors.New("test"))
	val = errResult.UnwrapOrElse(func(err error) int {
		return 99
	})
	if val != 99 {
		t.Error("UnwrapOrElse on Err should call function")
	}
}

func TestMap(t *testing.T) {
	okResult := Ok(21)
	mapped := okResult.Map(func(x int) int { return x * 2 })

	if !mapped.IsOk() {
		t.Error("Map on Ok should return Ok")
	}
	if mapped.Unwrap() != 42 {
		t.Errorf("Mapped value = %d, want 42", mapped.Unwrap())
	}

	errResult := Err[int](errors.New("test"))
	mapped = errResult.Map(func(x int) int { return x * 2 })

	if !mapped.IsErr() {
		t.Error("Map on Err should return Err")
	}
}

func TestMapErr(t *testing.T) {
	okResult := Ok(42)
	mapped := okResult.MapErr(func(err error) error {
		return errors.New("modified")
	})

	if !mapped.IsOk() {
		t.Error("MapErr on Ok should return Ok")
	}

	errResult := Err[int](errors.New("original"))
	mapped = errResult.MapErr(func(err error) error {
		return errors.New("modified")
	})

	if !mapped.IsErr() {
		t.Error("MapErr on Err should return Err")
	}
	if mapped.Error().Error() != "modified" {
		t.Error("MapErr should transform error")
	}
}

func TestAndThen(t *testing.T) {
	okResult := Ok(21)
	chained := okResult.AndThen(func(x int) Result[int] {
		return Ok(x * 2)
	})

	if !chained.IsOk() {
		t.Error("AndThen on Ok should return Ok")
	}
	if chained.Unwrap() != 42 {
		t.Errorf("Chained value = %d, want 42", chained.Unwrap())
	}

	errResult := Err[int](errors.New("test"))
	chained = errResult.AndThen(func(x int) Result[int] {
		return Ok(x * 2)
	})

	if !chained.IsErr() {
		t.Error("AndThen on Err should return Err")
	}
}

func TestOrElse(t *testing.T) {
	okResult := Ok(42)
	recovered := okResult.OrElse(func(err error) Result[int] {
		return Ok(99)
	})

	if recovered.Unwrap() != 42 {
		t.Error("OrElse on Ok should return original Ok")
	}

	errResult := Err[int](errors.New("test"))
	recovered = errResult.OrElse(func(err error) Result[int] {
		return Ok(99)
	})

	if !recovered.IsOk() {
		t.Error("OrElse on Err should call function")
	}
	if recovered.Unwrap() != 99 {
		t.Error("OrElse should recover with new value")
	}
}

func TestMatch(t *testing.T) {
	okCalled := false
	errCalled := false

	okResult := Ok(42)
	okResult.Match(
		func(x int) { okCalled = true },
		func(err error) { errCalled = true },
	)

	if !okCalled || errCalled {
		t.Error("Match on Ok should call ok function")
	}

	okCalled = false
	errCalled = false

	errResult := Err[int](errors.New("test"))
	errResult.Match(
		func(x int) { okCalled = true },
		func(err error) { errCalled = true },
	)

	if okCalled || !errCalled {
		t.Error("Match on Err should call err function")
	}
}

func TestToTuple(t *testing.T) {
	okResult := Ok(42)
	val, err := okResult.ToTuple()

	if err != nil {
		t.Error("Ok ToTuple should have nil error")
	}
	if val != 42 {
		t.Errorf("Ok ToTuple value = %d, want 42", val)
	}

	errResult := Err[int](errors.New("test"))
	val, err = errResult.ToTuple()

	if err == nil {
		t.Error("Err ToTuple should have error")
	}
	if val != 0 {
		t.Error("Err ToTuple should have zero value")
	}
}

func TestFromTuple(t *testing.T) {
	r := FromTuple(42, nil)
	if !r.IsOk() || r.Unwrap() != 42 {
		t.Error("FromTuple with nil error should create Ok")
	}

	r = FromTuple(0, errors.New("test"))
	if !r.IsErr() {
		t.Error("FromTuple with error should create Err")
	}
}
