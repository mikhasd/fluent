package fluent

import "fmt"

type err[T any] struct {
	e error
}

func (e err[T]) IsOk() bool {
	return false
}

func (e err[T]) IsErr() bool {
	return true
}

func (e err[T]) Ok() Option[T] {
	return Empty[T]()
}

func (e err[T]) Err() Option[error] {
	return Present(e.e)
}

func (e err[T]) Map(mapper func(T) T) Result[T] {
	return e
}

func (e err[T]) MapErr(mapper func(error) T) Result[T] {
	return ok[T]{mapper(e.e)}
}

func (e err[T]) Get() T {
	panic("error result")
}

func (e err[T]) GetErr() error {
	return e.e
}

func (e err[T]) OrElse(other T) T {
	return other
}

func (e err[T]) OrElseGet(supplier func() T) T {
	return supplier()
}

func (e err[T]) Or(supplier func() Result[T]) Result[T] {
	return supplier()
}

func (e err[T]) String() string {
	return fmt.Sprintf("Err[%+v]", e.e)
}
