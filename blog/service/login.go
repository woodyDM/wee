package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"wee-server/blog/entity"
	"wee-server/blog/midware"
	"wee-server/blog/repository"
	"wee-server/wee"
)

const (
	maxFailLogin = 5
)

func TryLogin(context *wee.Context, name, password string) *entity.User {
	count := getLastLoginCount(name)
	if count >= maxFailLogin {
		return nil
	}
	user := repository.GetUserByName(name)
	if user == nil {
		return nil
	}
	s := digest(password, user.Salt)
	if s == user.Password {
		bytes, _ := json.Marshal(user)
		context.Session.Set(midware.UserSessionKey, string(bytes))
		cacheKey := loginKey(name)
		CacheService.Del(cacheKey)
		return user
	}
	return nil
}

func getLastLoginCount(name string) int {
	key := loginKey(name)
	return CacheService.Incr(key, 3600*4)
}

func loginKey(name string) string {
	return fmt.Sprintf("Wee.loginKey:%s", name)
}

func Logout(ctx *wee.Context) {
	ctx.Session.Delete(midware.UserSessionKey)
}

func digest(raw, salt string) string {
	sum256 := sha256.Sum256([]byte(raw + salt))
	return fmt.Sprintf("%x", sum256)
}

func GetCurrentUser(ctx *wee.Context) *entity.User {
	s, ok := ctx.Session.Get(midware.UserSessionKey)
	if !ok {
		return nil
	}
	u := new(entity.User)
	e := json.Unmarshal([]byte(s), u)
	if e != nil {
		panic(e)
	}
	return u
}
