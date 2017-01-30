package structs

import (
	"fmt"
	"testing"
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

func Test_Reflect(t *testing.T) {
	s := Reflect(Foo{
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
	})

	if len(s) != 5 {
		t.Fatal("structs.Reflect Test 1 failed: wrong length\n")
	}

	if s["Query"] != "func()" || s["Id"] != "int" || s["Name"] != "string" {
		t.Fatalf("structs.Reflect Test 2 failed: wrong map struct \"%v\"\n", s)
	}
}
