package iterator

import "github.com/mikhasd/fluent"

func MapKeys[K comparable, V any](m map[K]V) Iterator[K] {
	keys := make([]K, len(m))
	index := 0
	for k := range m {
		keys[index] = k
		index++
	}
	return FromArray(keys)
}

func MapValues[K comparable, V any](m map[K]V) Iterator[V] {
	values := make([]V, len(m))
	index := 0
	for _, v := range m {
		values[index] = v
		index++
	}
	return FromArray(values)
}

type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

type mapIterator[K comparable, V any] struct {
	data  map[K]V
	keys  Iterator[K]
	index int
}

func FromMap[K comparable, V any](m map[K]V) Iterator[MapEntry[K, V]] {
	return mapIterator[K, V]{
		data:  m,
		keys:  MapKeys(m),
		index: 0,
	}
}

func (it mapIterator[K, V]) Next() fluent.Option[MapEntry[K, V]] {
	return fluent.MapOption(it.keys.Next(), func(key K) MapEntry[K, V] {
		return MapEntry[K, V]{
			Key:   key,
			Value: it.data[key],
		}
	})
}
