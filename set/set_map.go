package set

import "github.com/mikhasd/fluent/iterator"

type mapSet[T comparable] map[T]setEntry

func (s mapSet[T]) Contains(element T) bool {
	_, found := s[element]
	return found
}

func (s mapSet[T]) ContainsAll(iter iterator.Iterable[T]) bool {
	it := iter.Iterator()
	for o := it.Next(); o.IsPresent(); o = it.Next() {
		el := o.Get()
		if _, found := s[el]; !found {
			return false
		}
	}
	return true
}

func (s mapSet[T]) Add(element T) {
	s[element] = setEntry{}
}

func (s mapSet[T]) AddAll(iter iterator.Iterable[T]) {
	it := iter.Iterator()
	for o := it.Next(); o.IsPresent(); o = it.Next() {
		el := o.Get()
		s[el] = setEntry{}
	}
}

func (s mapSet[T]) Iterator() iterator.Iterator[T] {
	return iterator.MapKeys(s)
}

func (s mapSet[T]) ForEach(fn func(T)) {
	for el := range s {
		fn(el)
	}
}

func (s mapSet[T]) Remove(el T) {
	delete(s, el)
}

func (s mapSet[T]) Empty() bool {
	return len(s) == 0
}

func (s mapSet[T]) Size() int {
	return len(s)
}
