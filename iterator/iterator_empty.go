package iterator

import "github.com/mikhasd/fluent"

type emptyIterator[T any] struct{}

func (i *emptyIterator[T]) Next() fluent.Option[T] {
	return fluent.Empty[T]()
}

func (i emptyIterator[T]) Size() fluent.Option[int] {
	return fluent.Present(0)
}
