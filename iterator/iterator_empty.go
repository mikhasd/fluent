package iterator

import "github.com/mikhasd/fluent"

type emptyIterator[T any] struct{}

func (i *emptyIterator[T]) Next() fluent.Option[T] {
	return fluent.Empty[T]()
}

// Implements iterator.Sized interface
func (i emptyIterator[T]) Size() int {
	return 0
}
