package iterator

import "github.com/mikhasd/fluent"

type arrayIterator[T any] struct {
	data  []T
	index int
}

type arrayIterable[T any] []T

func ArrayIterable[T any](arr []T) Iterable[T] {
	return arrayIterable[T](arr)
}

func (it arrayIterable[T]) Iterator() Iterator[T] {
	return FromArray(it)
}

func (it *arrayIterator[T]) Next() fluent.Option[T] {
	if it.index < len(it.data) {
		value := it.data[it.index]
		it.index = it.index + 1
		return fluent.Present(value)
	} else {
		return fluent.Empty[T]()
	}
}

// Implements iterator.sized interface
func (it *arrayIterator[T]) Size() int {
	return len(it.data)
}
