// Package stack provides a concurrent Stack (LIFO/FIFO) container and a ring buffer.
package stack

import (
	"bytes"
	"encoding/gob"
	"sync"
)

// Stype indicates the stack ordering.
type Stype int

const (
	// NIL indicates an uninitialized Stack state.
	NIL Stype = iota
	// LiFo represents Last-In-First-Out ordering.
	LiFo
	// FiFo represents First-In-First-Out ordering.
	FiFo
)

// Stack holds values in LIFO or FIFO order depending on its constructor.
// Its methods are safe for concurrent use. Stored values remain caller-owned.
// A Stack must not be copied after first use.
type Stack struct {
	nodes []interface{}
	count int
	stype Stype
	mutex sync.RWMutex
}

// Lifo returns a pointer to a new LIFO stack.
func Lifo() *Stack {
	return &Stack{
		stype: LiFo,
	}
}

// Fifo returns a pointer to a new FIFO stack.
func Fifo() *Stack {
	return &Stack{
		stype: FiFo,
	}
}

// Stack2 is a type alias for Stack.
type Stack2 = Stack

// ToFifo is a package-level convenience function that converts a Stack to FIFO order.
func ToFifo(s *Stack) *Stack { return s.ToFifo() }

// ToFifoFromFifo is a package-level convenience function that creates a new FIFO copy
// of a FIFO Stack, restoring iteration order. Equivalent to ToFifo for FIFO stacks.
func ToFifoFromFifo(s *Stack) *Stack { return s.ToFifo() }

// ToLifo is a package-level convenience function that converts a Stack to LIFO order.
func ToLifo(s *Stack) *Stack { return s.ToLifo() }

// ToLifoFromLifo is a package-level convenience function that creates a new LIFO copy
// of a LIFO Stack, restoring iteration order. Equivalent to ToLifo for LIFO stacks.
func ToLifoFromLifo(s *Stack) *Stack { return s.ToLifo() }

// Unset resets the stack to its zero/uninitialized state.
func (s *Stack) Unset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.nodes = nil
	s.count = 0
	s.stype = NIL
}

// ToFifo returns an independent FIFO copy of the remaining values in insertion order.
// The receiver retains its contents and ordering.
func (s *Stack) ToFifo() *Stack {
	nodes, _ := s.snapshot()
	return &Stack{nodes: nodes, stype: FiFo}
}

// ToLifo returns an independent LIFO copy of the remaining values in insertion order.
// The receiver retains its contents and ordering.
func (s *Stack) ToLifo() *Stack {
	nodes, _ := s.snapshot()
	return &Stack{nodes: nodes, count: len(nodes), stype: LiFo}
}

// Val returns a snapshot of the remaining values in pop order.
// Modifying the returned slice does not change the stack.
func (s *Stack) Val() []interface{} {
	nodes, stype := s.snapshot()
	if stype == LiFo {
		for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
			nodes[i], nodes[j] = nodes[j], nodes[i]
		}
	}
	return nodes
}

// snapshot copies the available values in insertion order.
func (s *Stack) snapshot() ([]interface{}, Stype) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var nodes []interface{}
	if s.stype == LiFo {
		nodes = s.nodes[:s.count]
	} else if s.stype == FiFo {
		nodes = s.nodes[s.count:]
	}
	return append([]interface{}(nil), nodes...), s.stype
}

// Push adds a value to the Stack according to its ordering.
func (s *Stack) Push(n interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.stype == LiFo {
		s.nodes = append(s.nodes[:s.count], n)
		s.count++
	} else if s.stype == FiFo {
		s.nodes = append(s.nodes, n)
	}
}

// Add is an alias for Push.
func (s *Stack) Add(n interface{}) {
	s.Push(n)
}

// Pop removes and returns the next element according to the stack's
// ordering. When the stack is empty a zero-value (empty string) is
// returned (to preserve original behavior).
func (s *Stack) Pop() interface{} {
	if value, ok := s.TryPop(); ok {
		return value
	}
	return ""
}

// TryPop removes the next value and reports whether one was available.
// It distinguishes an empty stack from a stored empty string or nil value.
func (s *Stack) TryPop() (interface{}, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.stype == LiFo {
		if s.count == 0 {
			return nil, false
		}
		s.count--
		value := s.nodes[s.count]
		s.nodes[s.count] = nil
		s.nodes = s.nodes[:s.count]
		return value, true
	}
	if s.stype == FiFo {
		if s.count == len(s.nodes) {
			return nil, false
		}
		value := s.nodes[s.count]
		s.nodes[s.count] = nil
		s.count++
		if s.count == len(s.nodes) {
			s.nodes = nil
			s.count = 0
		} else if s.count >= 1024 && s.count >= len(s.nodes)/2 {
			// Bound the consumed prefix of queues that never become empty.
			s.nodes = append([]interface{}(nil), s.nodes[s.count:]...)
			s.count = 0
		}
		return value, true
	}
	return nil, false
}

// Peek returns the next value without removing it, reporting whether one is available.
func (s *Stack) Peek() (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.stype == LiFo && s.count > 0 {
		return s.nodes[s.count-1], true
	}
	if s.stype == FiFo && s.count < len(s.nodes) {
		return s.nodes[s.count], true
	}
	return nil, false
}

// Get is an alias for Pop.
func (s *Stack) Get() interface{} {
	return s.Pop()
}

// Len returns the number of elements available according to the
// current ordering. When uninitialized the function returns -1.
func (s *Stack) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.stype == LiFo {
		return s.count
	}
	if s.stype == FiFo {
		return len(s.nodes) - s.count
	}
	return -1
}

// IsEmpty reports whether the stack contains no elements.
func (s *Stack) IsEmpty() bool {
	return s.Len() == 0
}

// Marshal encodes a snapshot of the remaining values in insertion order
// using the gob []interface{} format. The stack ordering is not encoded.
func (s *Stack) Marshal() ([]byte, error) {
	nodes, _ := s.snapshot()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(nodes); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal replaces the stack contents with a gob-encoded []interface{}.
// The receiver retains its ordering. A decoding error leaves the stack unchanged.
func (s *Stack) Unmarshal(data []byte) error {
	var nodes []interface{}
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&nodes); err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.nodes = nodes
	s.count = 0
	if s.stype == LiFo {
		s.count = len(nodes)
	}
	return nil
}

// Iter returns a buffered, closed channel containing a snapshot of the
// remaining values in insertion order. Receiving does not consume stack values.
func (s *Stack) Iter() <-chan interface{} {
	nodes, _ := s.snapshot()
	c := make(chan interface{}, len(nodes))
	for _, value := range nodes {
		c <- value
	}
	close(c)
	return c
}

// Ring returns a pointer to a new ring buffer.
func Ring() *Rings {
	return &Rings{}
}

// Rings is a simple ring buffer storing strings. The buffer will
// overwrite older entries once the buffer size is reached.
type Rings struct {
	nodes  []string
	count  int
	xcount int
	size   int
}

// SetSize sets the maximum capacity of the ring. The value is applied
// only when greater than the total number of values ever added.
func (r *Rings) SetSize(i int) {
	if i > r.xcount {
		r.size = i
	}
}

// Init prepares the ring for use with an initial size and resets
// internal counters.
func (r *Rings) Init(i int) {
	r.size = i
	r.count = -1
	r.xcount = -1
}

// GetSize returns the configured maximum size of the ring.
func (r *Rings) GetSize() int {
	return r.size
}

// Push adds a string to the ring and returns the position where it
// was written.
func (r *Rings) Push(n string) int {
	r.count++
	r.xcount++
	if r.count > r.size {
		r.count = 0
	}
	if r.xcount > r.count {
		r.nodes[r.count] = n
	} else {
		r.nodes = append(r.nodes[:r.count], n)
	}

	return r.count
}

// Get returns a slice of strings from `from` up to the current
// position. The returned slice has a length equal to the ring size.
func (r *Rings) Get(from int) []string {
	ret := make([]string, r.size)
	var i int
	for from != r.count {
		ret[i] = string(r.nodes[from])
		from++
		i++
		if from > r.size {
			from = 0
		}
	}
	return ret
}

// Pos returns the current write position and the total count of
// values written to the ring.
func (r *Rings) Pos() (int, int) {
	return r.count, r.xcount
}
