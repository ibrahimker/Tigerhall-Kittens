// Package sightingv1 is an entry point for starting sighting module
package sightingv1

import (
	"context"

	goredis "github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	tigerv1 "github.com/ibrahimker/tigerhall-kittens/api/proto"
	"github.com/ibrahimker/tigerhall-kittens/common/config"
	"github.com/ibrahimker/tigerhall-kittens/common/logging"
)

// InitGrpc initializes gRPC user management modules.
func InitGrpc(server *grpc.Server, cfg *config.Config, pool *pgxpool.Pool, rds *goredis.Client, logger *logrus.Entry) {
	//scorecardBuilder := builder.BuildScorecardHandler(cfg, pool, rds, logger, guruTokenAuthClient)
	tigerv1.RegisterTigerhallKittensServiceServer(server, nil)
}

// InitRest initializes REST user management modules.
// If any error occurs, it logs the error and continue the process.
func InitRest(ctx context.Context, server *runtime.ServeMux, grpcPort string, logger *logrus.Entry, options ...grpc.DialOption) {
	if err := tigerv1.RegisterTigerhallKittensServiceHandlerFromEndpoint(ctx, server, grpcPort, options); err != nil {
		logging.WithError(err, logger).Error("RegisterTigerhallKittensServiceHandlerFromEndpoint failed to be registered")
	}
}
