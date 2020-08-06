package web

import (
	"wee-server/support/cache"
	"wee-server/wee"
)

type RedisSessionHolder struct {
	ExpireSeconds int
	Service       *cache.RedisService
}

type redisSession struct {
	sessionId     string
	redisHKey     string
	service       *cache.RedisService
	ExpireSeconds int
}

func (r *redisSession) Delete(key string) {
	_ = r.service.HDel(r.redisHKey, key)
}

func (r *redisSession) SessionId() string {
	return r.sessionId
}

func (r *redisSession) Get(key string) (string, bool) {
	return r.service.HGet(r.redisHKey, key)
}

func (r *redisSession) Set(key string, value string) {
	err := r.service.HSetEx(r.redisHKey, key, value, r.ExpireSeconds)
	if err != nil {
		panic("Failed to write to redis")
	}
}

func sessionKey(sessionId string) string {
	return "wee.Gsession" + ":" + sessionId
}

func (r *RedisSessionHolder) GetSession(sessionId string) (wee.Session, bool) {
	redisKey := sessionKey(sessionId)
	_, _ = r.Service.Do("EXPIRE", redisKey, r.ExpireSeconds)
	return &redisSession{
		sessionId:     sessionId,
		redisHKey:     redisKey,
		service:       r.Service,
		ExpireSeconds: r.ExpireSeconds,
	}, true
}

func (r *RedisSessionHolder) SetSession(sessionId string, s wee.Session) {
	//do nothing
}
