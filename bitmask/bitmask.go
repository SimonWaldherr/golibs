// bitmasks are a way to store multiple boolean values in a single integer
package bitmask

import "fmt"

type Bitmask struct {
	value int
}

// New creates a new Bitmask
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
