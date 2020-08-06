package wee

import "testing"

func TestCombine(t *testing.T) {
	method := "GET"

	p := splitAndCombine(method, "/1/2/3/")
	assertTrue(len(p) == 4, t)
	assertTrue(p[0] == "GET", t)
	assertTrue(p[1] == "1", t)
	assertTrue(p[2] == "2", t)
	assertTrue(p[3] == "3", t)

}

func TestCombineNoPrefix(t *testing.T) {
	method := "GET"

	p := splitAndCombine(method, "1/2/3/")
	assertTrue(len(p) == 4, t)
	assertTrue(p[0] == "GET", t)
	assertTrue(p[1] == "1", t)
	assertTrue(p[2] == "2", t)
	assertTrue(p[3] == "3", t)

}
func TestCombineNoPrefixAndSuffix(t *testing.T) {
	method := "GET"

	p := splitAndCombine(method, "1/2/3")
	assertTrue(len(p) == 4, t)
	assertTrue(p[0] == "GET", t)
	assertTrue(p[1] == "1", t)
	assertTrue(p[2] == "2", t)
	assertTrue(p[3] == "3", t)

}

func Test_should_insert_to_first_when_empty_list(t *testing.T) {
	n := getNewNode("ok")
	list := appendNormalNode(make([]*node, 0), n)

	assertTrue(len(list) == 1, t)
	assertTrue(list[0] == n, t)
}

func Test_should_insert_with_index_when_list(t *testing.T) {
	n := getNewNode("a")
	n2 := getNewNode("b")
	n3 := getNewNode("c")

	var list []*node
	list = appendNormalNode(make([]*node, 0), n2)
	list = appendNormalNode(list, n3)
	list = appendNormalNode(list, n)

	assertTrue(len(list) == 3, t)
	assertTrue(list[0] == n, t)
	assertTrue(list[1] == n2, t)
	assertTrue(list[2] == n3, t)
}

func Test_should_insert_with_index_when_wild(t *testing.T) {
	n := getNewNode("a")
	n2 := getNewNode("b")
	n3 := getNewNode("c")
	n4 := getNewNode(":id")
	n5 := getNewNode(":gid")

	var list = []*node{n4, n5}
	list = appendNormalNode(list, n2)
	list = appendNormalNode(list, n3)
	list = appendNormalNode(list, n)

	assertTrue(len(list) == 5, t)
	assertTrue(list[0] == n, t)
	assertTrue(list[1] == n2, t)
	assertTrue(list[2] == n3, t)
	assertTrue(list[3].wild, t)
	assertTrue(list[4].wild, t)
}

func Test_should_insert_with_index_when_wild_02(t *testing.T) {
	n := getNewNode("a")
	n2 := getNewNode("b")
	n3 := getNewNode("c")
	n4 := getNewNode(":id")

	var list = []*node{n, n4}
	list = appendNormalNode(list, n2)
	list = appendNormalNode(list, n3)

	assertTrue(len(list) == 4, t)
	assertTrue(list[0] == n, t)
	assertTrue(list[1] == n2, t)
	assertTrue(list[2] == n3, t)
	assertTrue(list[3].wild, t)
	assertTrue(list[3].pattern == ":id", t)
}

func getNewNode(pattern string) *node {
	return &node{
		pattern: pattern,
		wild:    isWild(pattern),
		value:   1,
	}
}

func TestRegister(t *testing.T) {
	tr := newTree()
	tr.register("GET", "/list", 1)
	tr.register("GET", "/list/:user", 2)
	tr.register("GET", "/list/:user/index", 22)
	tr.register("POST", "/list/:user", 6)

	//child is GET and post
	assertTrue(tr.root.childrenWildLen == 0, t)
	listRoot := tr.root.children[0]
	assertTrue(listRoot.pattern == "GET", t)
	assertTrue(listRoot.childrenWildLen == 0, t)
	assertTrue(len(listRoot.children) == 1, t)
	listNode := listRoot.children[0]
	assertTrue(listNode.pattern == "list", t)
	assertTrue(listNode.value == 1, t)
	assertTrue(len(listNode.children) == 1, t)
	assertTrue(listNode.childrenWildLen == 1, t)
	assertTrue(listNode.children[0].pattern == ":user", t)
	assertTrue(listNode.children[0].value == 2, t)
	assertTrue(len(listNode.children[0].children) == 1, t)
	assertTrue(listNode.children[0].children[0].value == 22, t)
}

func TestMatch01(t *testing.T) {
	tr := newTree()
	tr.register("GET", "/list/:date", 4)
	tr.register("GET", "/list/advise/index", 3)
	tr.register("GET", "/list/advise/index/p", 2)
	tr.register("GET", "/list", 100)
	tr.register("GET", "/list/:user/index", 1)
	tr.register("GET", "/list/content", 2)

	match := tr.match("GET", "/list/content")
	assertTrue(match.isMatches, t)
	match2 := tr.match("GET", "/list/advise/index")
	assertTrue(match2.isMatches, t)

	match3 := tr.match("GET", "/list/advise")
	assertTrue(match3.isMatches, t)
	match4 := tr.match("GET", "/list2/advise/w")
	assertTrue(!match4.isMatches, t)
	match5 := tr.match("POST", "/list/advise")
	assertTrue(!match5.isMatches, t)

	match6 := tr.match("GET", "/list/advise/index/page")
	assertTrue(!match6.isMatches, t)
	match7 := tr.match("GET", "/list/ok/index")
	assertTrue(match7.isMatches, t)
	assertTrue(match7.getMatchedPath() == "/GET/list/:user/index", t)
	match8 := tr.match("GET", "/list/advise/index")
	assertTrue(match8.isMatches, t)
	assertTrue(match8.getMatchedPath() == "/GET/list/advise/index", t)
	match9 := tr.match("GET", "/list/ok/index/p")
	assertTrue(!match9.isMatches, t)

}

func TestPathVariableMap(t *testing.T) {
	tr := newTree()
	tr.register("GET", "/list/:userId/:aId/index", 100)
	match := tr.match("GET", "/list/100/group/index")
	assertTrue(match.isMatches, t)
	assertTrueMsg(match.getMatchedPath() == "/GET/list/:userId/:aId/index", "should mathc path", t)
	assertTrueMsg(len(match.pathVariable) == 2, "length should be 2", t)
	assertTrueMsg(match.pathVariable["userId"] == "100", "100 eq", t)
	assertTrueMsg(match.pathVariable["aId"] == "group", "aid group", t)
}

func TestMatch02(t *testing.T) {
	tr := newTree()
	tr.register("GET", "/list/page", 4)

	m1 := tr.match("GET", "/list/page")
	assertTrue(m1.isMatches, t)
	assertTrue(m1.getMatchedPath() == "/GET/list/page", t)
	m2 := tr.match("GET", "/list/1")
	assertTrue(!m2.isMatches, t)
	m3 := tr.match("GET", "/list/page/2")
	assertTrue(!m3.isMatches, t)

}

func TestMatch03(t *testing.T) {
	tr := newTree()
	tr.register("GET", "/list/:id", 4)

	m1 := tr.match("GET", "/list")
	assertTrue(!m1.isMatches, t)
	m2 := tr.match("GET", "/list/1")
	assertTrue(m2.isMatches, t)
	assertTrue(m2.getMatchedPath() == "/GET/list/:id", t)
	m3 := tr.match("GET", "/list/id/1")
	assertTrue(!m3.isMatches, t)

}

func TestMatch04(t *testing.T) {
	tr := newTree()
	tr.register("GET", "/list/:id/page", 4)

	m1 := tr.match("GET", "/list/id")
	assertTrue(!m1.isMatches, t)
	m2 := tr.match("GET", "/list/1/page")
	assertTrue(m2.isMatches, t)
	matchPath := m2.getMatchedPath()
	assertTrue(matchPath == "/GET/list/:id/page", t)
	m3 := tr.match("GET", "/list/id/1")
	assertTrue(!m3.isMatches, t)

}

func TestRegisterForWild(t *testing.T) {
	tr := newTree()
	tr.register("GET", "/list", 100)
	tr.register("GET", "/list/content", 2)
	tr.register("GET", "/list/advise/index", 3)
	tr.register("GET", "/list/:user", 1)
	tr.register("GET", "/list/:userName/profile", 5)
	validateTheTree(tr, t)
}

func TestRegisterForWild02(t *testing.T) {
	tr := newTree()
	tr.register("GET", "/list/:user", 1)
	tr.register("GET", "/list/content", 2)
	tr.register("GET", "/list/:userName/profile", 5)
	tr.register("GET", "/list/advise/index", 3)
	tr.register("GET", "/list", 100)

	validateTheTree(tr, t)
}

func validateTheTree(tr *tree, t *testing.T) {
	methodNode := tr.root.children[0]
	assertTrue(methodNode.value == nil, t)
	assertTrue(methodNode.pattern == "GET", t)
	assertTrue(methodNode.childrenWildLen == 0, t)
	assertTrue(len(methodNode.children) == 1, t)
	listNode := methodNode.children[0]
	assertTrue(listNode.childrenWildLen == 2, t)
	assertTrue(listNode.value == 100, t)
	assertTrue(listNode.pattern == "list", t)
	assertTrue(len(listNode.children) == 4, t)
	assertTrue(listNode.children[0].value == nil, t)
	assertTrue(listNode.children[0].pattern == "advise", t)
	assertTrue(listNode.children[1].value == 2, t)
	assertTrue(listNode.children[1].pattern == "content", t)
	assertTrue(listNode.children[2].wild, t)
	assertTrue(listNode.children[3].wild, t)
	assertTrue(listNode.children[2].pattern != listNode.children[3].pattern, t)
	assertTrue(listNode.children[2].pattern == ":user" || listNode.children[2].pattern == ":userName", t)
	if listNode.children[2].pattern == ":user" {
		assertTrue(listNode.children[2].children == nil, t)
		assertTrue(len(listNode.children[3].children) == 1, t)
		assertTrue(listNode.children[3].children[0].value == 5, t)
	} else {
		assertTrue(listNode.children[3].children == nil, t)
		assertTrue(len(listNode.children[2].children) == 1, t)
		assertTrue(listNode.children[2].children[0].value == 5, t)
	}
	assertTrue(len(listNode.children[0].children) == 1, t)
	assertTrue(listNode.children[1].children == nil, t)
	assertTrue(listNode.children[2].children == nil, t)

	adIndexNode := listNode.children[0].children[0]
	assertTrue(adIndexNode.value == 3, t)
	assertTrue(adIndexNode.children == nil, t)
	assertTrue(adIndexNode.pattern == "index", t)
	assertTrue(adIndexNode.childrenWildLen == 0, t)
}

func TestSearch01(t *testing.T) {
	nodes := []*node{newNode("a")}
	r := search("a", 0, nodes)
	assertTrue(len(r) == 1, t)
	r2 := search("b", 0, nodes)
	assertTrue(len(r2) == 0, t)
}

func TestSearch02(t *testing.T) {
	nodes := []*node{newNode("a"), newNode("c")}
	r := search("a", 0, nodes)
	assertTrue(len(r) == 1, t)
	r2 := search("b", 0, nodes)
	assertTrue(len(r2) == 0, t)
}

func TestSearch03(t *testing.T) {
	nodes := []*node{newNode("a"), newNode("c"), newNode(":d")}
	r := search("a", 1, nodes)
	assertTrue(len(r) == 2, t)
	r2 := search("b", 1, nodes)
	assertTrue(len(r2) == 1, t)
}

func TestSearch04(t *testing.T) {
	nodes := []*node{newNode("a"), newNode("c"), newNode(":d"), newNode(":od")}
	r := search("a", 2, nodes)
	assertTrue(len(r) == 3, t)
	r2 := search(":od", 2, nodes)
	assertTrue(len(r2) == 2, t)
}

func TestSearch05(t *testing.T) {
	nodes := []*node{newNode(":a"), newNode(":c"), newNode(":d"), newNode(":od")}
	r := search("a", 4, nodes)
	assertTrue(len(r) == 4, t)
	r2 := search("w", 4, nodes)
	assertTrue(len(r2) == 4, t)
}

func TestBinarySearch01(t *testing.T) {
	nodes := []*node{newNode("a")}
	r := binarySearch("a", nodes)
	assertTrue(r != nil && r.pattern == "a", t)
	r2 := binarySearch("b", nodes)
	assertTrue(r2 == nil, t)
}

func TestBinarySearch02(t *testing.T) {
	nodes := []*node{newNode("a"), newNode("b")}
	r := binarySearch("a", nodes)
	assertTrue(r != nil && r.pattern == "a", t)
	r2 := binarySearch("b", nodes)
	assertTrue(r2 != nil && r2.pattern == "b", t)
	r3 := binarySearch("e", nodes)
	assertTrue(r3 == nil, t)
}

func TestBinarySearch03(t *testing.T) {
	nodes := []*node{newNode("a"), newNode("c"), newNode("d")}
	r := binarySearch("a", nodes)
	assertTrue(r != nil, t)
	r2 := binarySearch("c", nodes)
	assertTrue(r2 != nil, t)
	r3 := binarySearch("d", nodes)
	assertTrue(r3 != nil, t)
	r4 := binarySearch("b", nodes)
	assertTrue(r4 == nil, t)
}

func newNode(p string) *node {
	return &node{
		pattern: p,
		wild:    isWild(p),
	}
}

func assertTrueMsg(condition bool, msg string, t *testing.T) {
	if !condition {
		t.Fatal(msg)
	}
}

func assertMid(ware MidWare, name string, t *testing.T) {
	if m, ok := ware.(*MockMidWare); ok && m.name == name {
		return
	} else {
		t.Fatal("not match " + name)
	}
}

func assertTrue(condition bool, t *testing.T) {
	assertTrueMsg(condition, "no message", t)
}
