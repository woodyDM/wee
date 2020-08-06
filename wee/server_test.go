package wee

import "testing"

func TestRegistryGroup_Register(t *testing.T) {
	s := NewServer(10000)

	w := newMid("m1")
	w2 := newMid("m2")
	w3 := newMid("m3")
	w4 := newMid("m4")

	s.Register("GET", "common/page", h)
	s.Register("GET", "common/award", h)
	s.Group("/user", func(r *RegistryGroup) {

		r.Use(w)
		r.Use(w2)
		r.Register("GET", "/page", h)
		r.Register("GET", "/index", h)
		r.Register("GET", "", hIndex)
		r.Group("/:id", func(r *RegistryGroup) {
			r.Register("GET", "/", h)
			r.Register("GET", "/profile", h)
			r.Register("POST", "/profile", h)
			r.Group("/:activityId", func(r *RegistryGroup) {
				r.Register("GET", "/page", h)
				r.Use(w3)
			})
		})
		r.Group("/:orgId", func(r *RegistryGroup) {
			r.Register("GET", "/index", h)
			r.Use(w4)
		})

	})
	s.Use(newMid("root"))

	assertTrueMsg(s.group != nil, "server should init with default group", t)
	assertTrueMsg(len(s.group.pathHandlers) == 2, "server with 2 handler", t)
	assertTrueMsg(len(s.group.midWares) == 1, "server no handler", t)
	assertTrueMsg(s.group.midWares[0].(*MockMidWare).name == "root", "root", t)
	assertTrueMsg(len(s.group.subGroup) == 1, "group len 1", t)
	g1 := s.group.subGroup[0]
	assertTrueMsg(g1 != nil, "g1 not nil", t)
	assertTrueMsg(len(g1.midWares) == 2, "g1 len 2", t)
	assertTrueMsg(g1.midWares[0].(*MockMidWare).name == "m1", "m1", t)
	assertTrueMsg(g1.midWares[1].(*MockMidWare).name == "m2", "m2", t)
	assertTrueMsg(len(g1.subGroup) == 2, "g1 sub len 2", t)
	assertTrueMsg(len(g1.pathHandlers) == 3, "g1 handler len 3", t)
	assertTrueMsg(g1.pathHandlers[0].pattern == "/user/page", "/page", t)
	assertTrueMsg(g1.pathHandlers[1].pattern == "/user/index", "/index", t)
	assertTrueMsg(g1.pathHandlers[2].pattern == "/user", "/user", t)

	gid := g1.subGroup[0]
	assertTrueMsg(len(gid.midWares) == 0, "gid mid len 0", t)
	assertTrueMsg(len(gid.pathHandlers) == 3, "gid h len 3", t)
	assertTrueMsg(len(gid.subGroup) == 1, "gid sub len 1", t)
	assertTrueMsg(gid.pathHandlers[0].pattern == "/user/:id", "/:id", t)
	assertTrueMsg(gid.pathHandlers[1].pattern == "/user/:id/profile", "profile", t)
	assertTrueMsg(gid.pathHandlers[2].pattern == "/user/:id/profile", "profile", t)
	assertTrueMsg(gid.pathHandlers[2].method == "POST", "profile post", t)
	last := gid.subGroup[0]
	assertTrueMsg(len(last.subGroup) == 0, "last sub 0", t)
	assertTrueMsg(len(last.midWares) == 1, "last mid len 4", t)
	assertTrueMsg(last.midWares[0].(*MockMidWare).name == "m3", "last m3 ", t)
	assertTrueMsg(len(last.pathHandlers) == 1, "last h len 1", t)

	gid2 := g1.subGroup[1]
	assertTrueMsg(gid2.subGroup == nil, "gid2 sub group nil", t)
	assertTrueMsg(len(gid2.midWares) == 1, "gid2 mid 3", t)
	assertTrueMsg(gid2.midWares[0].(*MockMidWare).name == "m4", "gid2 m4", t)
	assertTrueMsg(len(gid2.pathHandlers) == 1, "gid2 h 1", t)
	assertTrueMsg(gid2.pathHandlers[0].pattern == "/user/:orgId/index", "gid2 h pattern", t)

	//when
	s.registerGroups()
	//check routers
	tr := s.router.(*TreeRouter).tr
	methodNode := tr.root.children
	assertTrueMsg(len(methodNode) == 2, "method len 2", t)
	assertTrueMsg(methodNode[0].pattern == "GET", "GET", t)
	assertTrueMsg(methodNode[1].pattern == "POST", "POST", t)
	assertTrueMsg(methodNode[0].value == nil, "GET value nil", t)
	assertTrueMsg(methodNode[1].value == nil, "POST value nil", t)
	//check POST
	postNodeC := methodNode[1].children
	assertTrueMsg(len(postNodeC) == 1, "postNodec len 1", t)
	assertTrueMsg(postNodeC[0].value == nil, "postNodec 0 v nil", t)
	assertTrueMsg(postNodeC[0].pattern == "user", "postNodec user", t)
	assertTrueMsg(len(postNodeC[0].children) == 1, "postNodeC c len 1", t)
	userNode := postNodeC[0].children[0]
	assertTrueMsg(userNode.value == nil, "userNode v nil", t)
	assertTrueMsg(len(userNode.children) == 1, "userNode c len 1", t)
	assertTrueMsg(userNode.pattern == ":id", ":id", t)
	profileNode := userNode.children[0]
	assertTrueMsg(profileNode.pattern == "profile", ":id post", t)
	profileMids := profileNode.value.(*MidWareChain).chain
	assertTrueMsg(len(profileMids) == 4, "pnode v not nil", t)
	assertMid(profileMids[0], "root", t)
	assertMid(profileMids[1], "m1", t)
	assertMid(profileMids[2], "m2", t)
	//check GET
	getNodes := methodNode[0].children
	assertTrueMsg(len(getNodes) == 2, "get c len 2", t)
	//check GET common
	commonNode := getNodes[0]
	assertTrueMsg(commonNode.pattern == "common", "c", t)
	assertTrueMsg(commonNode.value == nil, "common v nil", t)
	assertTrueMsg(len(commonNode.children) == 2, "commn c len 2", t)
	assertTrueMsg(commonNode.children[0].pattern == "award", "common award", t)
	assertTrueMsg(commonNode.children[1].pattern == "page", "common page", t)
	awardChain := commonNode.children[0].value.(*MidWareChain).chain
	assertMid(awardChain[0], "root", t)
	assertTrueMsg(len(awardChain) == 2, "award chain 2", t)
	pageChain := commonNode.children[1].value.(*MidWareChain).chain
	assertMid(pageChain[0], "root", t)
	assertTrueMsg(len(pageChain) == 2, "award chain 2", t)
	//check GET user
	uNode := getNodes[1]
	assertTrueMsg(uNode.pattern == "user", "user v ", t)
	uChain := uNode.value.(*MidWareChain).chain
	assertTrueMsg(len(uChain) == 4, "uChain len 4", t)
	assertMid(uChain[0], "root", t)
	assertMid(uChain[1], "m1", t)
	assertMid(uChain[2], "m2", t)
	assertTrueMsg(uNode.pattern == "user", "user p", t)
	//check GET user/index
	idxNode := uNode.children[0]
	idxChain := idxNode.value.(*MidWareChain).chain
	assertTrueMsg(len(idxChain) == 4, "uChain len 4", t)
	assertTrueMsg(idxNode.pattern == "index", "index n", t)
	assertMid(idxChain[0], "root", t)
	assertMid(idxChain[1], "m1", t)
	assertMid(idxChain[2], "m2", t)
	//check GET user/:id
	idNode := uNode.children[2]
	idChain := idNode.value.(*MidWareChain).chain
	assertTrueMsg(len(idChain) == 4, "idChain len 4", t)
	assertTrueMsg(idNode.pattern == ":id", ":id n", t)
	assertMid(idChain[0], "root", t)
	assertMid(idChain[1], "m1", t)
	assertMid(idChain[2], "m2", t)
	//check GET user/:orgId
	orgIdNode := uNode.children[3]
	assertTrueMsg(orgIdNode.value == nil, "orgIdNode v nil", t)
	//check user/:id/profile
	pNode := idNode.children[0]
	pChain := pNode.value.(*MidWareChain).chain
	assertTrueMsg(len(pChain) == 4, "pChain len 4", t)
	assertTrueMsg(pNode.pattern == "profile", "profile n", t)
	assertMid(pChain[0], "root", t)
	assertMid(pChain[1], "m1", t)
	assertMid(pChain[2], "m2", t)
	//check user/:id/:activityId
	acNode := idNode.children[1]
	assertTrueMsg(acNode.pattern == ":activityId", "acn p", t)
	assertTrueMsg(acNode.value == nil, "acn v nil", t)
	assertTrueMsg(len(acNode.children) == 1, "ac c len 1", t)
	//check user/:id/:activityId/page
	pageNode := acNode.children[0]
	assertTrueMsg(pageNode.pattern == "page", "page v", t)
	pgChain := pageNode.value.(*MidWareChain).chain
	assertTrueMsg(len(pgChain) == 5, "pgChain len 4", t)
	assertMid(pgChain[0], "root", t)
	assertMid(pgChain[1], "m1", t)
	assertMid(pgChain[2], "m2", t)
	assertMid(pgChain[3], "m3", t)
	//check user/:orgId/index
	orgIdxNode := orgIdNode.children[0]
	assertTrueMsg(orgIdxNode.pattern == "index", "index v", t)
	orgIdxChain := orgIdxNode.value.(*MidWareChain).chain
	assertTrueMsg(len(orgIdxChain) == 5, "orgIdxChain len 4", t)
	assertMid(orgIdxChain[0], "root", t)
	assertMid(orgIdxChain[1], "m1", t)
	assertMid(orgIdxChain[2], "m2", t)
	assertMid(orgIdxChain[3], "m4", t)
}

func h(ctx *Context) {
}

func hIndex(ctx *Context) {
}

type MockMidWare struct {
	name string
}

func (m MockMidWare) Init() {
}

func (m MockMidWare) Destroy() {
}

func newMid(n string) MidWare {
	return &MockMidWare{name: n}
}

func (m MockMidWare) Action(ctx *Context, chain *MidWareChain) {
	chain.Next(ctx)
}

func Test_normalize(t *testing.T) {
	assertTrueMsg(normalize("path/ok") == "/path/ok", "prefix ", t)
	assertTrueMsg(normalize("/path/ok") == "/path/ok", "prefix ", t)
	assertTrueMsg(normalize("/path/ok/") == "/path/ok", "prefix ", t)
}
