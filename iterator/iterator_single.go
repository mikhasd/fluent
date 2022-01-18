package iterator

import "github.com/mikhasd/fluent"

type singleItemIterator[T any] struct {
	consumed bool
	item     T
}

func (i *singleItemIterator[T]) Next() fluent.Option[T] {
	if i.consumed {
		return fluent.Empty[T]()
	} else {
		i.consumed = true
		return fluent.Present(i.item)
	}
}

// Implements iterator.Sized interface
func (i singleItemIterator[T]) Size() int {
	return 1
}
