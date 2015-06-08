// +build local

package main

import (
	"fmt"
	"log"
	"simonwaldherr.de/go/golibs/ansi"
	"simonwaldherr.de/go/golibs/as"
	"simonwaldherr.de/go/golibs/regex"
	"simonwaldherr.de/go/golibs/stack"
	"simonwaldherr.de/go/golibs/xmath"
	"strconv"
)

func main() {
	floats := []float64{1.4, 3.14, 9.81, 13.2, 23.42, 33.7, 44.11, 51}
	array := stack.Lifo()
	array.Push(fmt.Sprintf("Lorem %v", ansi.Color("Ipsum", ansi.Blue)))
	array.Push(fmt.Sprintf("Dolor %v", ansi.Bold("sit Amet")))
	array.Push(fmt.Sprintf("5th Prime: %v", xmath.Prime(5)))
	array.Push(fmt.Sprintf("\n\tMin: %v\n\tMax: %v\n\tMedian: %v\n\tArithmetic: %v\n\tHarmonic: %v\n\tGeometric: %v", xmath.Min(floats), xmath.Max(floats), xmath.Median(floats), xmath.Arithmetic(floats), xmath.Harmonic(floats), xmath.Geometric(floats)))
	array.Push(fmt.Sprintf("Date: %v", as.Time("11.01.2015")))
	for array.Len() > 0 {
		log.Println(array.Pop())
	}

	array.Push(string(as.Bytes(regex.ReplaceAllString("foobar", "o+", "u"))))
	array.Push(as.Int(23.0000))
	array.Push(as.Float(13.37))
	array.Push(as.String(23.0))
	array.Push(as.Bool(111111))
	array.Push(as.Bytes(12.34))
	array.Push(as.String(strconv.ParseInt("42", 10, 0)))
	array.Push(as.FloatFromXString("2,3"))
	array.Push(as.FloatFromXString(".23"))
	array.Push(as.String("\r\n\t\n"))
	for array.Len() > 0 {
		log.Println(array.Pop())
	}
}
