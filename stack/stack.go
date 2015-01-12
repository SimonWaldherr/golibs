package stack

func Lifo() *Stack {
	return &Stack{}
}

func Ring() *Rings {
	return &Rings{}
}

type Stack struct {
	nodes []interface{}
	count int
}

type Rings struct {
	nodes  []string
	count  int
	xcount int
	size   int
}

func (s *Stack) Push(n interface{}) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

func (s *Stack) Pop() interface{} {
	if s.count == 0 {
		return ""
	}
	s.count--
	return s.nodes[s.count]
}

func (s *Stack) Len() int {
	return s.count
}

func (r *Rings) SetSize(i int) {
	if i > r.xcount {
		r.size = i
	}
}

func (r *Rings) Init() {
	r.size = 1
	r.count = -1
	r.xcount = -1
}

func (r *Rings) GetSize() int {
	return r.size
}

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

func (r *Rings) Pos() (int, int) {
	return r.count, r.xcount
}
