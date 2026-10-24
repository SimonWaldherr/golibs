package set_test

import (
	"fmt"
	"sort"

	"simonwaldherr.de/go/golibs/set"
)

func Example() {
	var languages set.Set[string]
	languages.Add("Go", "Rust", "Go")
	languages.Remove("Rust")
	combined := languages.Union(set.New("Python", "Go"))
	values := combined.Values()
	sort.Strings(values) // Set iteration order is unspecified.
	fmt.Println(languages.Has("Go"), languages.Len())
	fmt.Println(values)
	// Output:
	// true 1
	// [Go Python]
}
