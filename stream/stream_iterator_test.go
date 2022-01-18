package stream

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/mikhasd/fluent"
	"github.com/mikhasd/fluent/iterator"
	"github.com/mikhasd/fluent/set"
	"github.com/stretchr/testify/assert"
)

var streamTestData []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func Test_iteratorStream_next(t *testing.T) {
	stream := FromArray(streamTestData)
	it := stream.Iterator()

	var o fluent.Option[int]
	for i := range streamTestData {
		val := streamTestData[i]
		o = it.Next()

		assert.NotNil(t, o, "option")
		assert.True(t, o.IsPresent(), "present")
		assert.Equal(t, val, o.Get(), "value")
	}

	o = it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_iteratorStream_Skip(t *testing.T) {
	count := 5
	stream := FromArray(streamTestData).Skip(count)
	it := stream.Iterator()
	data := streamTestData[count:]

	var o fluent.Option[int]
	for i := range data {
		val := data[i]
		o = it.Next()

		assert.NotNil(t, o, "option")
		assert.True(t, o.IsPresent(), "present")
		assert.Equal(t, val, o.Get(), "value")
	}

	o = it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_iteratorStream_Skip_short(t *testing.T) {
	count := 11
	stream := FromArray(streamTestData).Skip(count)
	it := stream.Iterator()

	o := it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_skip_Size_larger(t *testing.T) {
	count := 5
	expected := 5
	actual := skip[int]{
		count:  count,
		source: iterator.FromArray(streamTestData),
	}.Size()

	assert.NotNil(t, actual, "option")
	assert.True(t, actual.IsPresent(), "present")
	assert.Equal(t, expected, actual.Get(), "size")
}

func Test_skip_Size_short(t *testing.T) {
	count := 11
	expected := 0
	actual := skip[int]{
		count:  count,
		source: iterator.FromArray(streamTestData),
	}.Size()

	assert.NotNil(t, actual, "option")
	assert.True(t, actual.IsPresent(), "present")
	assert.Equal(t, expected, actual.Get(), "size")
}

func Test_iteratorStream_Limit(t *testing.T) {
	count := 5
	stream := FromArray(streamTestData).Limit(count)
	data := streamTestData[0:count]
	it := stream.Iterator()

	var o fluent.Option[int]
	for i := range data {
		val := data[i]
		o = it.Next()

		assert.NotNil(t, o, "option")
		assert.True(t, o.IsPresent(), "present")
		assert.Equal(t, val, o.Get(), "value")
	}

	o = it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_limit_Size_larger(t *testing.T) {
	max := 5
	expected := 5
	actual := limit[int]{
		max:    max,
		source: iterator.FromArray(streamTestData),
	}.Size()

	assert.NotNil(t, actual, "option")
	assert.True(t, actual.IsPresent(), "present")
	assert.Equal(t, expected, actual.Get(), "size")
}

func Test_limit_Size_short(t *testing.T) {
	max := 12
	expected := len(streamTestData)
	actual := limit[int]{
		max:    max,
		source: iterator.FromArray(streamTestData),
	}.Size()

	assert.NotNil(t, actual, "option")
	assert.True(t, actual.IsPresent(), "present")
	assert.Equal(t, expected, actual.Get(), "size")
}

func Test_iteratorStream_Filter(t *testing.T) {
	isEven := func(n int) bool {
		return n%2 == 0
	}
	stream := FromArray(streamTestData).Filter(isEven)
	it := stream.Iterator()
	data := []int{2, 4, 6, 8, 10}

	var o fluent.Option[int]
	for i := range data {
		val := data[i]
		o = it.Next()

		assert.NotNil(t, o, "option")
		assert.True(t, o.IsPresent(), "present")
		assert.Equal(t, val, o.Get(), "value")
	}

	o = it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_iteratorStream_Map(t *testing.T) {
	double := func(n int) int {
		return n * 2
	}
	stream := FromArray(streamTestData).Map(double)
	it := stream.Iterator()

	var o fluent.Option[int]
	for i := range streamTestData {
		val := streamTestData[i] * 2
		o = it.Next()

		assert.NotNil(t, o, "option")
		assert.True(t, o.IsPresent(), "present")
		assert.Equal(t, val, o.Get(), "value")
	}

	o = it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_mapper_Size(t *testing.T) {
	expected := len(streamTestData)
	actual := mapper[int]{
		mapper: func(i int) int { return i },
		source: iterator.FromArray(streamTestData),
	}.Size()

	assert.NotNil(t, actual, "option")
	assert.True(t, actual.IsPresent(), "present")
	assert.Equal(t, expected, actual.Get(), "size")
}

func Test_iteratorStream_Count(t *testing.T) {
	size := FromArray(streamTestData).Count()
	assert.Equal(t, len(streamTestData), size, "size")
}

func Test_iteratorStream_Count_parallel(t *testing.T) {
	size := FromArray(streamTestData).Parallel().Count()
	assert.Equal(t, len(streamTestData), size, "size")
}

func Test_iteratorStream_Array_even(t *testing.T) {
	isEven := func(n int) bool {
		return n%2 == 0
	}
	arr := FromArray(streamTestData).Filter(isEven).Array()
	assert.Equal(t, len(streamTestData)/2, len(arr), "size")

	for _, val := range arr {
		assert.True(t, val%2 == 0, "event val")
	}
}

func Test_iteratorStream_Array(t *testing.T) {
	arr := FromArray(streamTestData).Array()
	for i := 0; i < len(streamTestData); i++ {
		actual := int(arr[i])
		expected := int(streamTestData[i])
		assert.Equal(t, expected, actual, "values")
	}
	assert.Equal(t, streamTestData, arr, "arrays")
}

func Test_iteratorStream_Array_parallel(t *testing.T) {
	arr := FromArray(streamTestData).Parallel().Array()
	for i := 0; i < len(streamTestData); i++ {
		actual := int(streamTestData[i])
		expected := int(arr[i])
		assert.Equal(t, expected, actual, "values")
	}
	assert.Equal(t, streamTestData, arr, "arrays")
}

func Test_iteratorStream_Peek(t *testing.T) {
	count := 0
	arr := make([]int, len(streamTestData))

	counter := func(val int) {
		arr[count] = val
		count++
	}

	FromArray(streamTestData).Peek(counter).Count()

	assert.Equal(t, len(streamTestData), count, "size")
	assert.Equal(t, streamTestData, arr, "size")
}

func Test_peek_Size(t *testing.T) {
	expected := len(streamTestData)
	actual := peek[int]{
		consumer: func(i int) {},
		source:   iterator.FromArray(streamTestData),
	}.Size()

	assert.NotNil(t, actual, "option")
	assert.True(t, actual.IsPresent(), "present")
	assert.Equal(t, expected, actual.Get(), "size")
}

func Test_iteratorStream_ForEach(t *testing.T) {
	count := 0
	arr := make([]int, len(streamTestData))

	counter := func(_, val int) {
		arr[count] = val
		count++
	}

	FromArray(streamTestData).ForEach(counter)

	assert.Equal(t, len(streamTestData), count, "size")
	assert.Equal(t, streamTestData, arr, "size")
}

func Test_iteratorStream_ForEach_Parallel(t *testing.T) {
	var count int32 = 0
	arr := make([]int, len(streamTestData))

	counter := func(_, val int) {
		fmt.Println("parallel", val)
		arr[atomic.LoadInt32(&count)] = val
		atomic.AddInt32(&count, 1)
	}

	FromArray(streamTestData).Parallel().ForEach(counter)

	originalSet := set.FromArray(streamTestData)
	processedSet := set.FromArray(arr)

	assert.Equal(t, int32(len(streamTestData)), count, "size")
	assert.True(t, processedSet.ContainsAll(originalSet), "content")
}
