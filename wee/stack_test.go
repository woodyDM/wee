package wee

import (
	"testing"
)

func Test_newStack(t *testing.T) {
	s := newStack()
	assertTrueMsg(len(s.elements) == 0, "ele empty", t)
	assertTrueMsg(s.elements != nil, "ele not nil", t)
}

func Test_stack_push(t *testing.T) {
	s := newStack()
	s.push(1)
	s.push(2)
	s.push(3)

	assertTrueMsg(len(s.elements) == 3, "len 3", t)
	assertTrueMsg(s.peek() == 3, "p 3", t)
	p := s.pop()
	assertTrueMsg(p == 3, "pop 3", t)
	p = s.pop()
	assertTrueMsg(p == 2, "pop 2", t)
	assertTrueMsg(len(s.elements) == 1, "len 1", t)
	assertTrueMsg(s.peek() == 1, "peek 1", t)
	s.push(10)
	assertTrueMsg(s.peek() == 10, "peek 10", t)
	assertTrueMsg(len(s.elements) == 2, "len 2", t)
	s.pop()
	p = s.pop()
	assertTrueMsg(p == 1, "pop 1", t)
	assertTrueMsg(len(s.elements) == 0, "empty stack", t)
}
