// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache struct {
	Cache *redis.Client
}

func NewCache(cacheClient *redis.Client) *RedisCache {
	return &RedisCache{Cache: cacheClient}
}

func (cache *RedisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := cache.Cache.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (cache *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := cache.Cache.Set(ctx, key, value, expiration).Err()
	return err
}

func (cache *RedisCache) Delete(ctx context.Context, key string) error {
	err := cache.Cache.Del(ctx, key).Err()
	return err
}
