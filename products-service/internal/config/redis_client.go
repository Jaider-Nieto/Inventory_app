package config

import "github.com/redis/go-redis/v9"

func InitRedisClient(adr, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     adr,
		Password: password,
	})
}
