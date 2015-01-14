package stack

// Lifo returns a pointer to a new stack.
func Lifo() *Stack {
	return &Stack{}
}

// Ring returns a pointer to a new ring.
func Ring() *Rings {
	return &Rings{}
}

// struct Stack contains nodes as slice of interfaces
// and a counter for the current position.
type Stack struct {
	nodes []interface{}
	count int
}

// struct Rings contains nodes as slice of strings,
// count for the current ring position,
// xcount for the total amount of added entries
// and size for the maximum size of the "Ring".
type Rings struct {
	nodes  []string
	count  int
	xcount int
	size   int
}

// Push adds a value to the Stack
func (s *Stack) Push(n interface{}) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

// Pop returns the last added value and decrease the position counter.
func (s *Stack) Pop() interface{} {
	if s.count == 0 {
		return ""
	}
	s.count--
	return s.nodes[s.count]
}

// Len returns the current position in the Stack.
func (s *Stack) Len() int {
	return s.count
}

// SetSize sets the maximum size of the Ring,
// this size must be greater then the current counter.
func (r *Rings) SetSize(i int) {
	if i > r.xcount {
		r.size = i
	}
}

// Init sets the default values of the Ring.
func (r *Rings) Init(i int) {
	r.size = i
	r.count = -1
	r.xcount = -1
}

// GetSize returns the max size of the Ring.
func (r *Rings) GetSize() int {
	return r.size
}

// Push adds a string to the Ring and returns it position
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

// Get returns a slice of strings from the given
// to the current position
func (r *Rings) Get(from int) []string {
	ret := make([]string, r.size)
	var i int = 0
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

// Pos returns the current position and the
// number of overall added values
func (r *Rings) Pos() (int, int) {
	return r.count, r.xcount
}
