package bitmask

import (
	"math/bits"
	"testing"
)

func TestSignBitOperations(t *testing.T) {
	position := bits.UintSize - 1
	b := New(0)
	b.Set(position, true)
	if !b.Get(position) {
		t.Error("highest bit was not set")
	}
	if got := b.Count(); got != 1 {
		t.Errorf("Count() = %d, want 1", got)
	}
	b.Set(0, true)
	if got := b.Count(); got != 2 {
		t.Errorf("Count() = %d, want 2", got)
	}
	b.Toggle(position)
	if b.Get(position) || b.Int() != 1 {
		t.Errorf("Toggle sign bit left value %d", b.Int())
	}
	b.Set(position, true)
	b.Set(position, false)
	if b.Get(position) || b.Int() != 1 {
		t.Errorf("clearing sign bit left value %d", b.Int())
	}
}

func TestNegativeMaskCounts(t *testing.T) {
	b := New(-1)
	if got := b.Count(); got != bits.UintSize {
		t.Errorf("Count(-1) = %d, want %d", got, bits.UintSize)
	}
	for position := 0; position < bits.UintSize; position++ {
		if !b.Get(position) {
			t.Errorf("bit %d of -1 was not set", position)
		}
	}
	b.Set(0, false)
	if got := b.Count(); got != bits.UintSize-1 {
		t.Errorf("count after clearing low bit = %d, want %d", got, bits.UintSize-1)
	}
	b.Reset()
	if b.Count() != 0 || b.Int() != 0 {
		t.Errorf("Reset() left value %d", b.Int())
	}
}

func TestBitwiseOperationsPreserveOperands(t *testing.T) {
	a, b := New(-1), New(0b1010)
	if got := a.And(b).Int(); got != 0b1010 {
		t.Errorf("And() = %d, want 10", got)
	}
	if got := a.Or(b).Int(); got != -1 {
		t.Errorf("Or() = %d, want -1", got)
	}
	if got := a.Xor(b).Int(); got != ^0b1010 {
		t.Errorf("Xor() = %d, want %d", got, ^0b1010)
	}
	if a.Int() != -1 || b.Int() != 0b1010 {
		t.Error("bitwise operation mutated an operand")
	}
}
