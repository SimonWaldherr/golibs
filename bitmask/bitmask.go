// Package bitmask provides a way to store multiple boolean values in a single integer.
package bitmask

import (
	"fmt"
)

// Bitmask represents an integer that stores multiple boolean values using bit manipulation.
type Bitmask struct {
	value int
}

// New creates a new Bitmask with the given initial value.
func New(init int) *Bitmask {
	return &Bitmask{value: init}
}

func (b *Bitmask) setBit(pos int) int {
	b.value |= (1 << pos)
	return b.value
}

func (b *Bitmask) clearBit(pos int) int {
	mask := ^(1 << pos)
	b.value &= mask
	return b.value
}

// Set sets the bit at position pos to val
func (b *Bitmask) Set(pos int, val bool) int {
	if val == true {
		return b.setBit(pos)
	}
	return b.clearBit(pos)
}

// hasBit checks if the bit at position pos is set
func (b *Bitmask) hasBit(pos int) bool {
	return ((b.value & (1 << pos)) > 0)
}

// Get returns the value of the bit at position pos
func (b *Bitmask) Get(pos int) bool {
	return b.hasBit(pos)
}

// Int returns the integer value of the bitmask
func (b *Bitmask) Int() int {
	return b.value
}

// String returns the string representation of the bitmask
func (b *Bitmask) String() string {
	return fmt.Sprintf("%b", b.value)
}

// Byte returns the byte representation of the bitmask
func (b *Bitmask) Byte() []byte {
	return []byte{byte(b.value)}
}

// FromByte initializes the bitmask from a byte.
func FromByte(b byte) *Bitmask {
	return &Bitmask{value: int(b)}
}

// Toggle toggles the bit at the given position.
func (b *Bitmask) Toggle(pos int) {
	b.value ^= (1 << pos)
}

// Count returns the number of set bits (1s) in the bitmask.
func (b *Bitmask) Count() int {
	count := 0
	value := b.value
	for value > 0 {
		count += value & 1
		value >>= 1
	}
	return count
}

// Reset resets all bits to 0.
func (b *Bitmask) Reset() {
	b.value = 0
}

// And performs bitwise AND operation with another bitmask.
func (b *Bitmask) And(other *Bitmask) *Bitmask {
	return &Bitmask{value: b.value & other.value}
}

// Or performs bitwise OR operation with another bitmask.
func (b *Bitmask) Or(other *Bitmask) *Bitmask {
	return &Bitmask{value: b.value | other.value}
}

// Xor performs bitwise XOR operation with another bitmask.
func (b *Bitmask) Xor(other *Bitmask) *Bitmask {
	return &Bitmask{value: b.value ^ other.value}
}
