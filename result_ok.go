package fluent

import "fmt"

type ok[T any] struct {
	value T
}

func (o ok[T]) IsOk() bool {
	return true
}

func (o ok[T]) IsErr() bool {
	return false
}

func (o ok[T]) Ok() Option[T] {
	return Present(o.value)
}

func (o ok[T]) Err() Option[error] {
	return Empty[error]()
}

func (o ok[T]) Map(mapper func(T) T) Result[T] {
	return Ok(mapper(o.value))
}

func (o ok[T]) MapErr(mapper func(error) T) Result[T] {
	return o
}

func (o ok[T]) Get() T {
	return o.value
}

func (o ok[T]) GetErr() error {
	panic("ok result")
}

func (o ok[T]) OrElse(other T) T {
	return o.value
}

func (o ok[T]) OrElseGet(supplier func() T) T {
	return o.value
}

func (o ok[T]) Or(supplier func() Result[T]) Result[T] {
	return o
}

func (o ok[T]) String() string {
	return fmt.Sprintf("Ok[%+v]", o.value)
}
