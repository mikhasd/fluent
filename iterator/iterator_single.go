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

func (i singleItemIterator[T]) Size() fluent.Option[int] {
	return fluent.Present(1)
}
