// Package deque provides a generic growable double-ended queue.
package deque

// Deque stores values in a ring buffer. Its zero value is ready to use.
// Push and pop operations take amortized constant time. A Deque is not safe for
// concurrent mutation and must not be copied after first use; use Clone instead.
// Elements themselves are copied shallowly.
type Deque[T any] struct {
	buffer []T
	head   int
	size   int
}

// Len returns the number of stored elements.
func (d *Deque[T]) Len() int { return d.size }

// Cap returns the current ring capacity.
func (d *Deque[T]) Cap() int { return len(d.buffer) }

// PushFront adds value before the first element.
func (d *Deque[T]) PushFront(value T) {
	d.reserve()
	d.head = (d.head - 1) & (len(d.buffer) - 1)
	d.buffer[d.head] = value
	d.size++
}

// PushBack adds value after the last element.
func (d *Deque[T]) PushBack(value T) {
	d.reserve()
	d.buffer[(d.head+d.size)&(len(d.buffer)-1)] = value
	d.size++
}

// Front returns the first value and whether the deque is nonempty.
func (d *Deque[T]) Front() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	return d.buffer[d.head], true
}

// Back returns the last value and whether the deque is nonempty.
func (d *Deque[T]) Back() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	return d.buffer[(d.head+d.size-1)&(len(d.buffer)-1)], true
}

// PopFront removes the first element, releasing its storage reference.
// An empty deque returns the zero value and false.
func (d *Deque[T]) PopFront() (T, bool) {
	value, ok := d.Front()
	if !ok {
		return value, false
	}
	var zero T
	d.buffer[d.head] = zero
	d.head = (d.head + 1) & (len(d.buffer) - 1)
	d.size--
	if d.size == 0 {
		d.head = 0
	}
	return value, true
}

// PopBack removes the last element, releasing its storage reference.
// An empty deque returns the zero value and false.
func (d *Deque[T]) PopBack() (T, bool) {
	value, ok := d.Back()
	if !ok {
		return value, false
	}
	d.size--
	var zero T
	d.buffer[(d.head+d.size)&(len(d.buffer)-1)] = zero
	if d.size == 0 {
		d.head = 0
	}
	return value, true
}

// Clear removes all values and their references, retaining the current capacity.
func (d *Deque[T]) Clear() {
	first := min(d.size, len(d.buffer)-d.head)
	clear(d.buffer[d.head : d.head+first])
	clear(d.buffer[:d.size-first])
	d.head, d.size = 0, 0
}

// Values returns an independent slice in front-to-back order, or nil when empty.
func (d *Deque[T]) Values() []T {
	if d.size == 0 {
		return nil
	}
	values := make([]T, d.size)
	d.copyTo(values)
	return values
}

// Clone returns an independent deque with the same values. Element references
// are copied shallowly. Excess source capacity is not retained, and empty deques
// clone to the zero value.
func (d *Deque[T]) Clone() Deque[T] {
	var result Deque[T]
	if d.size > 0 {
		capacity := 8
		for capacity < d.size {
			capacity *= 2
		}
		result.buffer = make([]T, capacity)
		d.copyTo(result.buffer)
		result.size = d.size
	}
	return result
}

func (d *Deque[T]) copyTo(target []T) {
	first := min(d.size, len(d.buffer)-d.head)
	copy(target, d.buffer[d.head:d.head+first])
	copy(target[first:], d.buffer[:d.size-first])
}

func (d *Deque[T]) reserve() {
	if d.size < len(d.buffer) {
		return
	}
	capacity := 8
	if len(d.buffer) > 0 {
		if len(d.buffer) > int(^uint(0)>>1)/2 {
			panic("deque: capacity overflow")
		}
		capacity = len(d.buffer) * 2
	}
	buffer := make([]T, capacity)
	d.copyTo(buffer)
	d.buffer, d.head = buffer, 0
}
