// Package builder provides functionality to build the representative flows.
// Imagine this package as an implementation of builder design pattern.
// Read more: https://sourcemaking.com/design_patterns/builder.
package builder

import (
	goredis "github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/ibrahimker/tigerhall-kittens/common/config"
	"github.com/ibrahimker/tigerhall-kittens/driver/redis"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/internal/grpc/handler"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/repository/postgres"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/service"
)

// BuildTigerSightingHandler builds tiger sighting handler including all of its dependencies.
func BuildTigerSightingHandler(cfg *config.Config, pool *pgxpool.Pool, rds *goredis.Client, logger *logrus.Entry) *handler.TigerSighting {
	redisRepo := redis.NewRedisClient(rds)
	tigerSightingRepo := postgres.NewTigerSightingRepo(pool)
	tigerSightingService := service.NewTigerSightingService(tigerSightingRepo, redisRepo)
	return handler.NewTigerSighting(logger, tigerSightingService)
}
