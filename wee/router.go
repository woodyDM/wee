package wee

type TreeRouter struct {
	tr *tree
}

func newTreeRouter() *TreeRouter {
	return &TreeRouter{tr: newTree()}
}

func (t *TreeRouter) Register(method string, pattern string, chain *MidWareChain) {
	t.tr.register(method, pattern, chain)
}

func (t *TreeRouter) match(method string, pattern string) MatchedRoute {
	return t.tr.match(method, pattern)
}

func (m *matchResult) matches() bool {
	return m.isMatches
}

func (m *matchResult) getPathVariable() map[string]string {
	return m.pathVariable
}

func (m *matchResult) chain() *MidWareChain {
	level := m.level[len(m.parts)-1]
	h := level.target[level.pos].value
	if chain, ok := h.(*MidWareChain); ok {
		return chain
	} else {
		panic("can't convert to chain, check your register!")
	}
}
