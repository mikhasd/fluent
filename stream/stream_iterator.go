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
			if !o.IsPresent() {
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

// While

type while[T any] struct {
	condition func(T) bool
	source    iterator.Iterator[T]
}

func (w while[T]) Size() fluent.Option[int] {
	return iterator.Size(w.source)
}

func (w *while[T]) Next() fluent.Option[T] {
	next := w.source.Next()
	if next.IsPresent() && w.condition(next.Get()) {
		return next
	}
	return fluent.Empty[T]()
}

func (s *iteratorStream[T]) While(condition func(T) bool) Stream[T] {
	return &iteratorStream[T]{
		parallel: s.parallel,
		iterator: &while[T]{
			condition: condition,
			source:    s.iterator,
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
	for o.IsPresent() && !f.filter(o.Get()) {
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

func (c *concurrent[T]) Size() fluent.Option[int] {
	return iterator.Size(c.source)
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

func (s *iteratorStream[T]) ForEach(fn func(int, T)) {
	it := s.iterator
	if s.parallel {
		size := iterator.Size(s.iterator)
		if size.IsPresent() {
			s.sizedParallelForEach(fn, size.Get())
		} else {
			s.parallelForEach(fn)
		}
	} else {
		var index int = 0
		for o := it.Next(); o.IsPresent(); o = it.Next() {
			fn(index, o.Get())
			index++
		}
	}

}

func (s *iteratorStream[T]) parallelForEach(fn func(int, T)) {
	it := s.iterator
	var done int32 = 0
	var index int32 = 0
	var wg sync.WaitGroup

	for atomic.LoadInt32(&done) == 0 {
		wg.Add(1)
		go func() {
			if o := it.Next(); o.IsPresent() {
				i := atomic.LoadInt32(&index)
				fn(int(i), o.Get())
				atomic.AddInt32(&index, 1)
			} else {
				atomic.CompareAndSwapInt32(&done, 0, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func (s *iteratorStream[T]) sizedParallelForEach(fn func(int, T), size int) {
	it := s.iterator
	var wg sync.WaitGroup

	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(index int) {
			if o := it.Next(); o.IsPresent() {
				fn(int(index), o.Get())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// Count

func (s *iteratorStream[T]) Count() int {
	var counter int32
	s.ForEach(func(_ int, _ T) {
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
	size := iterator.Size(s.iterator)

	var arr []T

	if size.IsPresent() {
		arr = make([]T, size.Get())
		s.ForEach(func(index int, val T) {
			arr[index] = val
		})
	} else if s.parallel {
		var mtx sync.Mutex
		arr = make([]T, 0, 10)
		s.ForEach(func(_ int, val T) {
			mtx.Lock()
			arr = append(arr, val)
			mtx.Unlock()
		})
	} else {
		arr = make([]T, 0, 10)
		s.ForEach(func(_ int, val T) {
			arr = append(arr, val)
		})
	}

	return arr
}
