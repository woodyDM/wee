package web

import (
	"fmt"
	"log"
	"runtime"
	"time"
	"wee-server/wee"
)

type PanicHandler struct {
}

type response struct {
	Message string
	Detail  string
	Time    time.Time
}

func (h *PanicHandler) Action(ctx *wee.Context, chain *wee.MidWareChain) {
	defer func() {
		switch p := recover(); p {
		case nil:
		//
		default:
			ctx.Response.WriteHeader(500)
			detail := fmt.Sprintf("%v", p)
			r := &response{
				Message: "Internal Server Error",
				Detail:  detail,
				Time:    time.Now(),
			}
			log.Printf("[500]%s\n%s\n\n", detail, stack())
			ctx.Json(r)
		}
	}()
	chain.Next(ctx)
}

func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}
