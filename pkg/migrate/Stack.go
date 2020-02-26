package migrate

// Stack represents a collection of SQL queries
type Stack struct {
	stack []string
}

// NewStack creates a new Stack with an empty slice of strings for queries
func NewStack() Stack {
	return Stack{
		stack: make([]string, 0),
	}
}

// Slice returns the slice of strings from the Stack
func (s *Stack) Slice() []string {
	return s.stack
}

// Empty returns true if the stack has no queries in it
func (s *Stack) Empty() bool {
	return len(s.stack) == 0
}

// Len returns the list of queries in the stack
func (s *Stack) Len() int {
	return len(s.stack)
}

// Push appends a query to the Stack
func (s *Stack) Push(str string) {
	s.stack = append(s.stack, str)
}

// Pop removes the last item from the stack
func (s *Stack) Pop() {
	s.stack = s.stack[0 : len(s.stack)-1]
}

// Last returns the last query from the stack
func (s *Stack) Last() string {
	return s.stack[len(s.stack)-1]
}
