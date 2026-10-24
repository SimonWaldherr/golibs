package xmath

import (
	"math"
	"reflect"

	"simonwaldherr.de/go/golibs/as"
)

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type reduction uint8

const (
	sumReduction reduction = iota
	minReduction
	maxReduction
	arithmeticReduction
	rmsReduction
	harmonicReduction
	geometricReduction
)

// Dispatch once per call, so numeric loops need neither reflection nor boxed values.
// Named types and mixed values retain the conversion behavior of as.Float.
func reduce(input interface{}, mode reduction) float64 {
	switch values := input.(type) {
	case []float64:
		return reduceNumbers(values, mode)
	case []float32:
		return reduceNumbers(values, mode)
	case []int:
		return reduceNumbers(values, mode)
	case []int8:
		return reduceNumbers(values, mode)
	case []int16:
		return reduceNumbers(values, mode)
	case []int32:
		return reduceNumbers(values, mode)
	case []int64:
		return reduceNumbers(values, mode)
	case []uint:
		return reduceNumbers(values, mode)
	case []uint8:
		return reduceNumbers(values, mode)
	case []uint16:
		return reduceNumbers(values, mode)
	case []uint32:
		return reduceNumbers(values, mode)
	case []uint64:
		return reduceNumbers(values, mode)
	default:
		return reduceNumbers(numericCopy(input), mode)
	}
}

func reduceNumbers[T number](values []T, mode reduction) float64 {
	var result float64
	// Select the operation once, keeping the numeric loops free of mode branches.
	switch mode {
	case minReduction:
		result = float64(values[0])
		for _, value := range values[1:] {
			x := float64(value)
			if x < result {
				result = x
			}
		}
	case maxReduction:
		result = float64(values[0])
		for _, value := range values[1:] {
			x := float64(value)
			if x > result {
				result = x
			}
		}
	case rmsReduction:
		for _, value := range values {
			x := float64(value)
			result += x * x
		}
		return math.Sqrt(result / float64(len(values)))
	case harmonicReduction:
		for _, value := range values {
			result += 1 / float64(value)
		}
		return float64(len(values)) / result
	case geometricReduction:
		result = 1
		for _, value := range values {
			result *= float64(value)
		}
		return math.Pow(result, 1/float64(len(values)))
	default:
		for _, value := range values {
			result += float64(value)
		}
		if mode == arithmeticReduction {
			return result / float64(len(values))
		}
	}
	return result
}

func numericCopy(input interface{}) []float64 {
	switch values := input.(type) {
	case []float64:
		return copyNumbers(values)
	case []float32:
		return copyNumbers(values)
	case []int:
		return copyNumbers(values)
	case []int8:
		return copyNumbers(values)
	case []int16:
		return copyNumbers(values)
	case []int32:
		return copyNumbers(values)
	case []int64:
		return copyNumbers(values)
	case []uint:
		return copyNumbers(values)
	case []uint8:
		return copyNumbers(values)
	case []uint16:
		return copyNumbers(values)
	case []uint32:
		return copyNumbers(values)
	case []uint64:
		return copyNumbers(values)
	default:
		slice := reflect.ValueOf(input)
		out := make([]float64, slice.Len())
		for i := range out {
			out[i] = as.Float(slice.Index(i).Interface())
		}
		return out
	}
}

func copyNumbers[T number](values []T) []float64 {
	out := make([]float64, len(values))
	for i, value := range values {
		out[i] = float64(value)
	}
	return out
}
