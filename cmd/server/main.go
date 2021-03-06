package main

import (
	"context"
	"fmt"
	"log"

	goredis "github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ibrahimker/tigerhall-kittens/common/config"
	"github.com/ibrahimker/tigerhall-kittens/common/healthcheck"
	"github.com/ibrahimker/tigerhall-kittens/common/logging"
	"github.com/ibrahimker/tigerhall-kittens/common/postgres"
	"github.com/ibrahimker/tigerhall-kittens/common/redis"
	sightingv1 "github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1"
	"github.com/ibrahimker/tigerhall-kittens/server"
)

const (
	envDevelopment     = "development"
	maxCallRecvMsgSize = 20000000
)

func main() {
	cfg, cerr := config.NewConfig(".env")
	checkError(cfg, cerr)

	logger := initLogger(cfg)
	logger.Info("Initiating rapor-pendidikan-be server")

	pgpool, perr := postgres.NewPool(&cfg.Postgres)
	checkError(cfg, perr)
	rds, rerr := redis.NewClient(&cfg.Redis)
	checkError(cfg, rerr)

	grpcServer := createGrpcServer(cfg, logger)
	registerGrpcHandlers(grpcServer.Server, cfg, pgpool, rds, logger)

	restServer := createRestServer(cfg.Port.REST, cfg)
	registerRestHandlers(context.Background(), restServer.ServeMux, fmt.Sprintf(":%s", cfg.Port.GRPC), logger, grpc.WithTransportCredentials(insecure.NewCredentials()))

	healthcheck.RegisterHealthHandler(grpcServer.Server)

	_ = grpcServer.Run()
	_ = restServer.Run()
	_ = grpcServer.AwaitTermination()
}

func initLogger(cfg *config.Config) *logrus.Entry {
	l := logging.NewLogger()
	var logLevel logrus.Level

	env := cfg.Env
	switch env {
	case "production":
		logLevel = logrus.InfoLevel
	default:
		logLevel = logrus.DebugLevel
	}

	l.SetLevel(logLevel)
	return l.WithFields(logrus.Fields{
		"service": cfg.ServiceName,
		"version": 1,
	})
}

func createGrpcServer(cfg *config.Config, logger *logrus.Entry) *server.Grpc {
	if cfg.Env == envDevelopment {
		return server.NewDevelopmentGrpc(cfg.Port.GRPC, logger)
	}
	srv, err := server.NewProductionGrpc(cfg, logger)
	checkError(cfg, err)
	return srv
}

func createRestServer(port string, cfg *config.Config) *server.Rest {
	return server.NewProductionRest(port, cfg)
}

func registerGrpcHandlers(server *grpc.Server, cfg *config.Config, pgpool *pgxpool.Pool,
	rds *goredis.Client, logger *logrus.Entry) {
	// start register all module's gRPC handlers

	sightingv1.InitGrpc(server, cfg, pgpool, rds, logger)
	// end of register all module's gRPC handlers
}

func registerRestHandlers(ctx context.Context, server *runtime.ServeMux, grpcPort string, logger *logrus.Entry, options ...grpc.DialOption) {
	// start register all module's REST handlers
	options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize)))
	// examplev2.InitRest(ctx, server, grpcPort, options...)
	sightingv1.InitRest(ctx, server, grpcPort, logger, options...)
	// end of register all module's REST handlers
}

func checkError(cfg *config.Config, err error) {
	if err != nil {
		if cfg.IsDevelopment() {
			log.Printf("Error %+v", err)
		} else {
			log.Fatal(err)
		}
	}
}
