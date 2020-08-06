package midware

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"wee-server/blog/entity"
	"wee-server/blog/repository"
	"wee-server/support/cache"
	"wee-server/support/web"
	"wee-server/wee"
)

type LogMidWare struct {
}
type AuthMidWare struct {
}
type AspectMidWare struct {
	Name string
}
type RateLimitMidWare struct {
	Redis    *cache.RedisService
	Duration int
	Segment  int
	Limit    int
}

const (
	UserSessionKey = "currentUser"
)

func (r *RateLimitMidWare) Action(ctx *wee.Context, chain *wee.MidWareChain) {
	ua := ctx.Request.Header.Get("User-Agent")
	ip := web.RemoteIp(ctx.Request)
	key := fmt.Sprintf("%s_%s", ua, ip)
	sum256 := sha256.Sum256([]byte(key))
	key = fmt.Sprintf("%s_%x", "common", sum256)
	limiter := cache.NewRateLimiter(key, r.Duration, r.Segment, r.Limit, r.Redis)
	ok := limiter.IsAcquired()
	if ok {
		chain.Next(ctx)
	} else {
		ctx.Response.WriteHeader(http.StatusTooManyRequests)
	}
}

func (a *AuthMidWare) Action(ctx *wee.Context, chain *wee.MidWareChain) {
	path := ctx.Request.URL.Path
	_, ok := ctx.Session.Get(UserSessionKey)
	if !ok {
		ctx.Response.WriteHeader(401)
		log.Printf("[%s][auth=false]:path:%s\n", ctx.Session.SessionId(), path)
	} else {
		log.Printf("[%s][auth=true]:path:%s\n\n", ctx.Session.SessionId(), path)
		chain.Next(ctx)
	}
}

func (m *AspectMidWare) Action(ctx *wee.Context, chain *wee.MidWareChain) {
	path := ctx.Request.URL.Path
	log.Printf("*********************************************************\n")
	log.Printf("%s --->[%s]\n", m.Name, path)
	if "/health" != path {
		go func() {
			m.saveHistoryTodataBase(ctx, path)
		}()
	}
	chain.Next(ctx)
	log.Printf("%s <---[%s]\n\n", m.Name, path)
}

func (m *AspectMidWare) saveHistoryTodataBase(ctx *wee.Context, path string) {
	ua := ctx.Request.Header.Get("User-Agent")
	fullIp := web.RemoteIp(ctx.Request)
	hash := web.HashCode(ua + "@" + fullIp)
	history := &entity.History{
		Id:         0,
		CreateTime: time.Now(),
		UserAgent:  ua,
		FullIp:     fullIp,
		TrimIp:     strings.Split(fullIp, ",")[0],
		VisitPath:  path,
		UserHash:   hash,
	}
	log.Printf("Hash is %d", hash)
	repository.SaveHistory(history)
}

func (l *LogMidWare) Action(ctx *wee.Context, chain *wee.MidWareChain) {
	path := ctx.Request.URL.Path
	start := time.Now()
	chain.Next(ctx)
	cost := time.Now().Sub(start)
	log.Printf("处理:%s 花费 %d ms\n", path, cost.Milliseconds())
}
