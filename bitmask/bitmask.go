package bitmask

import "fmt"

type Bitmask struct {
	value int
}

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

func (b *Bitmask) Set(pos int, val bool) int {
	if val == true {
		return b.setBit(pos)
	}
	return b.clearBit(pos)
}

func (b *Bitmask) hasBit(pos int) bool {
	return ((b.value & (1 << pos)) > 0)
}

func (b *Bitmask) Get(pos int) bool {
	return b.hasBit(pos)
}

func (b *Bitmask) Int() int {
	return b.value
}

func (b *Bitmask) String() string {
	return fmt.Sprintf("%b", b.value)
}
