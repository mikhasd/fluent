package stream

import (
	"github.com/mikhasd/fluent"
	"github.com/mikhasd/fluent/iterator"
)

type Stream[T any] interface {
	Skip(count int) Stream[T]
	Limit(max int) Stream[T]
	Filter(mapper func(T) bool) Stream[T]
	Map(mapper func(T) T) Stream[T]
	ForEach(consumer func(T))
	Peek(consumer func(T)) Stream[T]
	Count() int
	Iterator() iterator.Iterator[T]
	Array() []T
}

func FromArray[T any](arr []T) Stream[T] {
	it := iterator.FromArray(arr)
	return FromIterator(it)
}

func Of[T any](items ...T) Stream[T] {
	return FromArray(items)
}

func FromIterable[T any](iter iterator.Iterable[T]) Stream[T] {
	return FromIterator(iter.Iterator())
}

func FromIterator[T any](it iterator.Iterator[T]) Stream[T] {
	return &iteratorStream[T]{
		iterator: it,
	}
}

func Map[A any, R any](s Stream[A], mapper func(A) R) Stream[R] {
	it := s.Iterator()
	return FromIterator(iterator.Func(func() fluent.Option[R] {
		return fluent.MapOption(it.Next(), mapper)
	}))
}
