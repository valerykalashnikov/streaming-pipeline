package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const dbFileName = "processed_files"

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{
		client: redisClient,
	}
}

func (rc *RedisClient) AddFilename(ctx context.Context, filename string) error {
	return rc.client.LPush(ctx, dbFileName, filename).Err()
}

func (rc *RedisClient) GetProcessedFilesList(ctx context.Context) ([]string, error) {
	return rc.client.LRange(ctx, dbFileName, 0, -1).Result()
}
