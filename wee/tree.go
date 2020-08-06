package wee

import (
	"bytes"
	"strings"
)

const (
	urlSep            = "/"
	wildPatternPrefix = ":"
)

/**
root is always "/"
seconds level is http method
wildcard is always at end of children
non-wildcard child is sorted by string
children looks like this:  ["1","2","3",":user",":id"]
*/
type tree struct {
	root *node
}

type node struct {
	pattern         string
	wild            bool
	children        []*node
	value           interface{}
	childrenWildLen int
}

type matchResult struct {
	isMatches    bool
	method       string
	path         string
	parts        []string
	level        []*searchLevel
	pathVariable map[string]string
}

type searchLevel struct {
	target []*node
	pos    int
}

func newTree() *tree {
	return &tree{
		root: &node{
			pattern: urlSep,
		},
	}
}

func (t *tree) register(method string, path string, v interface{}) {
	p := splitAndCombine(method, path)
	var current = t.root
	max := len(p)
	for i := 0; i < max; i++ {
		newNode := createNewNode(i, max, p[i], v)
		if current.children == nil {
			current.children = make([]*node, 1)
			current.children[0] = newNode
			if newNode.wild {
				current.childrenWildLen += 1
			}
			current = newNode
		} else {
			//if not in children then insert to children
			next := appendChildrenIfNeed(current, newNode, i == max-1)
			checkDuplicateWildChild(current)
			current = next
		}
	}
}

func checkDuplicateWildChild(n *node) {
	var counter = 0
	for _, cn := range n.children {
		if cn.wild && cn.value != nil {
			counter++
		}
		if counter > 1 {
			panic("Duplicated child wild card mapping for node :[" + n.pattern + "]")
		}
	}
}

func (m *matchResult) generatePathVariableMap() {
	result := make(map[string]string)
	if m.isMatches {
		for i, n := range m.level {
			node := n.target[n.pos]
			if node.wild {
				key := strings.TrimPrefix(node.pattern, wildPatternPrefix)
				result[key] = m.parts[i]
			}
		}
	}
	m.pathVariable = result
}

func (m *matchResult) getMatchedPath() string {
	if !m.isMatches {
		return ""
	}
	b := new(bytes.Buffer)
	for _, n := range m.level {
		b.WriteString(urlSep)
		b.WriteString(n.target[n.pos].pattern)
	}
	return b.String()
}

func (t *tree) match(method string, path string) *matchResult {
	parts := splitAndCombine(method, path)
	result := &matchResult{
		method: method,
		path:   path,
		parts:  parts,
		level:  make([]*searchLevel, len(parts)),
	}
	doSearch(parts, 0, t.root, result)
	result.generatePathVariableMap()
	return result

}

/**
to search url in the tree
parts: [GET, list , index]
i:  0  root /
    1  GET /
    2  list /
    3  index
*/
func doSearch(parts []string, i int, n *node, result *matchResult) {
	maxDepth := len(parts)
	if result.isMatches {
		return //no longer search if ok
	}
	if i < maxDepth {
		target := search(parts[i], n.childrenWildLen, n.children)
		if result.level[i] == nil {
			result.level[i] = &searchLevel{
				target: target,
				pos:    0,
			}
		}
	} else if i == maxDepth {
		//reach final
		if n.matches(parts[maxDepth-1]) && n.value != nil {
			result.isMatches = true
		}
		return
	} else {
		//exceed max depth
		return
	}
	//at this position i < maxDepth
	thisLevel := result.level[i]
	for ; thisLevel.pos < len(thisLevel.target); thisLevel.pos++ {
		doSearch(parts, i+1, thisLevel.target[thisLevel.pos], result)
		if result.isMatches {
			break
		} else {
			//clear bottom level search data
			if i+1 < maxDepth {
				result.level[i+1] = nil
			}
		}
	}
}

func (n *node) matches(part string) bool {
	return n.pattern == part || n.wild

}

/**
to find p in some node.children (exact match or wildCard match)
*/
func search(p string, childWild int, nodes []*node) []*node {
	if len(nodes) == 0 {
		return nil
	}
	exact := binarySearch(p, nodes[:len(nodes)-childWild])
	if exact == nil {
		return nodes[len(nodes)-childWild:]
	} else {
		result := make([]*node, childWild+1)
		result[0] = exact
		copy(result[1:], nodes[len(nodes)-childWild:])
		return result
	}
}

func binarySearch(p string, nodes []*node) *node {
	if len(nodes) == 0 {
		return nil
	}
	left := 0
	right := len(nodes) - 1
	for mid := (left + right) / 2; left <= right; mid = (left + right) / 2 {
		n := nodes[mid]
		if n.pattern == p {
			return n
		} else if n.pattern > p {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return nil
}

func isWild(pattern string) bool {
	return strings.HasPrefix(pattern, wildPatternPrefix)
}

func createNewNode(depth int, max int, pattern string, v interface{}) *node {
	var nv interface{}
	if depth == max-1 {
		nv = v
	}
	return &node{
		pattern: pattern,
		value:   nv,
		wild:    isWild(pattern),
	}
}

/**
insert n to list if not exist or else:try update the existing node.
return nextCurrentNode
*/
func appendChildrenIfNeed(current *node, n *node, last bool) *node {
	list := current.children
	if ex := exist(list, n, last); ex == nil {
		if n.wild {
			//wild pattern node use append
			current.childrenWildLen += 1
			current.children = append(list, n)
		} else {
			//normal pattern should insert by string order
			current.children = appendNormalNode(list, n)
		}
		return n
	} else {
		//replace the existing node with n.value
		//in the case : first register /path/node/1  with V1 then register /path/node with V2 in this condition ,the node value should be V2
		if ex.value == nil {
			ex.value = n.value
		}
		return ex
	}
}

/**
find existing node in list
*/
func exist(list []*node, n *node, last bool) *node {
	for i := 0; i < len(list); i++ {
		if list[i].pattern == n.pattern {
			if last && list[i].value != nil {
				panic("duplicated mapping for pattern: " + n.pattern)
			}
			return list[i]
		}
	}
	return nil
}

func appendNormalNode(list []*node, n *node) []*node {
	var end = getWildBound(list)
	var i = 0
	for ; i < end; i++ {
		if list[i].pattern > n.pattern {
			break
		}
	}
	//insert to position i at list
	result := make([]*node, len(list)+1)
	copy(result[:i], list[:i])
	result[i] = n
	copy(result[i+1:], list[i:])
	return result
}

/**
return min index of wild pattern child,
if no wild pattern, return length of list.
*/
func getWildBound(list []*node) int {
	for i := 0; i < len(list); i++ {
		if list[i].wild {
			return i
		}
	}
	return len(list)
}

func splitAndCombine(method string, path string) []string {
	path = strings.TrimPrefix(path, urlSep)
	path = strings.TrimSuffix(path, urlSep)
	parts := strings.Split(path, urlSep)

	result := make([]string, 1+len(parts))
	result[0] = method
	copy(result[1:], parts)
	return result
}
