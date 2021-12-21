package fluent

import (
	"fmt"
)

func Divide(a, b int) Option[int] {
	if b == 0 {
		return Empty[int]()
	} else {
		return Present(a / b)
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
	var option Option[int] = Divide(6, 3).Map(Double)
	// But mapping operations with different input and output types must use
	// package level functions due to language limitations.
	message := MapOption(option, String)

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
	var option Option[int] = Divide(100, 0).Map(Double)
	message := MapOption(option, String)

	var result string
	if message.Present() {
		result = message.Get()
	} else {
		result = "empty"
	}

	fmt.Println(result)
	// Output: empty
}
