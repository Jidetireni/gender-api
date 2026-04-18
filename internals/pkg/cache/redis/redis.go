package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Jidetireni/gender-api/config"
	rds "github.com/redis/go-redis/v9"
)

type Redis struct {
	client *rds.Client
}

func New(cfg *config.Config) (*Redis, error) {
	opts, err := rds.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	client := rds.NewClient(opts)

	Redis := &Redis{
		client: client,
	}

	if err := Redis.Ping(); err != nil {
		return nil, err
	}

	return Redis, nil
}

// Ping checks database connection
func (r *Redis) Ping() error {
	return r.client.Ping(context.Background()).Err()
}

// Set JSON value in Redis
func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, v, expiration).Err()
}

// Get JSON value from Redis
func (r *Redis) Get(ctx context.Context, key string, dest any) error {
	v, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, rds.Nil) {

		}
		return err
	}

	return json.Unmarshal([]byte(v), dest)
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
