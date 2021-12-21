# fluent [![Build status][ci-img]][ci-url] [![Test coverage][cov-img]][cov-url]
Functional constructs (lent from other languages) for Go lang 1.18

[ci-img]: https://github.com/mikhasd/fluent/actions/workflows/go.yml/badge.svg
[ci-url]: https://github.com/mikhasd/fluent/actions/workflows/go.yml
[cov-img]: https://codecov.io/gh/mikhasd/fluent/branch/main/graph/badge.svg
[cov-url]: https://codecov.io/gh/mikhasd/fluent/branch/main

# Option Examples

Option represents an optional value.

Option can be used as alternative for:
  - Uninitialized values
  - Invalid/empty return values
  - os.IsNotExist(e error)

This implementation is based on Java's java.util.Optional and Rust's std::option.

```go
package main

import (
  "github.com/mikhasd/fluent" 
  "fmt"
)

func Divide(a, b int) fluent.Option[int] {
	if b == 0 {
		return fluent.OptionEmpty[int]()
	} else {
		return fluent.OptionPresent(a / b)
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

# Result Examples

Result represents the output of a function which may have been computed successfully (`Ok`) or failed with an error (`Err`).

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
	return fluent.ResultFromCall(func() ([]byte, error) {
		data, err := sourceFiles.ReadFile(fileName)
		return data, err
	})
}

func ExampleResult_goodFile() {
	r := FileBytes("README.md")
	msg := fluent.ResultMap(r, func(b []byte) string {
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
	msg := fluent.ResultMap(r, func(b []byte) string {
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