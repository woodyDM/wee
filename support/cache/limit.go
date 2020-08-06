package cache

import (
	"fmt"
	"strconv"
	"time"
)

//an approximate rate limiter.
type RateLimiter struct {
	Name     string
	Duration int //total seconds
	Segment  int // slide windows segment
	Limit    int
	redis    *RedisService
	d        int64
	s        int64
}

func NewRateLimiter(name string, duration, segment, limit int, redis *RedisService) *RateLimiter {
	return &RateLimiter{
		Name:     name,
		Duration: duration,
		Segment:  segment,
		Limit:    limit,
		redis:    redis,
		d:        int64(duration),
		s:        int64(segment),
	}
}

func (r *RateLimiter) IsAcquired() bool {
	return r.Acquire() <= r.Limit
}

func (r *RateLimiter) Acquire() int {
	nowSeg := r.toUnixSegment(time.Now().Unix())
	segKey := segCacheKey(r.Name, nowSeg)
	init := r.redis.Incr(segKey, r.Duration)
	total := init
	for i := nowSeg - 1; i >= nowSeg-int64(r.Segment); i-- {
		s, ok := r.redis.Get(segCacheKey(r.Name, i))
		if ok {
			num, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			total += num
		}
	}
	return total
}

func segCacheKey(name string, seg int64) string {
	return fmt.Sprintf("WeeRateLimit:%s_%d", name, seg)
}

func (r *RateLimiter) toUnixSegment(stamp int64) int64 {
	s := stamp / r.d * r.s
	return s
}
