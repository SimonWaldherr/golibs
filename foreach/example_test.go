package foreach_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/foreach"
)

type Foobar struct {
	X int
	Y int
}
type Foo struct {
	Id    int
	Name  string
	Lorem Foobar
	Ipsum interface{}
	Query func()
}

func ExampleStruct() {
	foreach.Struct(Foo{
		Id:   5,
		Name: "Foobar",
		Lorem: Foobar{
			X: 42,
			Y: 23,
		},
		Ipsum: float64(2),
		Query: func() {
			fmt.Println("---")
		},
	}, func(key string, index string, value interface{}, depth int) {
		fmt.Printf("name: %v, type: %v, depth: %v\n", key, index, depth)
	})

	// Output:
	//name: Id, type: int, depth: 0
	//name: Name, type: string, depth: 0
	//name: Lorem, type: foreach_test.Foobar, depth: 0
	//name: X, type: int, depth: 1
	//name: Y, type: int, depth: 1
	//name: Ipsum, type: interface {}, depth: 0
	//name: Query, type: func(), depth: 0
}
