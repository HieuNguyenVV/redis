package cache

import (
	"context"
	"fmt"
	redisv8 "github.com/go-redis/redis/v8"
	"time"
)

type Redis interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(cxt context.Context, key string) (interface{}, error)
	Exist(ctx context.Context, key string) (int64, error)
	MSet(ctx context.Context, entries []Entry, ttl time.Duration) error
	MGet(ctx context.Context, keys []string) ([]interface{}, error)
	Del(ctx context.Context, key string) error
	MDel(ctx context.Context, keys []string) error
	Close()
}

type redis struct {
	client *redisv8.Client
}

type Entry struct {
	Key   string
	Value interface{}
}

func NewRedis(host string, username string, password string, writeTimeout time.Duration,
	readTimeout time.Duration) (Redis, error) {
	redisClient := redisv8.NewClient(&redisv8.Options{
		Addr:         host,
		Username:     username,
		Password:     password,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("pig redis error,err: %v", err)
	}
	return &redis{client: redisClient}, nil
}

func (r *redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	_, err := r.client.Set(ctx, key, value, ttl).Result()
	return err
}

func (r *redis) Get(ctx context.Context, key string) (interface{}, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, err
}

func (r *redis) Exist(ctx context.Context, key string) (int64, error) {
	exist, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return exist, nil
}

func (r *redis) MSet(ctx context.Context, entries []Entry, ttl time.Duration) error {
	pipe := r.client.Pipeline()

	for _, e := range entries {
		pipe.Set(ctx, e.Key, e.Value, ttl)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *redis) MGet(ctx context.Context, keys []string) ([]interface{}, error) {
	result, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *redis) Del(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redis) MDel(ctx context.Context, keys []string) error {
	pipe := r.client.Pipeline()
	for _, key := range keys {
		pipe.Del(ctx, key)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return err
}

func (r *redis) Close() {
	r.client.Close()
}
