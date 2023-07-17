package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"hezzl/config"
)

func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.Db,
	})

	ping := rdb.Ping()

	_, err := ping.Result()
	if err != nil {
		return nil, fmt.Errorf("error while trying to connect to redis: %w", err)
	}

	return rdb, nil
}
