package migrate

type Stack struct {
	stack []string
}

func NewStack() Stack {
	return Stack{
		stack: make([]string, 0),
	}
}

func (s *Stack) Slice() []string {
	return s.stack
}

func (s *Stack) Empty() bool {
	return len(s.stack) == 0
}

func (s *Stack) Len() int {
	return len(s.stack)
}

func (s *Stack) Push(str string) {
	s.stack = append(s.stack, str)
}

func (s *Stack) Pop() {
	s.stack = s.stack[0 : len(s.stack)-1]
}

func (s *Stack) Last() string {
	return s.stack[len(s.stack)-1]
}
