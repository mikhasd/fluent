package fluent

import "fmt"

type empty[T any] struct{}

func (e empty[T]) IsPresent() bool {
	return false
}

func (e empty[T]) Get() T {
	panic("empty option")
}

func (e empty[T]) Map(func(T) T) Option[T] {
	return e
}

func (e empty[T]) OrElse(other T) T {
	return other
}

func (e empty[T]) OrElseGet(supplier func() T) T {
	return supplier()
}

func (e empty[T]) Or(supplier func() Option[T]) Option[T] {
	return supplier()
}

func (e empty[T]) OrError(err error) Result[T] {
	return Err[T](err)
}

func (e empty[T]) IfPresent(func(T)) {
	// Do nothing
}

func (e empty[T]) Filter(func(T) bool) Option[T] {
	return e
}

func (e empty[T]) String() string {
	return fmt.Sprintf("Empty[]")
}
