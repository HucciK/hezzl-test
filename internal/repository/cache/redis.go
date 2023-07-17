package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisCache struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func NewRedisCache(cli *redis.Client, defaultTTL int) *RedisCache {
	return &RedisCache{
		client:     cli,
		defaultTTL: time.Duration(defaultTTL) * time.Second,
	}
}

func (r RedisCache) Set(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error while trying convert value to json: %w", err)
	}

	status := r.client.Set(key, jsonData, r.defaultTTL)
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

func (r RedisCache) Get(key string) ([]byte, error) {
	res := r.client.Get(key)
	if err := res.Err(); err != nil {
		return nil, err
	}

	return res.Bytes()
}

func (r RedisCache) Delete(key string) error {
	res := r.client.Del(key)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}
