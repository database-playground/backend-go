package dbrunnerservice

import (
	"os"
	"strconv"

	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
	"github.com/redis/go-redis/v9"
)

type DBRunnerService struct {
	cacheModule *cacheModule

	dbrunnerv1connect.UnimplementedDbRunnerServiceHandler
}

type DBRunnerServiceOptions struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

type NewDBRunnerServiceOptionFn func(*DBRunnerServiceOptions)

func NewDBRunnerService(optfns ...NewDBRunnerServiceOptionFn) *DBRunnerService {
	options := &DBRunnerServiceOptions{}
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

	return &DBRunnerService{
		cacheModule: newCacheModule(redis),
	}
}

func WithRedisAddr(addr string) NewDBRunnerServiceOptionFn {
	return func(o *DBRunnerServiceOptions) {
		o.RedisAddr = addr
	}
}

func WithRedisPassword(password string) NewDBRunnerServiceOptionFn {
	return func(o *DBRunnerServiceOptions) {
		o.RedisPassword = password
	}
}

func WithRedisDB(db int) NewDBRunnerServiceOptionFn {
	return func(o *DBRunnerServiceOptions) {
		o.RedisDB = db
	}
}

func WithRedisEnvironments() NewDBRunnerServiceOptionFn {
	return func(o *DBRunnerServiceOptions) {
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
