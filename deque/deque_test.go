package deque

import (
	"reflect"
	"testing"
)

func TestZeroAndStoredZeroValues(t *testing.T) {
	var values Deque[*int]
	if values.Len() != 0 || values.Cap() != 0 || values.Values() != nil {
		t.Fatal("zero deque is not empty")
	}
	for _, get := range []func() (*int, bool){values.Front, values.Back, values.PopFront, values.PopBack} {
		if value, ok := get(); value != nil || ok {
			t.Fatal("empty operation did not return nil, false")
		}
	}
	values.Clear()
	copy := values.Clone()
	if copy.Len() != 0 || copy.Cap() != 0 {
		t.Fatal("empty clone is not the zero value")
	}
	values.PushFront(nil)
	if value, ok := values.Front(); value != nil || !ok {
		t.Fatal("stored nil is indistinguishable from an empty deque")
	}
	if value, ok := values.PopBack(); value != nil || !ok {
		t.Fatal("stored nil was not removed")
	}
}

func TestMixedEnds(t *testing.T) {
	var values Deque[int]
	values.PushBack(1)
	values.PushBack(2)
	values.PushFront(0)
	values.PushFront(-1)
	if got := values.Values(); !reflect.DeepEqual(got, []int{-1, 0, 1, 2}) {
		t.Fatalf("order = %v", got)
	}
	if value, ok := values.Front(); value != -1 || !ok {
		t.Fatalf("front = %v, %v", value, ok)
	}
	if value, ok := values.Back(); value != 2 || !ok {
		t.Fatalf("back = %v, %v", value, ok)
	}
	if value, ok := values.PopFront(); value != -1 || !ok {
		t.Fatalf("pop front = %v, %v", value, ok)
	}
	if value, ok := values.PopBack(); value != 2 || !ok {
		t.Fatalf("pop back = %v, %v", value, ok)
	}
	if value, _ := values.PopFront(); value != 0 {
		t.Fatalf("next front = %v", value)
	}
	if value, _ := values.PopBack(); value != 1 {
		t.Fatalf("last value = %v", value)
	}
	if values.Len() != 0 || values.head != 0 {
		t.Fatal("empty state was not restored")
	}
}

func TestWrappedGrowthAndPushFront(t *testing.T) {
	var values Deque[int]
	for i := 0; i < 8; i++ {
		values.PushBack(i)
	}
	for i := 0; i < 5; i++ {
		if got, ok := values.PopFront(); got != i || !ok {
			t.Fatalf("initial pop = %d, %v; want %d", got, ok, i)
		}
	}
	for i := 8; i < 13; i++ {
		values.PushBack(i)
	}
	if values.head == 0 {
		t.Fatal("test did not exercise a wrapped ring")
	}
	capacity := values.Cap()
	values.PushFront(4)
	if values.Cap() <= capacity {
		t.Fatal("full wrapped ring did not grow")
	}
	for want := 12; want >= 4; want-- {
		if got, ok := values.PopBack(); got != want || !ok {
			t.Fatalf("after growth pop = %d, %v; want %d", got, ok, want)
		}
	}
	for i := 0; i < 80; i++ {
		values.PushFront(i)
	}
	for want := 79; want >= 0; want-- {
		if got, ok := values.PopFront(); got != want || !ok {
			t.Fatalf("front-only growth pop = %d, %v; want %d", got, ok, want)
		}
	}
}

func TestValuesAndCloneAreIndependent(t *testing.T) {
	var original Deque[int]
	for i := 0; i < 8; i++ {
		original.PushBack(i)
	}
	for i := 0; i < 3; i++ {
		original.PopFront()
	}
	original.PushBack(8)
	original.PushBack(9)
	snapshot := original.Values()
	copy := original.Clone()
	want := []int{3, 4, 5, 6, 7, 8, 9}
	if !reflect.DeepEqual(copy.Values(), want) {
		t.Fatalf("wrapped clone = %v", copy.Values())
	}
	snapshot[0] = 99
	copy.PopFront()
	copy.PushBack(10)
	if !reflect.DeepEqual(original.Values(), want) {
		t.Fatalf("clone/snapshot mutated original: %v", original.Values())
	}
	original.Clear()
	if !reflect.DeepEqual(copy.Values(), []int{4, 5, 6, 7, 8, 9, 10}) {
		t.Fatalf("clear original mutated clone: %v", copy.Values())
	}

	value := 1
	var pointers Deque[*int]
	pointers.PushBack(&value)
	shallow := pointers.Clone()
	if got, _ := shallow.Front(); got != &value {
		t.Fatal("Clone unexpectedly changed element identity")
	}
}

func TestRemovedAndClearedReferencesAreReleased(t *testing.T) {
	var values Deque[*int]
	for i := 0; i < 8; i++ {
		value := i
		values.PushBack(&value)
	}
	values.PopFront()
	if values.buffer[0] != nil {
		t.Fatal("PopFront retained its removed element")
	}
	values.PopBack()
	if values.buffer[7] != nil {
		t.Fatal("PopBack retained its removed element")
	}
	for i := 0; i < 3; i++ {
		values.PopFront()
	}
	value := 9
	values.PushBack(&value)
	values.PushBack(&value)
	values.PushBack(&value)
	values.PushBack(&value)
	if values.head+values.Len() <= values.Cap() {
		t.Fatal("test did not exercise a wrapped clear")
	}
	capacity := values.Cap()
	values.Clear()
	if values.Len() != 0 || values.Cap() != capacity || values.Values() != nil {
		t.Fatal("Clear failed to retain only empty capacity")
	}
	for i, value := range values.buffer {
		if value != nil {
			t.Fatalf("Clear retained reference in slot %d", i)
		}
	}
	values.PushBack(&value)
	if got, ok := values.PopFront(); got != &value || !ok {
		t.Fatal("deque could not be reused after Clear")
	}
}

func TestCloneDoesNotRetainExcessCapacity(t *testing.T) {
	var values Deque[int]
	for i := 0; i < 1024; i++ {
		values.PushBack(i)
	}
	full := values.Clone()
	first, _ := full.Front()
	last, _ := full.Back()
	if full.Len() != 1024 || first != 0 || last != 1023 {
		t.Fatalf("full clone lost values: length %d, endpoints %d/%d", full.Len(), first, last)
	}
	for values.Len() > 3 {
		values.PopFront()
	}
	clone := values.Clone()
	if clone.Cap() >= values.Cap() || !reflect.DeepEqual(clone.Values(), []int{1021, 1022, 1023}) {
		t.Fatalf("sparse clone: values %v, cap %d; original cap %d", clone.Values(), clone.Cap(), values.Cap())
	}
	// Continue across the compact ring boundary to verify it remains a deque.
	for i := 0; i < 20; i++ {
		clone.PushBack(i)
	}
	if !reflect.DeepEqual(values.Values(), []int{1021, 1022, 1023}) {
		t.Fatal("growing the compact clone modified its source")
	}
}

func FuzzDequeAgainstSlice(f *testing.F) {
	for _, seed := range [][]byte{nil, {0, 1, 2, 3, 4, 5, 6, 7}, {255, 128, 64, 0}, make([]byte, 128)} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, operations []byte) {
		if len(operations) > 2048 {
			operations = operations[:2048]
		}
		var values Deque[int]
		var model []int
		for _, operation := range operations {
			value := int(int8(operation))
			switch operation & 7 {
			case 0:
				if len(model) < 64 {
					values.PushFront(value)
					model = append([]int{value}, model...)
				}
			case 1:
				if len(model) < 64 {
					values.PushBack(value)
					model = append(model, value)
				}
			case 2:
				got, ok := values.PopFront()
				if len(model) == 0 {
					if got != 0 || ok {
						t.Fatalf("empty PopFront = %d, %v", got, ok)
					}
				} else {
					if !ok || got != model[0] {
						t.Fatalf("PopFront = %d, %v; model %v", got, ok, model)
					}
					model = model[1:]
				}
			case 3:
				got, ok := values.PopBack()
				if len(model) == 0 {
					if got != 0 || ok {
						t.Fatalf("empty PopBack = %d, %v", got, ok)
					}
				} else {
					if !ok || got != model[len(model)-1] {
						t.Fatalf("PopBack = %d, %v; model %v", got, ok, model)
					}
					model = model[:len(model)-1]
				}
			case 4:
				got, ok := values.Front()
				if ok != (len(model) > 0) || ok && got != model[0] {
					t.Fatalf("Front = %d, %v; model %v", got, ok, model)
				}
			case 5:
				got, ok := values.Back()
				if ok != (len(model) > 0) || ok && got != model[len(model)-1] {
					t.Fatalf("Back = %d, %v; model %v", got, ok, model)
				}
			case 6:
				values.Clear()
				model = nil
			case 7:
				copy := values.Clone()
				if !equalValues(copy.Values(), model) {
					t.Fatalf("Clone = %v; model %v", copy.Values(), model)
				}
				copy.Clear()
			}
			if values.Len() != len(model) || !equalValues(values.Values(), model) {
				t.Fatalf("deque = %v; model %v", values.Values(), model)
			}
		}
	})
}

func equalValues(a, b []int) bool {
	return len(a) == len(b) && reflect.DeepEqual(append([]int(nil), a...), append([]int(nil), b...))
}
