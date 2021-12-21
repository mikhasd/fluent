package fluent

import "fmt"

type present[T any] struct {
	value T
}

func (p present[T]) Present() bool {
	return true
}

func (p present[T]) Get() T {
	return p.value
}

func (p present[T]) Map(mapper func(T) T) Option[T] {
	return present[T]{mapper(p.value)}
}

func (p present[T]) OrElse(other T) T {
	return p.value
}

func (p present[T]) OrElseGet(supplier func() T) T {
	return p.value
}

func (p present[T]) Or(supplier func() Option[T]) Option[T] {
	return p
}

func (p present[T]) OrError(e error) Result[T] {
	return ResultOk(p.value)
}

func (p present[T]) String() string {
	return fmt.Sprintf("Present[%+v]", p.value)
}
