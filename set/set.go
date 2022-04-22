package set

import "github.com/mikhasd/fluent/iterator"

type Set[T any] interface {
	Contains(T) bool
	ContainsAll(iterator.Iterable[T]) bool
	Add(T)
	AddAll(iterator.Iterable[T])
	Iterator() iterator.Iterator[T]
	Remove(T)
	Empty() bool
	Size() int
	ForEach(fn func(T))
}
