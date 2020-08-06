package blog

import (
	"log"
	"wee-server/blog/midware"
	"wee-server/blog/repository"
	"wee-server/blog/service"
	"wee-server/support"
	"wee-server/support/cache"
	"wee-server/support/database"
	web2 "wee-server/support/web"
)
import "wee-server/wee"

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
}

func Start(config *support.Configuration) {
	//component define
	mysql := database.NewMysql(config.Mysql.Url)
	defer mysql.Close()
	redisService := cache.NewRedisService(config.Redis.Host, config.Redis.Port, config.Redis.Auth)
	defer redisService.Close()
	repository.Mysql = mysql
	service.CacheService = redisService
	service.CacheTimeout = config.Web.CacheTimeoutSeconds
	//midWare define
	sessionExpire := config.Web.SessionSeconds
	sessionMidWare := wee.NewSimpleSessionMidware(sessionExpire, &web2.RedisSessionHolder{
		ExpireSeconds: sessionExpire,
		Service:       redisService,
	})
	authMidWare := new(midware.AuthMidWare)
	rateLimitMidWare := &midware.RateLimitMidWare{
		Redis:    redisService,
		Duration: 12,
		Segment:  3,
		Limit:    30,
	}
	//server config
	server := wee.NewServer(config.Web.Port)
	server.Use(&web2.PanicHandler{})
	server.Use(&midware.AspectMidWare{Name: "请求"})
	server.Use(&midware.LogMidWare{})
	server.Use(rateLimitMidWare)
	server.Use(sessionMidWare)
	server.Get("/", ApiHealthIndex)
	server.Get("/health", ApiHealthIndex)
	server.Group("/api", func(g *wee.RegistryGroup) {
		g.Group("/article", func(g2 *wee.RegistryGroup) {
			g2.Register("GET", "/list/:userId", ArticleListController)
			g2.Register("GET", "/:id", ArticleViewController)
		})
		g.Group("/user", func(g2 *wee.RegistryGroup) {
			g2.Use(authMidWare)
			g2.Register("GET", "/article/:id", UserArticleViewController)
			g2.Register("POST", "/article.json", UserSaveArticleController)
			g2.Register("POST", "/article/:id", UserUpdateArticleController)
		})
		g.Register("POST", "/login.json", Login)
		g.Register("POST", "/logout.json", Logout)
	})
	log.Printf("Configuration Web with session:%ds, cache timeout %ds.", config.Web.SessionSeconds, config.Web.CacheTimeoutSeconds)
	server.Start()
}
