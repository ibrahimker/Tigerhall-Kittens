package sightingv1_test

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"

	"github.com/ibrahimker/tigerhall-kittens/common/config"
	"github.com/ibrahimker/tigerhall-kittens/common/logging"
	sightingv1 "github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1"
)

func TestInitGrpc(t *testing.T) {
	t.Run("successfully build Tiger Sighting GRPC", func(t *testing.T) {
		pool := &pgxpool.Pool{}
		cfg, _ := config.NewConfig("../../../../test/fixture/env.valid")
		rds, _ := redismock.NewClientMock()

		sightingv1.InitGrpc(grpc.NewServer(), cfg, pool, rds, logging.NewTestLogger())
	})
}

func TestInitRest(t *testing.T) {
	t.Run("successfully build Tiger Sighting Rest", func(t *testing.T) {
		options := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(20000000))
		sightingv1.InitRest(context.Background(), runtime.NewServeMux(), ":8081", logging.NewTestLogger(), options)
	})
}
