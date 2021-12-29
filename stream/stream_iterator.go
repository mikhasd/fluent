package stream

import (
	"github.com/mikhasd/fluent"
	"github.com/mikhasd/fluent/iterator"
)

type iteratorStream[T any] struct {
	iterator iterator.Iterator[T]
}

func (s iteratorStream[T]) next() fluent.Option[T] {
	return s.iterator.Next()
}

// Skip

type skip[T any] struct {
	count   int
	skipped bool
	source  iterator.Iterator[T]
}

func (s skip[T]) Size() fluent.Option[int] {
	return iterator.Size(s.source).Map(func(size int) int {
		c := size - s.count
		if c < 0 {
			return 0
		}
		return c
	})
}

func (s *skip[T]) Next() fluent.Option[T] {
	if !s.skipped {
		for i := 0; i < s.count; i++ {
			o := s.source.Next()
			if !o.Present() {
				s.skipped = true
				return o
			}
		}
		s.skipped = true
	}
	return s.source.Next()
}

func (s *iteratorStream[T]) Skip(count int) Stream[T] {
	return &iteratorStream[T]{
		iterator: &skip[T]{
			count:   count,
			skipped: false,
			source:  s.iterator,
		},
	}
}

// Limit

type limit[T any] struct {
	max     int
	current int
	source  iterator.Iterator[T]
}

func (l limit[T]) Size() fluent.Option[int] {
	return iterator.Size(l.source).Map(func(size int) int {
		if size < l.max {
			return size
		}
		return l.max
	})
}

func (l *limit[T]) Next() fluent.Option[T] {
	if l.current < l.max {
		l.current++
		return l.source.Next()
	} else {
		return fluent.Empty[T]()
	}
}

func (s *iteratorStream[T]) Limit(max int) Stream[T] {
	return &iteratorStream[T]{
		iterator: &limit[T]{
			max:     max,
			current: 0,
			source:  s.iterator,
		},
	}
}

// Filter

type filter[T any] struct {
	filter func(T) bool
	source iterator.Iterator[T]
}

func (f filter[T]) Next() fluent.Option[T] {
	o := f.source.Next()
	for o.Present() && !f.filter(o.Get()) {
		o = f.source.Next()
	}
	return o
}

func (s *iteratorStream[T]) Filter(fn func(T) bool) Stream[T] {
	return &iteratorStream[T]{
		iterator: filter[T]{
			filter: fn,
			source: s.iterator,
		},
	}
}

// Map

type mapper[T any] struct {
	mapper func(T) T
	source iterator.Iterator[T]
}

func (m mapper[T]) Size() fluent.Option[int] {
	return iterator.Size(m.source)
}

func (m mapper[T]) Next() fluent.Option[T] {
	return m.source.Next().Map(m.mapper)
}

func (s *iteratorStream[T]) Map(fn func(T) T) Stream[T] {
	return &iteratorStream[T]{
		iterator: mapper[T]{
			mapper: fn,
			source: s.iterator,
		},
	}
}

// Peek

type peek[T any] struct {
	consumer func(T)
	source   iterator.Iterator[T]
}

func (p peek[T]) Size() fluent.Option[int] {
	return iterator.Size(p.source)
}

func (p peek[T]) Next() fluent.Option[T] {
	o := p.source.Next()
	o.IfPresent(p.consumer)
	return o
}

func (s *iteratorStream[T]) Peek(consumer func(T)) Stream[T] {
	return &iteratorStream[T]{
		iterator: peek[T]{
			consumer: consumer,
			source:   s.iterator,
		},
	}
}

// For Each

func (s *iteratorStream[T]) ForEach(fn func(T)) {
	it := s.iterator
	for o := it.Next(); o.Present(); o = it.Next() {
		fn(o.Get())
	}
}

// Count

func (s *iteratorStream[T]) Count() int {
	var counter int
	it := s.iterator
	for o := it.Next(); o.Present(); o = it.Next() {
		counter++
	}
	return counter
}

// Iterator

func (s *iteratorStream[T]) Iterator() iterator.Iterator[T] {
	return s.iterator
}

// Array
func (s *iteratorStream[T]) Array() []T {
	knownSize := func(size int) []T {
		arr := make([]T, size)

		for i := 0; i < size; i++ {
			arr[i] = s.iterator.Next().Get()
		}

		return arr
	}

	unknownSize := func() []T {
		arr := make([]T, 0, 10)
		for o := s.iterator.Next(); o.Present(); o = s.iterator.Next() {
			arr = append(arr, o.Get())
		}
		return arr
	}

	return fluent.MapOption(
		iterator.Size(s.iterator),
		knownSize,
	).OrElseGet(unknownSize)
}
