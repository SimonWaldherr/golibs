// +build local

package main

import (
	"fmt"
	"os"
	"simonwaldherr.de/go/golibs/as"
	"simonwaldherr.de/go/golibs/xmath"
	"strings"
)

func asString(number float64, decimalplaces int) string {
	var r string
	number = xmath.FloatRound(number, decimalplaces)
	str := strings.Split(as.String(number), ".")
	if len(str) == 1 {
		r = str[0] + ".000"
	}
	if len(str) == 2 {
		r = str[0] + "." + as.FixedLengthAfter(str[1], "0", 3)
	}
	return as.FixedLengthBefore(r, " ", 12)
}

func main() {
	var l = len(os.Args)
	var f = make([]float64, 0)
	for i := 1; i < l; i++ {
		f = append(f, as.Float(os.Args[i]))
	}

	fmt.Printf("Arithmetic:     %v\n", asString(xmath.Mean(f, xmath.ArithmeticMean), 3))
	fmt.Printf("Geometric:      %v\n", asString(xmath.Mean(f, xmath.GeometricMean), 3))
	fmt.Printf("Harmonic:       %v\n", asString(xmath.Mean(f, xmath.HarmonicMean), 3))
	fmt.Printf("Median:         %v\n", asString(xmath.Mean(f, xmath.MedianMean), 3))
	fmt.Printf("Rootmeansquare: %v\n", asString(xmath.Mean(f, xmath.RmsMean), 3))
}
