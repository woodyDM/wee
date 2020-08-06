package wee

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	router Router
	port   int
	group  *RegistryGroup //default group
}

type Handler func(ctx *Context)

type MatchedRoute interface {
	matches() bool
	getPathVariable() map[string]string
	chain() *MidWareChain
}

type Router interface {
	Register(method string, patter string, chain *MidWareChain)
	match(method string, pattern string) MatchedRoute
}

type PathHandler struct {
	method  string
	pattern string
	h       Handler
}

type RegistryGroup struct {
	midWares     []MidWare
	subGroup     []*RegistryGroup
	pathHandlers []*PathHandler
	pathPrefix   string
}

func (h Handler) handle(ctx *Context) {
	h(ctx)
}

func (r *RegistryGroup) Use(ware MidWare) {
	if r.midWares == nil {
		r.midWares = make([]MidWare, 0)
	}
	r.midWares = append(r.midWares, ware)
}

func (r *RegistryGroup) Register(method string, path string, h Handler) {
	if path == "" || path == urlSep {
		path = ""
	} else {
		path = normalize(path)
	}
	if r.pathHandlers == nil {
		r.pathHandlers = make([]*PathHandler, 0)
	}
	pathHandler := &PathHandler{
		method:  method,
		pattern: r.pathPrefix + path,
		h:       h,
	}
	r.pathHandlers = append(r.pathHandlers, pathHandler)
}

func (r *RegistryGroup) Group(pathPrefix string, ops func(*RegistryGroup)) {
	newRegistry := &RegistryGroup{
		pathPrefix: r.pathPrefix + pathPrefix,
	}
	if r.subGroup == nil {
		r.subGroup = make([]*RegistryGroup, 0)
	}
	r.subGroup = append(r.subGroup, newRegistry)
	if ops != nil {
		ops(newRegistry)
	}
}

func NewServer(port int) *Server {
	return &Server{
		port:   port,
		router: newTreeRouter(),
		group:  &RegistryGroup{},
	}
}

func (s *Server) Group(path string, ops func(*RegistryGroup)) {
	s.group.Group(path, ops)
}

func (s *Server) Get(path string, h Handler) {
	s.group.Register("GET", path, h)
}

func (s *Server) Use(m MidWare) {
	s.group.Use(m)
}

func (s *Server) Register(method string, path string, h Handler) {
	s.group.Register(method, path, h)
}

func (s *Server) mapping(method string, path string, chain *MidWareChain) {
	s.router.Register(method, path, chain)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	match := s.router.match(method, path)
	if match.matches() {
		match.chain().copy().doChain(&Context{
			Server:       s,
			Response:     w,
			Request:      r,
			PathVariable: match.getPathVariable(),
		})
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Page Not Found"))
	}
}

func (s *Server) Start() {
	s.registerGroups()
	http.Handle("/", s)
	log.Printf("Wee server start at port: %d\n", s.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
	if err != nil {
		log.Printf("server start failed. %v\n", err)
	}
}

//expand group to tree
func (s *Server) registerGroups() {
	stack := newStack()
	stack.push(s.group)
	doRegisterGroups(s, stack)
}

func (s *Server) AllUniqueMidWares() []MidWare {
	result := make([]MidWare, 0)
	return searchMidWares(s.group, result)
}

func searchMidWares(r *RegistryGroup, list []MidWare) []MidWare {
	for _, m := range r.midWares {
		exist := false
		for _, it := range list {
			if it == m {
				exist = true
				break
			}
		}
		if !exist {
			list = append(list, m)
		}
	}
	for _, sub := range r.subGroup {
		list = searchMidWares(sub, list)
	}
	return list
}

func doRegisterGroups(s *Server, stack *stack) {
	if stack.isEmpty() {
		return
	}
	currentGroup := stack.peek().(*RegistryGroup)
	if len(currentGroup.pathHandlers) > 0 {
		midWares := generateMidWare(stack)
		for _, pathHandler := range currentGroup.pathHandlers {
			s.mapping(pathHandler.method, pathHandler.pattern, appendRenderMidWare(midWares, pathHandler.h))
		}
	}
	if len(currentGroup.subGroup) > 0 {
		for _, subGroup := range currentGroup.subGroup {
			stack.push(subGroup)
			doRegisterGroups(s, stack)
		}
	}
	stack.pop()
}

func generateMidWare(stack *stack) []MidWare {
	result := make([]MidWare, 0)
	stack.forEach(func(g interface{}) {
		group := g.(*RegistryGroup)
		for _, m := range group.midWares {
			result = append(result, m)
		}
	})
	return result
}

func appendRenderMidWare(midWares []MidWare, h Handler) *MidWareChain {
	if midWares == nil {
		return &MidWareChain{
			chain: []MidWare{&RenderMidWare{h: h}},
		}
	} else {
		maxLen := len(midWares) + 1
		copyMidWare := make([]MidWare, maxLen)
		copy(copyMidWare, midWares)
		copyMidWare[maxLen-1] = &RenderMidWare{h: h}
		return &MidWareChain{
			chain: copyMidWare,
		}
	}
}

func normalize(path string) string {
	var result = path
	if !strings.HasPrefix(result, urlSep) {
		result = urlSep + result
	}
	return strings.TrimSuffix(result, urlSep)
}
