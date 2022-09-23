package redis

import (
	"github.com/go-redis/redis"
)

func InitialRedis(host, pass string) *redis.Client {
	rds := redis.NewClient(&redis.Options{
		Addr: host,
		Password: pass,
		DB: 0,
	})

	return rds
}
