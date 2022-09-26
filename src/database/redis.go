package database

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

var Cache *redis.Client
var CacheChannel chan string

func SetupRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
}

func SetupCacheChannel() {
	CacheChannel = make(chan string)

	go func(ch chan string) {
		for {
			time.Sleep(3 * time.Second)

			Cache.Del(context.Background(), <-ch)
		}
	}(CacheChannel)
}

func ClearCache(keys ...string) {
	for _, key := range keys {
		CacheChannel <- key
	}
}

func SetCache(key string, value []byte, ttl time.Duration) error {
	return Cache.Set(context.Background(), key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
	result, err := Cache.Get(context.Background(), key).Result()

	if err != nil {
		return "", err
	}

	return result, nil
}
