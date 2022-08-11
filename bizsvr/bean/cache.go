// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package bean

import (
	"entry-task/conf"
	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
)

func InitCache(conf *conf.CacheServerConf) error {
	r := redis.NewClient(&redis.Options{
		Addr:        conf.Host,
		DB:          conf.DbIndex,
		PoolSize:    conf.PoolSize,
		MaxRetries:  conf.MaxRetries,
		IdleTimeout: conf.IdleTimeout,
	})

	RedisClient = r

	return nil
}
