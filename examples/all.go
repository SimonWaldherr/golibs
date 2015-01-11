package main

import (
	"fmt"
	"github.com/simonwaldherr/golibs/as"
	"github.com/simonwaldherr/golibs/cli"
	"github.com/simonwaldherr/golibs/stack"
	"github.com/simonwaldherr/golibs/xmath"
)

func main() {
	floats := []float64{1.4, 3.14, 9.81, 13.2, 23.42, 33.7, 44.11, 51}
	array := stack.Lifo()
	array.Push(fmt.Sprintf("Lorem %v", cli.Color("Ipsum", cli.Blue)))
	array.Push(fmt.Sprintf("Dolor %v", cli.Bold("sit Amet")))
	array.Push(fmt.Sprintf("5th Prime: %v", xmath.Prime(5)))
	array.Push(fmt.Sprintf("Min: %v, Max: %v, Median: %v, Avg: %v", xmath.MinFloat(floats), xmath.MaxFloat(floats), xmath.MedianFloat(floats), xmath.AvgFloat(floats)))
	array.Push(fmt.Sprintf("Date: %v", as.Time("11.01.2015")))
	for array.Len() > 0 {
		fmt.Println(array.Pop())
	}
}
