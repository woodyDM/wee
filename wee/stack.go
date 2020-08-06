package wee

type stack struct {
	elements []interface{}
}

func newStack() *stack {
	return &stack{
		elements: make([]interface{}, 0),
	}
}

func (s *stack) push(o interface{}) {
	s.elements = append(s.elements, o)
}

func (s *stack) isEmpty() bool {
	return len(s.elements) == 0

}

func (s *stack) forEach(consumer func(i interface{})) {
	for _, it := range s.elements {
		consumer(it)
	}
}
func (s *stack) peek() interface{} {
	l := len(s.elements)
	if l == 0 {
		panic("Can't peek on an empty stack")
	}
	return s.elements[l-1]
}

func (s *stack) pop() interface{} {
	l := len(s.elements)
	if l == 0 {
		panic("Can't pop on an empty stack")
	}
	r := s.elements[l-1]
	s.elements[l-1] = nil
	s.elements = s.elements[:l-1]
	return r
}
