package stream

import (
	"github.com/mikhasd/fluent"
	"github.com/mikhasd/fluent/iterator"
)

// Stream is a lazy sequence of elements of a generic type that can be
// manipulated thru a pipeline of operations.
//
// A stream pipeline can be used only once.
type Stream[T any] interface {
	// Skip discards the first `n` elements of the stream.
	Skip(n int) Stream[T]

	// Limit truncates the stream to be no larger than `max` elements.
	Limit(max int) Stream[T]

	// Filter returns a stream consisting of only the elements that matches the
	// provided `condition`.
	Filter(condition func(T) bool) Stream[T]

	// Map returns a stream with the results of applying the provided `mapper`
	// function to each item of the original stream.
	//
	// Due to go language limitations, the return type of the mapper function
	// must be the same type as the input.
	//
	// If different input and output types are required, use `stream.Map`.
	Map(mapper func(T) T) Stream[T]

	// ForEach executes the provided `consumer` function on each element of the
	// stream.
	ForEach(consumer func(int, T))

	// Peek executes the provided `consumer` function with each element of the
	// stream and returns a stream with the same elements.
	Peek(consumer func(T)) Stream[T]

	// Count returns the number of elements in the stream.
	Count() int

	// Iterator returns an Iterator with the result of the stream pipeline
	// processing.
	Iterator() iterator.Iterator[T]

	// Array collects into an array the result of the stream pipeline processing.
	Array() []T

	While(func(T) bool) Stream[T]

	Parallel() Stream[T]
}

// FromArray creates a stream with an array as data source.
func FromArray[T any](arr []T) Stream[T] {
	it := iterator.FromArray(arr)
	return FromIterator(it)
}

// Of creates a stream with the given elements.
func Of[T any](items ...T) Stream[T] {
	return FromArray(items)
}

// FromIterable creates a stream from an iterator.Iterable.
func FromIterable[T any](iter iterator.Iterable[T]) Stream[T] {
	return FromIterator(iter.Iterator())
}

// FromIterator cerates a stream from an iterator.Iterator.
func FromIterator[T any](it iterator.Iterator[T]) Stream[T] {
	return &iteratorStream[T]{
		iterator: it,
	}
}

// Map applies the `mapper` function to the elements of the source stream and
// returns a new stream with the results.
//
// This function is useful if the input and output types of the mapper function
// are different.
func Map[A any, R any](s Stream[A], mapper func(A) R) Stream[R] {
	it := s.Iterator()
	return FromIterator(iterator.Func(func() fluent.Option[R] {
		return fluent.MapOption(it.Next(), mapper)
	}))
}

// MapArray is a shortcut to create a stream with the input array and apply
// stream.Map to the stream.
//
//  stream.MapArray(inputArray, myMapperFunction)
//
// is equivalent to:
//
//  stream.Map(stream.FromArray(inputArray), myMapperFunction)
func MapArray[I any, O any](in []I, mapper func(I) O) Stream[O] {
	return Map(FromArray(in), mapper)
}
