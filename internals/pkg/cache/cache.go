package cache

import (
	"github.com/Jidetireni/gender-api/config"
	"github.com/Jidetireni/gender-api/internals/pkg/cache/redis"
)

type Cache struct {
	Redis *redis.Redis
}

func New(cfg *config.Config) (*Cache, error) {
	rds, err := redis.New(cfg)
	if err != nil {
		return nil, err
	}
	return &Cache{Redis: rds}, nil
}
