package wee

/**
Filters or MidWare of wee Context
*/
type MidWareChain struct {
	pos   int
	chain []MidWare
}

func (mc *MidWareChain) doChain(ctx *Context) {
	if len(mc.chain) != 0 {
		mc.chain[0].Action(ctx, mc)
	}
}

func (mc *MidWareChain) Next(ctx *Context) {
	if mc.pos < len(mc.chain)-1 {
		mc.pos++
		mc.chain[mc.pos].Action(ctx, mc)
	}
}

func (mc *MidWareChain) copy() *MidWareChain {
	return &MidWareChain{
		chain: mc.chain,
	}
}

type MidWare interface {
	Action(ctx *Context, chain *MidWareChain)
}

type RenderMidWare struct {
	h Handler
}

func (r *RenderMidWare) Action(ctx *Context, chain *MidWareChain) {
	r.h.handle(ctx)
	chain.Next(ctx)
}
