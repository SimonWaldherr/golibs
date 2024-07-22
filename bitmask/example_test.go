package bitmask_test

import (
	"fmt"

	"simonwaldherr.de/go/golibs/bitmask"
)

func Example() {
	// Initialize a bitmask with all bits set
	i := bitmask.New(0b11111111)

	// Clear the bit at position 0
	i.Set(0, false)
	fmt.Println(i) // Expected output: 11111110

	// Clear the bit at position 3
	i.Set(3, false)
	fmt.Println(i) // Expected output: 11110110

	// Set the bit at position 0
	i.Set(0, true)
	fmt.Println(i) // Expected output: 11110111

	// Get the value of the bit at position 2
	fmt.Println("2:", i.Get(2)) // Expected output: 2: true

	// Print the byte representation of the bitmask
	fmt.Printf("[]byte: %08b\n", i.Byte()) // Expected output: []byte: [11110111]

	// Toggle the bit at position 0
	i.Toggle(0)
	fmt.Println(i) // Expected output: 11110110

	// Toggle the bit at position 0 again
	i.Toggle(0)
	fmt.Println(i) // Expected output: 11110111

	// Count the number of set bits
	fmt.Println("Count:", i.Count()) // Expected output: Count: 7

	// Reset all bits to 0
	i.Reset()
	fmt.Println(i) // Expected output: 0

	// Initialize bitmask from a byte
	b := bitmask.FromByte(0b10101010)
	fmt.Println(b) // Expected output: 10101010

	// Bitwise AND with another bitmask
	a := bitmask.New(0b11110000)
	andResult := a.And(b)
	fmt.Println("AND:", andResult) // Expected output: AND: 10100000

	// Bitwise OR with another bitmask
	orResult := a.Or(b)
	fmt.Println("OR:", orResult) // Expected output: OR: 11111010

	// Bitwise XOR with another bitmask
	xorResult := a.Xor(b)
	fmt.Println("XOR:", xorResult) // Expected output: XOR: 01011010

	// Output:
	// 11111110
	// 11110110
	// 11110111
	// 2: true
	// []byte: [11110111]
	// 11110110
	// 11110111
	// Count: 7
	// 0
	// 10101010
	// AND: 10100000
	// OR: 11111010
	// XOR: 1011010
}
