package dbrunnerservice

import (
	"os"
	"strconv"

	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	cacheModule *cacheModule

	dbrunnerv1connect.UnimplementedDbRunnerServiceHandler
}

type Options struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

type OptionFn func(*Options)

func New(optfns ...OptionFn) *Service {
	options := &Options{}
	for _, optfn := range optfns {
		optfn(options)
	}

	if options.RedisAddr == "" {
		panic("missing RedisAddr – set it with WithRedisAddr or WithRedisEnvironments")
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     options.RedisAddr,
		Password: options.RedisPassword,
		DB:       options.RedisDB,
	})

	return &Service{
		cacheModule: newCacheModule(redis),
	}
}

func WithRedisAddr(addr string) OptionFn {
	return func(o *Options) {
		o.RedisAddr = addr
	}
}

func WithRedisPassword(password string) OptionFn {
	return func(o *Options) {
		o.RedisPassword = password
	}
}

func WithRedisDB(db int) OptionFn {
	return func(o *Options) {
		o.RedisDB = db
	}
}

func WithRedisEnvironments() OptionFn {
	return func(o *Options) {
		redisDBStr := os.Getenv("REDIS_DB")
		if redisDBStr == "" {
			redisDBStr = "0"
		}
		parsedDBNum, err := strconv.Atoi(redisDBStr)
		if err != nil {
			panic("invalid REDIS_DB: " + err.Error())
		}

		o.RedisAddr = os.Getenv("REDIS_ADDR")
		o.RedisPassword = os.Getenv("REDIS_PASSWORD")
		o.RedisDB = parsedDBNum
	}
}
