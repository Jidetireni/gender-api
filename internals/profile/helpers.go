package profile

import (
	"fmt"
	"time"
)

type RedisKeyPrefix string

const (
	RedisProfileNamePrefix    RedisKeyPrefix = "profile:name:%s"
	RedisProfileIDPrefix      RedisKeyPrefix = "profile:id:%s"
	RedisProfileExpirationTTL time.Duration  = 60 * time.Minute
)

func RedisProfileNameKey(name string) string {
	return fmt.Sprintf(string(RedisProfileNamePrefix), name)
}

func RedisProfileIDKey(id string) string {
	return fmt.Sprintf(string(RedisProfileIDPrefix), id)
}
