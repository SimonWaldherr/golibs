package set

import (
	"math"
	"sort"
	"testing"
)

func assertValues(t *testing.T, got Set[int], want ...int) {
	t.Helper()
	values := got.Values()
	sort.Ints(values)
	if len(values) != len(want) || got.Len() != len(want) {
		t.Fatalf("set values = %v, Len = %d, want %v", values, got.Len(), want)
	}
	for i, value := range values {
		if value != want[i] {
			t.Fatalf("set values = %v, want %v", values, want)
		}
	}
}

func TestZeroValueAndMembership(t *testing.T) {
	var values Set[int]
	if values.Len() != 0 || values.Has(0) || values.Values() != nil {
		t.Fatal("zero set was not empty")
	}
	values.Add()
	values.Remove(3)
	if values.values != nil {
		t.Fatal("empty operations allocated storage")
	}
	values.Add(3, 1, 3, 2, 0)
	assertValues(t, values, 0, 1, 2, 3)
	if !values.Has(0) || values.Has(99) {
		t.Fatal("membership does not distinguish present and absent values")
	}
	values.Remove(0, 3, 99, 3)
	assertValues(t, values, 1, 2)
	values.Remove(1, 2)
	if values.Values() != nil {
		t.Fatal("empty populated set should also return nil Values")
	}
}

func TestCopies(t *testing.T) {
	original := New(1, 2)
	shared := original
	shared.Add(3)
	if !original.Has(3) {
		t.Fatal("copy of populated set did not share its documented storage")
	}
	clone := original.Clone()
	clone.Remove(1)
	clone.Add(4)
	assertValues(t, original, 1, 2, 3)
	assertValues(t, clone, 2, 3, 4)
	values := original.Values()
	for i := range values {
		values[i] = 99
	}
	assertValues(t, original, 1, 2, 3)
	empty := New[int]().Clone()
	empty.Add(1)
	assertValues(t, empty, 1)
}

func TestSetAlgebra(t *testing.T) {
	cases := []struct {
		name                            string
		left, right                     []int
		union, intersection, difference []int
	}{
		{"both empty", nil, nil, nil, nil, nil},
		{"empty left", nil, []int{2}, []int{2}, nil, nil},
		{"empty right", []int{1}, nil, []int{1}, nil, []int{1}},
		{"equal", []int{1, 1, 2}, []int{2, 1}, []int{1, 2}, []int{1, 2}, nil},
		{"disjoint", []int{1, 2}, []int{3}, []int{1, 2, 3}, nil, []int{1, 2}},
		{"overlap", []int{1, 2, 3}, []int{2, 3, 4}, []int{1, 2, 3, 4}, []int{2, 3}, []int{1}},
		{"right larger", []int{2}, []int{1, 2, 3}, []int{1, 2, 3}, []int{2}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			left, right := New(tc.left...), New(tc.right...)
			union, intersection, difference := left.Union(right), left.Intersection(right), left.Difference(right)
			assertValues(t, union, tc.union...)
			assertValues(t, right.Union(left), tc.union...)
			assertValues(t, intersection, tc.intersection...)
			assertValues(t, right.Intersection(left), tc.intersection...)
			assertValues(t, difference, tc.difference...)
			// All algebra results must be independently mutable.
			union.Add(99)
			intersection.Add(98)
			difference.Add(97)
			for _, key := range []int{97, 98, 99} {
				if left.Has(key) || right.Has(key) {
					t.Fatal("algebra result mutated an operand")
				}
			}
		})
	}
}

func TestComparableStructs(t *testing.T) {
	type key struct {
		Name string
		ID   int
	}
	values := New(key{"item", 1}, key{"item", 1}, key{"item", 2})
	if values.Len() != 2 || !values.Has(key{"item", 2}) {
		t.Fatal("struct keys lost equality semantics")
	}
}

func TestEmptyOperandCopiesPreserveFloatKeys(t *testing.T) {
	const nanBits = 0x7ff8000000000042
	negativeZero := math.Copysign(0, -1)
	original := New(math.Float64frombits(nanBits), negativeZero)
	var empty Set[float64]
	for name, result := range map[string]Set[float64]{
		"union":      empty.Union(original),
		"difference": original.Difference(empty),
	} {
		t.Run(name, func(t *testing.T) {
			foundNaN, foundZero := false, false
			for _, value := range result.Values() {
				switch math.Float64bits(value) {
				case nanBits:
					foundNaN = true
				case math.Float64bits(negativeZero):
					foundZero = true
				default:
					t.Fatalf("float key changed: %v", value)
				}
			}
			if result.Len() != 2 || !foundNaN || !foundZero {
				t.Fatal("empty operand copy lost a stored float key")
			}
			result.Remove(0)
			result.Add(42)
			if original.Len() != 2 || !original.Has(0) || original.Has(42) {
				t.Fatal("empty operand result shares its source storage")
			}
		})
	}
}
