package fluent

import (
	"embed"
	"fmt"
)

//go:embed README.md
var sourceFiles embed.FS

func FileBytes(fileName string) Result[[]byte] {
	// ResultFromCall returns an error Result if the result err is equal to nil
	// or Ok if the error is not present.
	return ResultFromCall(func() ([]byte, error) {
		data, err := sourceFiles.ReadFile(fileName)
		return data, err
	})
}

func ExampleResult_goodFile() {
	r := FileBytes("README.md")
	msg := ResultMap(r, func(b []byte) string {
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
	msg := ResultMap(r, func(b []byte) string {
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
