package set

import (
	"github.com/mikhasd/fluent/iterator"
)

type mapSet[K comparable, V any] struct {
	items  map[K]V
	hasher func(V) K
}

func (s mapSet[K, V]) Contains(element V) bool {
	_, found := s.items[s.hasher(element)]
	return found
}

func (s mapSet[K, V]) ContainsAll(iter iterator.Iterable[V]) bool {
	it := iter.Iterator()
	for o := it.Next(); o.IsPresent(); o = it.Next() {
		el := o.Get()
		if _, found := s.items[s.hasher(el)]; !found {
			return false
		}
	}
	return true
}

func (s mapSet[K, V]) Add(element V) {
	s.items[s.hasher(element)] = element
}

func (s mapSet[K, V]) AddAll(iter iterator.Iterable[V]) {
	it := iter.Iterator()
	for o := it.Next(); o.IsPresent(); o = it.Next() {
		el := o.Get()
		hash := s.hasher(el)
		s.items[hash] = el
	}
}

func (s mapSet[K, V]) Iterator() iterator.Iterator[V] {
	return iterator.MapValues(s.items)
}

func (s mapSet[K, V]) ForEach(fn func(V)) {
	for _, v := range s.items {
		fn(v)
	}
}

func (s mapSet[K, V]) Remove(el V) {
	delete(s.items, s.hasher(el))
}

func (s mapSet[K, V]) Empty() bool {
	return len(s.items) == 0
}

func (s mapSet[K, V]) Size() int {
	return len(s.items)
}

func New[T comparable]() Set[T] {
	return WithSize[T](16)
}

func WithSize[T comparable](size int) Set[T] {
	return WithSizeAndHasher(size, func(t T) T {
		return t
	})
}

func WithSizeAndHasher[K comparable, V any](size int, hasher func(V) K) Set[V] {
	return &mapSet[K, V]{
		items:  make(map[K]V, size),
		hasher: hasher,
	}
}

func FromArray[T comparable](arr []T) Set[T] {
	set := WithSize[T](len(arr))
	set.AddAll(iterator.ArrayIterable(arr))
	return set
}

func FromIterable[T comparable](iter iterator.Iterable[T]) Set[T] {
	s := New[T]()
	s.AddAll(iter)
	return s
}
