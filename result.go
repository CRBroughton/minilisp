package main

import "fmt"

type Result[T any] struct {
	value *T
	err   error
	isOk  bool
}

func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: &value,
		err:   nil,
		isOk:  true,
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		value: nil,
		err:   err,
		isOk:  false,
	}
}

func (r Result[T]) IsOk() bool {
	return r.isOk
}

func (r Result[T]) IsErr() bool {
	return !r.isOk
}

func (r Result[T]) Unwrap() T {
	if !r.isOk {
		panic(fmt.Sprintf("called Unwrap on an Err value: %v", r.err))
	}
	return *r.value
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.isOk {
		return *r.value
	}
	return defaultValue
}

func (r Result[T]) UnwrapOrElse(fn func(error) T) T {
	if r.isOk {
		return *r.value
	}
	return fn(r.err)
}

func (r Result[T]) Error() error {
	return r.err
}

func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.isOk {
		return Ok(fn(*r.value))
	}
	return Err[T](r.err)
}

func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	if r.isErr() {
		return Err[T](fn(r.err))
	}
	return r
}

func (r Result[T]) AndThen(fn func(T) Result[T]) Result[T] {
	if r.isOk {
		return fn(*r.value)
	}
	return Err[T](r.err)
}

func (r Result[T]) OrElse(fn func(error) Result[T]) Result[T] {
	if r.isOk {
		return r
	}
	return fn(r.err)
}

func (r Result[T]) Match(okFn func(T), errFn func(error)) {
	if r.isOk {
		okFn(*r.value)
	} else {
		errFn(r.err)
	}
}

func (r Result[T]) ToTuple() (T, error) {
	if r.isOk {
		return *r.value, nil
	}
	var zero T
	return zero, r.err
}

func FromTuple[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

func (r Result[T]) isErr() bool {
	return !r.isOk
}
