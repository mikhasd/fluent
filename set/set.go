package set

import "github.com/mikhasd/fluent/iterator"

type Set[T comparable] interface {
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

type setEntry struct{}

func New[T comparable]() Set[T] {
	return make(mapSet[T], 16)
}

func WithSize[T comparable](size uint) Set[T] {
	return make(mapSet[T], size)
}

func FromArray[T comparable](arr []T) Set[T] {
	s := make(mapSet[T], len(arr))
	for _, el := range arr {
		s[el] = setEntry{}
	}
	return s
}

func FromIterable[T comparable](iter iterator.Iterable[T]) Set[T] {
	s := New[T]()
	s.AddAll(iter)
	return s
}
