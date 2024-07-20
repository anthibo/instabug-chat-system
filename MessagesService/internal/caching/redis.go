package caching

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCacheManager struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string, password string, db int) *RedisCacheManager {
	log.Println("Connecting to Redis...")
	log.Println("Address: ", addr)
	log.Println("Password: ", password)
	log.Println("DB: ", db)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCacheManager{
		client: redisClient,
		ctx:    context.Background(),
	}
}

func (r *RedisCacheManager) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisCacheManager) Get(key string) (interface{}, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *RedisCacheManager) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
