package stack

func Lifo() *Stack {
	return &Stack{}
}

type Stack struct {
	nodes []interface{}
	count int
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
