package dbrunnerservice

import (
	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var FxModule = fx.Module("dbrunner-service", fx.Provide(New))

type Service struct {
	cacheModule *CacheModule

	dbrunnerv1connect.UnimplementedDbRunnerServiceHandler
}

func New(redis *redis.Client) *Service {
	return &Service{
		cacheModule: NewCacheModule(redis),
	}
}
