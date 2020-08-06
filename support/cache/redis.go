package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"reflect"
	"time"
)

type RedisService struct {
	pool *redis.Pool
	host string
	port int
	auth string
}

const (
	hsetex = `
local result = redis.call('hset', KEYS[1], KEYS[2], ARGV[1])
if result < 0
then
	return -1
else
	return redis.call('expire',KEYS[1], ARGV[2])
end
`

	nilCache = "CACHE_VALUE_NIL"
	incre    = `
local result = redis.call('incr', KEYS[1])
redis.call('expire',KEYS[1], ARGV[1])
return result
`
)

func (s *RedisService) Get(key string) (string, bool) {
	v, e := s.doString("GET", key)
	if e != nil {
		if e == redis.ErrNil {
			return "", false
		}
		panic(e)
	}
	return v, true
}

func (s *RedisService) HGet(key, hashKey string) (string, bool) {
	v, e := s.doString("HGET", key, hashKey)
	if e != nil {
		return "", false
	}
	return v, true
}
func (s *RedisService) HSet(key, hashKey, value string) error {
	_, e := s.Do("HSET", key, hashKey, value)
	return e
}
func (s *RedisService) HDel(key, hashKey string) error {
	_, e := s.Do("HDEL", key, hashKey)
	return e
}

func (s *RedisService) HSetEx(key, hashKey, value string, ttl int) error {
	conn := s.pool.Get()
	defer conn.Close()
	lua := redis.NewScript(2, hsetex)
	result, e := lua.Do(conn, key, hashKey, value, ttl)
	if e != nil {
		return e
	}
	if result.(int64) != 1 {
		return errors.New("Failed to hsetex. ")
	}
	return nil
}

func (s *RedisService) Incr(key string, ttl int) int {
	conn := s.pool.Get()
	defer conn.Close()
	lua := redis.NewScript(1, incre)
	result, e := lua.Do(conn, key, ttl)
	if e != nil {
		panic(e)
	}
	if r, ok := result.(int64); ok {
		return int(r)
	} else {
		panic("Incr is expected to return int64")
	}
}

func (s *RedisService) Set(key, value string, ttl int) error {
	_, e := s.Do("SETEX", key, ttl, value)
	return e
}

func (s *RedisService) SetNx(key, value string, ttl int) error {
	_, e := s.Do("SET", key, value, "EX", ttl, "NX")
	return e
}

func (s *RedisService) Del(key string) {
	_, _ = s.Do("DEL", key)
}

func (s *RedisService) WrapCacheKey(key string) string {
	return "Wee.Cache:" + key
}

//return false result is nil or true when result is not nil
func (s *RedisService) GetCache(key string, ttl int, class interface{}, valueProvider func() interface{}) (interface{}, bool) {

	jsonString, ok := s.Get(key)
	if ok {
		if jsonString == nilCache {
			return nil, false
		} else {
			var result = reflect.New(reflect.ValueOf(class).Type()).Interface()
			err := json.Unmarshal([]byte(jsonString), result)
			if err != nil {
				panic(err)
			}
			return result, true
		}
	}
	var cacheString string
	cacheTarget := valueProvider()
	nilTarget := IsNil(cacheTarget)
	if nilTarget {
		cacheString = nilCache
	} else {
		bytes, e := json.Marshal(cacheTarget)
		if e != nil {
			panic(e)
		}
		cacheString = string(bytes)
		if cacheString == nilCache {
			panic(fmt.Errorf("GetCache don't not support value same as nil key:%s ", nilCache))
		}
	}
	//write to cache
	e := s.SetNx(key, cacheString, ttl)
	if e != nil {
		panic(e)
	}
	return cacheTarget, !nilTarget
}

func IsNil(o interface{}) bool {
	defer func() {
		recover()
	}()
	v := reflect.ValueOf(o)
	return v.IsNil()

}

func (s *RedisService) doString(commandName string, args ...interface{}) (string, error) {
	return redis.String(s.Do(commandName, args...))
}

func (s *RedisService) Do(commandName string, args ...interface{}) (interface{}, error) {
	conn := s.pool.Get()
	defer conn.Close()
	if e := conn.Err(); e != nil {
		return nil, e
	}
	i, e := conn.Do(commandName, args...)
	return i, e
}

func (s *RedisService) Close() {
	_ = s.pool.Close()
}

func NewRedisService(host string, port int, auth string) *RedisService {
	service := &RedisService{
		pool: &redis.Pool{
			Dial: func() (conn redis.Conn, e error) {
				return dial(host, port, auth)
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, e := c.Do("PING")
				return e
			},
			MaxIdle:         1,
			MaxActive:       8,
			IdleTimeout:     240 * time.Second,
			Wait:            true,
			MaxConnLifetime: 0,
		},
		host: host,
		port: port,
		auth: auth,
	}
	service.Get("dummy")
	log.Printf("Create new RedisService to %s. \n", host)
	return service
}

func dial(host string, port int, auth string) (redis.Conn, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	if auth != "" {
		if _, err := conn.Do("AUTH", auth); err != nil {
			_ = conn.Close()
			log.Printf("Invalid redis auth, please check your configuraion.")
			panic(err)
		}
	}
	return conn, nil
}
