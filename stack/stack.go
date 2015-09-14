// stack implements a stack with "last in, first out" functionality.
// it also provides a ring memory type which overrides itself after n write ops.
package stack

type Stype int

const (
	NIL Stype = iota
	LiFo
	FiFo
)

// struct Stack contains nodes as slice of interfaces
// and a counter for the current position.
type Stack struct {
	nodes []interface{}
	count int
	stype Stype
}

// Lifo returns a pointer to a new stack.
func Lifo() *Stack {
	return &Stack{
		stype: LiFo,
	}
}

func Fifo() *Stack {
	return &Stack{
		stype: FiFo,
	}
}

func (s *Stack) Unset() {
	*s = Stack{
		nodes: []interface{}{},
		count: 0,
		stype: NIL,
	}
}

func (s *Stack) ToFifo() *Stack {
	var x *Stack
	array := Fifo()
	if s.stype == FiFo {
		x = Fifo()
		*x = *s
		for s.Len() > 0 {
			array.Push(s.Pop())
		}
	}
	if s.stype == LiFo {
		x = Lifo()
		*x = *s
		helper := Lifo()
		for s.Len() > 0 {
			helper.Push(s.Pop())
		}
		for helper.Len() > 0 {
			array.Push(helper.Pop())
		}
	}
	*s = *x
	return array
}

func (s *Stack) ToLifo() *Stack {
	var x *Stack
	array := Lifo()
	if s.stype == FiFo {
		x = Fifo()
		*x = *s
		for s.Len() > 0 {
			array.Push(s.Pop())
		}
	}
	if s.stype == LiFo {
		x = Lifo()
		*x = *s
		helper := Lifo()
		for s.Len() > 0 {
			helper.Push(s.Pop())
		}
		for helper.Len() > 0 {
			array.Push(helper.Pop())
		}
	}
	*s = *x
	return array
}

func (s *Stack) Val() []interface{} {
	var a *Stack
	var r []interface{}

	if s.stype == FiFo {
		a = s.ToFifo()
	} else if s.stype == LiFo {
		a = s.ToLifo()
	}

	for s.Len() > 0 {
		r = append(r, s.Pop())
	}
	*s = *a
	return r
}

// Push adds a value to the Stack
func (s *Stack) Push(n interface{}) {
	if s.stype == LiFo {
		s.nodes = append(s.nodes[:s.count], n)
		s.count++
	} else if s.stype == FiFo {
		s.nodes = append(s.nodes, n)
	}
}

func (s *Stack) Add(n interface{}) {
	s.Push(n)
}

// Pop returns the last added value and decrease the position counter.
func (s *Stack) Pop() interface{} {
	if s.stype == LiFo {
		if s.count == 0 {
			return ""
		}
		s.count--
		return s.nodes[s.count]
	}
	if s.stype == FiFo {
		if s.count == len(s.nodes) {
			return ""
		}
		s.count++
		return s.nodes[s.count-1]
	}
	return ""
}

func (s *Stack) Get() interface{} {
	return s.Pop()
}

// Len returns the current position in the Stack.
func (s *Stack) Len() int {
	if s.stype == LiFo {
		return s.count
	}
	if s.stype == FiFo {
		return len(s.nodes) - s.count
	}
	return -1
}

func (s *Stack) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

// Ring returns a pointer to a new ring.
func Ring() *Rings {
	return &Rings{}
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

// Pos returns the current position and the
// number of overall added values
func (r *Rings) Pos() (int, int) {
	return r.count, r.xcount
}
