package redismodule

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var FxModule = fx.Module("redis", fx.Provide(New))

func New() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil, errors.New("missing REDIS_ADDR")
	}
	password := os.Getenv("REDIS_PASSWORD")
	db := 0

	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
		}
	}

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}), nil
}
