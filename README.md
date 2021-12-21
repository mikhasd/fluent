# fluent-go
Functional constructs (lent from other languages) for Go lang 1.18

# Option Examples

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
