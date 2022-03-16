# fluent [![Build status][ci-img]][ci-url] [![Test coverage][cov-img]][cov-url]
Functional constructs (lent from other languages) for Go lang 1.18

[ci-img]: https://github.com/mikhasd/fluent/actions/workflows/go.yml/badge.svg
[ci-url]: https://github.com/mikhasd/fluent/actions/workflows/go.yml
[cov-img]: https://codecov.io/gh/mikhasd/fluent/branch/main/graph/badge.svg
[cov-url]: https://codecov.io/gh/mikhasd/fluent/branch/main

- [Option](#option)
  - [API](#option-api)
  - [Examples](#option-examples)
- [Result](#result)
  - [API](#result-api)
  - [Examples](#result-examples)
- [array](#array)
  - [API](#array-api)
- [Iterator](#iterator)
  - [API](#iterator-api)

# Option

Option represents an optional value.

Option can be used as alternative for:
  - Uninitialized values
  - Invalid/empty return values
  - os.IsNotExist(e error)

This implementation is based on Java's java.util.Optional and Rust's std::option.

## Option API

```go
type Option[T any] interface {

        // Returns true if a value is present.
        IsPresent() bool

        // Gets the Option value or panics if empty
        Get() T

        // Applies the provided mapper function over the Option value, if present.
        Map(mapper func(T) T) Option[T]

        // Returns the Option value if present, or the provided value it empty.
        OrElse(other T) T

        // Returns the Option value if present, or invoke the provided function if
        // empty.
        OrElseGet(func() T) T

        // If empty, calls the provided function and returns it result or returns
        // the current option if present.
        Or(func() Option[T]) Option[T]

        // Returns the Option value wrapped into a result if present. If empty,
        // returns a Result with a wrapped error.
        OrError(e error) Result[T]

        // Executes the provided function if a value is present
        IfPresent(func(T))

        // If a value is present and the value is accepted by the provided filter,
        // returns an Option representing the value, otherwise will return an empty
        // Option
        Filter(func(T) bool) Option[T]

        String() string
}  

// Present returns a new Option wrapping the T value.
func Present[T any](value T) Option[T]

// Empty returns an empty Option.
func Empty[T any]() Option[T]

// MapOption executes the mapper function over the Option value if it is not
// empty.
//
// If the Option is empty, the mapper is not executed and an empty Option is
// returned.
func MapOption[T any, R any](o Option[T], mapper func(T) R) Option[R]

// OfNillable returns an Option describing the provided reference.
//
// If the provided reference of T is non-nil, a Option representing the value
// is returned.
//
// If the value is nil, an empty option is returned.
func OfNillable[T any](ref *T) Option[*T]
```

## Option Examples

```go
package main

import (
  "github.com/mikhasd/fluent" 
  "fmt"
)

func Divide(a, b int) fluent.Option[int] {
	if b == 0 {
		return fluent.Empty[int]()
	} else {
		return fluent.Present(a / b)
	}
}

func Double(a int) int {
	return a * 2
}

func String(a int) string {
	return fmt.Sprintf("%d", a)
}

func ExampleOption_goodDivision() {
	// Mapping operations can be chained if input and output types are the same
	var option fluent.Option[int] = Divide(6, 3).Map(Double)
	// But mapping operations with different input and output types must use
	// package level functions due to language limitations.
	message := fluent.MapOption(option, String)

	var result string
	if message.Present() {
		result = message.Get()
	} else {
		result = "empty"
	}

	fmt.Println(result)
	// Output: 4
}

func ExampleOption_badDivision() {
	var option fluent.Option[int] = Divide(100, 0).Map(Double)
	message := fluent.MapOption(option, String)

	var result string
	if message.Present() {
		result = message.Get()
	} else {
		result = "empty"
	}

	fmt.Println(result)
	// Output: empty
}
```

# Result

Result represents the output of a function which may have been computed successfully (`Ok`) or failed with an error (`Err`).

## Result API

```go
type Result[T any] interface {
	// Returns true if the Result is Ok
	IsOk() bool

	// Returns true if the Result is Err
	IsErr() bool

	// Converts the Result into an Option.
	//
	// If the Result is Ok, an Option will be returned, wrapping the result
	// value.
	//
	// If the Result is Err, an empty Option will be returned.
	Ok() Option[T]

	// Converts the Result into an Option.
	//
	// If the Result is Err, an Option will be returned, wrapping the result
	// error.
	//
	// If the Result is Ok, an empty Option will be returned.
	Err() Option[error]

	// Executes the mapper function over the Result value.
	//
	// If the Result is Err, the function is not executed.
	Map(mapper func(T) T) Result[T]

	// Executes the mapper function over the Result error.
	//
	// If the Result is Ok, the function is not executed.
	MapErr(func(error) T) Result[T]

	// Gets the Result value.
	//
	// This function will panic if the Result is Err.
	Get() T

	// Gets the Result error.
	//
	// This function will panic if the Result is Ok.
	GetErr() error

	// Gets the Result value or	the provided value if if the Result is Err.
	OrElse(value T) T

	// Calls the provided function if the Result is Err. Returns Result value
	// if Result is Ok.
	OrElseGet(func() T) T

	// Calls the provided function if the Result is Err. Returns the current
	// Result if Ok.
	Or(func() Result[T]) Result[T]

	String() string
}

// Ok returns a result representing a successful operation resulting in
// the value T.
func Ok[T any](value T) Result[T]

// Err returns a result representing a failed operation which resulted in
// on an error.
func Err[T any](e error) Result[T]

// MapResult executes the mapper function over the Result value if it Ok.
//
// If the Result is Err, the mapper is not executed and an Err Result is
// returned.
func MapResult[T any, R any](r Result[T], mapper func(T) R) Result[R]

// CallResult executes the provided function and returns an Ok result if
// the result error is nil.
//
// If the provided function returns an error, an Err result is returned, even
// if a result is returned.
func CallResult[T any](fn func() (T, error)) Result[T]
```

## Result Examples

```go
package main

import (
	"github.com/mikhasd/fluent"
	"embed"
	"fmt"
)

//go:embed README.md
var sourceFiles embed.FS

func FileBytes(fileName string) fluent.Result[[]byte] {
	// ResultFromCall returns an error Result if the result err is equal to nil
	// or Ok if the error is not present.
	return fluent.CallResult(func() ([]byte, error) {
		data, err := sourceFiles.ReadFile(fileName)
		return data, err
	})
}

func ExampleResult_goodFile() {
	r := FileBytes("README.md")
	msg := fluent.MapResult(r, func(b []byte) string {
		return "has file"
	})

	var result string
	if msg.IsOk() {
		result = msg.Get()
	} else {
		result = "empty"
	}

	fmt.Println(result)
	// Output: has file
}

func ExampleResult_badFile() {
	r := FileBytes("badfile")
	msg := fluent.MapResult(r, func(b []byte) string {
		return "has file"
	})

	var result string
	if msg.IsOk() {
		result = msg.Get()
	} else {
		result = "empty"
	}

	fmt.Println(result)
	// Output: empty
}

```

# array

The `array` package contains functions to facilitate working with arrays.

## array API

```go
// Map creates a new array populated with the result of calling the provided
// `mapper` function on every element of the input array.
func Map[I any, O any](in []I, mapper func(I) O) []O

// Filter creates a new array with all elements that pass the provided
// `condition` function.
func Filter[T any](in []T, condition func(T) bool) []T

```

# Iterator

An Iterator facilitates traversing a collection of elements of know or unknown size. 

## Iterator API

```go
type Iterator[T any] interface {
	// Next advances the iteration and return the next value.
	// An empty `fluent.Option` will be returned when the iteration finishes.
	Next() fluent.Option[T]
}

// FromArray creates a new iterator for a given array.
func FromArray[T any](elements []T) Iterator[T]

// Of creates a new `Iterator` for the elements provided as arguments.
func Of[T any](elements ...T) Iterator[T]

// MapKeys creates an `Iterator` with the keys of a given map.
func MapKeys[K comparable, V any](m map[K]V) Iterator[K]

// MapValues creates an `Iterator` with the values of a given map.
func MapValues[K comparable, V any](m map[K]V) Iterator[V]

// FromMap creates a new `Iterator` for the keys and values of a given map.
func FromMap[K comparable, V any](m map[K]V) Iterator[MapEntry[K, V]]
```

# stream