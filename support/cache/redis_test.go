package cache

import (
	"testing"
	"time"
)

func TestRedisService_Do(t *testing.T) {
	service := NewRedisService("localhost", 6379, "123456")
	result, e := service.Get("Key")
	defer service.Close()
	t.Log(result)
	t.Log(e)
	e3 := service.Set("Key", "msg", 60)
	t.Log(e3)
	result2, _ := service.Get("Key")
	t.Log(result2)
	r3, _ := service.Do("TTL", "Key")
	t.Log(r3)
	e3 = service.HSetEx("k", "inK1", "qwert", 566)
	e4 := service.HSetEx("k", "inK2", "abcdef", 566)
	t.Log(e3)
	t.Log(e4)
	r4, _ := service.HGet("k", "inK1")
	r5, _ := service.HGet("k", "inK2")
	s, b := service.HGet("k", "inKN")
	t.Log(s, b)
	t.Log(r4)
	t.Log(r5)

}

func TestRedisService_SetNx(t *testing.T) {
	service := NewRedisService("localhost", 6379, "123456")
	defer service.Close()
	e := service.SetNx("Key1", "v", 100)
	if e != nil {
		t.Fatal("Failed to set ex")
	}
	s, ok := service.Get("Key1")
	if !ok {
		t.Fatal("Failedto get")
	}
	if s != "v" {
		t.Fatal("failed to set.")
	}
	service.Del("Key1")
	_, ok2 := service.Get("Key1")
	if ok2 {
		t.Fatal("should not get")
	}

}

func TestIncr(t *testing.T) {
	service := NewRedisService("localhost", 6379, "123456")
	defer service.Close()
	defer service.Del("a")
	i := service.Incr("a", 100)
	if i != 1 {
		t.Fail()
	}
	i = service.Incr("a", 100)
	if i != 2 {
		t.Fail()
	}

}

func TestRateLimiter_Acquire(t *testing.T) {
	service := NewRedisService("localhost", 6379, "123456")
	defer service.Close()
	limiter := NewRateLimiter("rlimit", 6, 2, 3, service)
	i1 := limiter.Acquire()
	if i1 != 1 {
		t.Errorf("1")
	}
	time.Sleep(2 * time.Second)

	i1 = limiter.Acquire()
	i1 = limiter.Acquire()
	if i1 != 3 {
		t.Errorf("2")
	}
	time.Sleep(2 * time.Second)

	i1 = limiter.Acquire()
	i1 = limiter.Acquire()
	if i1 != 5 {
		t.Errorf("3:should be 5,but %d", i1)
	}
	time.Sleep(3 * time.Second)

	i1 = limiter.Acquire()
	if i1 != 5 {
		t.Errorf("4:should be 5,but %d", i1)
	}
	time.Sleep(2 * time.Second)

	i1 = limiter.Acquire()
	if i1 != 6 {
		t.Errorf("5:shoud be 6 but %d", i1)
	}

}

type CacheObject struct {
	Name  string
	Age   int
	Words []string
}

func TestRedisService_GetCache(t *testing.T) {
	service := NewRedisService("localhost", 6379, "123456")
	defer service.Close()
	//defer service.Del("CacheKey")
	result, ok := service.GetCache("CacheKey", 100, CacheObject{}, func() interface{} {
		return &CacheObject{
			Name:  "Kitty",
			Age:   100,
			Words: []string{"a", "A", "b", "B"},
		}
	})
	if !ok {
		t.Fatal("should ok")
	}
	if result == nil {
		t.Fatal("should not nil")
	}

}
func TestRedisService_GetCache_Nil(t *testing.T) {
	service := NewRedisService("localhost", 6379, "123456")
	defer service.Close()
	//defer service.Del("CacheKey")
	result, ok := service.GetCache("CacheKey", 100, CacheObject{}, func() interface{} {
		return nil
	})
	if ok {
		t.Fatal("should not ok")
	}
	if result != nil {
		t.Fatal("should nil")
	}

}
