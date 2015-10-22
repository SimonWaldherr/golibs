package file_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/file"
)

func ExampleGetAbsolutePath() {
	converted, err := file.GetAbsolutePath(".")
	fmt.Println(converted, err)
}

func ExampleRead() {
	str, _ := file.Read("test.txt")
	fmt.Println(str)

	// Output:
	// Lorem ipsum dolor sit amet, consectetur adipisici elit, sed eiusmod tempor incidunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquid ex ea commodi consequat. Quis aute iure reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint obcaecat cupiditat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
}

func ExampleSize() {
	size, _ := file.Size("test.txt")
	fmt.Println(size)

	// Output: 429
}
