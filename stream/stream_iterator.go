package stream

import (
	"sync"
	"sync/atomic"

	"github.com/mikhasd/fluent"
	"github.com/mikhasd/fluent/iterator"
)

type iteratorStream[T any] struct {
	iterator iterator.Iterator[T]
	parallel bool
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
		parallel: s.parallel,
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
		parallel: s.parallel,
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
		parallel: s.parallel,
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
		parallel: s.parallel,
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
		parallel: s.parallel,
		iterator: peek[T]{
			consumer: consumer,
			source:   s.iterator,
		},
	}
}

type concurrent[T any] struct {
	lock   sync.Mutex
	source iterator.Iterator[T]
}

func (c *concurrent[T]) Next() fluent.Option[T] {
	c.lock.Lock()
	defer c.lock.Unlock()
	o := c.source.Next()
	return o
}

func (s *iteratorStream[T]) Parallel() Stream[T] {
	return &iteratorStream[T]{
		parallel: true,
		iterator: &concurrent[T]{
			source: s.iterator,
		},
	}
}

// For Each

func (s *iteratorStream[T]) ForEach(fn func(T)) {
	it := s.iterator
	if s.parallel {
		s.parallelForEach(fn)
	} else {
		for o := it.Next(); o.Present(); o = it.Next() {
			fn(o.Get())
		}
	}

}

func (s *iteratorStream[T]) parallelForEach(fn func(T)) {
	it := s.iterator
	var iteratorDone int32 = 0
	var wg sync.WaitGroup

	for atomic.LoadInt32(&iteratorDone) == 0 {
		go func(done *int32) {
			wg.Add(1)
			if o := it.Next(); o.Present() {
				fn(o.Get())
			} else {
				atomic.StoreInt32(done, 1)
			}
			wg.Done()
		}(&iteratorDone)
	}
	wg.Wait()

}

// Count

func (s *iteratorStream[T]) Count() int {
	var counter int32
	s.ForEach(func(_ T) {
		atomic.AddInt32(&counter, 1)
	})
	return int(counter)
}

// Iterator

func (s *iteratorStream[T]) Iterator() iterator.Iterator[T] {
	return s.iterator
}

// Array
func (s *iteratorStream[T]) Array() []T {
	if s.parallel {
		return s.toArrayParallel()
	} else {
		return s.toArray()
	}
}

func (s *iteratorStream[T]) toArray() []T {
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

func (s *iteratorStream[T]) toArrayParallel() []T {
	knownSize := func(size int) []T {
		arr := make([]T, size)

		for i := 0; i < size; i++ {
			go func(index int) {
				arr[index] = s.iterator.Next().Get()
			}(i)
		}

		return arr
	}

	unknownSize := func() []T {
		var wq sync.WaitGroup
		var arr []T

		done := new(bool)
		*done = false

		elements := make(chan T)

		for !*done {
			go func() {
				if o := s.iterator.Next(); o.Present() {
					wq.Add(1)
					elem := o.Get()
					elements <- elem
				} else {
					*done = true
				}
			}()
		}

		go func() {
			for elem := range elements {
				arr = append(arr, elem)
				wq.Done()
			}
		}()

		wq.Wait()
		return arr
	}

	return fluent.MapOption(
		iterator.Size(s.iterator),
		knownSize,
	).OrElseGet(unknownSize)
}
