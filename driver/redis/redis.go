// Package redis provides connection to redis
package redis

import (
	"context"
	"encoding/json"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

// RedisClient global variable to execute redis command
var RedisClient *goredis.Client

// Redis defines an interface for redis repository
type Redis interface {
	Del(ctx context.Context, keys ...string) error
	Fetch(ctx context.Context, key string, value interface{}, expiration time.Duration, callback func() (interface{}, error)) error
}

// RepoProvider holds client for redis
type RepoProvider struct {
	rds *goredis.Client
}

// NewRedisClient initialize redis client instance
func NewRedisClient(rds *goredis.Client) *RepoProvider {
	RedisClient = rds

	return &RepoProvider{rds: RedisClient}
}

// Del deletes data from redis
func (r *RepoProvider) Del(ctx context.Context, keys ...string) error {
	cacheKeys := make([]string, 0, len(keys))

	cacheKeys = append(cacheKeys, keys...)
	err := r.rds.Del(ctx, cacheKeys...).Err()
	return err
}

// Fetch retrieves data from redis
// If data not found, it repopulate the data into redis
func (r *RepoProvider) Fetch(ctx context.Context, key string, value interface{}, expiration time.Duration, callback func() (interface{}, error)) error {
	if err := get(ctx, key, value); err != nil && err != goredis.Nil {
		return err
	} else if err == nil {
		return nil
	}

	// When redis get doesn't raise error and returns nil
	callbackValue, err := callback()
	if err != nil {
		return err
	}

	setBytes, err := json.Marshal(callbackValue)
	if err != nil {
		return err
	}

	if err = set(ctx, key, setBytes, expiration); err != nil {
		return err
	}

	return json.Unmarshal(setBytes, value)
}

func get(ctx context.Context, key string, value interface{}) error {
	bytes, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, value)
	return err
}

func set(ctx context.Context, key string, bytes []byte, expiration time.Duration) error {
	err := RedisClient.Set(ctx, key, bytes, expiration).Err()
	return err
}
