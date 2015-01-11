package stack

func Lifo() *Stack {
	return &Stack{}
}

type Stack struct {
	nodes []string
	count int
}

func (s *Stack) Push(n string) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

func (s *Stack) Pop() string {
	if s.count == 0 {
		return ""
	}
	s.count--
	return s.nodes[s.count]
}
