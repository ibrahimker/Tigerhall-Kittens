// Package server provides HTTP/2 gRCP server functionality.
package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ibrahimker/tigerhall-kittens/common/config"
	"github.com/ibrahimker/tigerhall-kittens/common/logging"
)

const (
	connProtocol = "tcp"
)

// Grpc is responsible to act as gRPC server.
// It composes grpc.Server.
type Grpc struct {
	*grpc.Server
	listener net.Listener
	port     string
}

// NewGrpc creates an instance of Grpc.
func NewGrpc(port string, options ...grpc.ServerOption) *Grpc {
	srv := grpc.NewServer(options...)
	return &Grpc{
		Server: srv,
		port:   port,
	}
}

// NewDevelopmentGrpc creates an instance of Grpc for used in development environment.
//
// These are list of interceptors that are attached (from innermost to outermost):
// 	- Metrics, using Prometheus.
// 	- Logging, using logrus/zap.
// 	- Recoverer, using grpc_recovery.
func NewDevelopmentGrpc(port string, logger *logrus.Entry) *Grpc {
	options := grpc_middleware.WithUnaryServerChain(defaultUnaryServerInterceptors(logger)...)

	srv := NewGrpc(port, options)
	grpc_prometheus.Register(srv.Server)
	return srv
}

// NewProductionGrpc creates an instance of Grpc with default production options attached.
// Actually, it can be used for non-production environment (such as staging or sandbox) as long as the environment satisfies all prerequisites.
//
// These are list of interceptors that are attached (from innermost to outermost):
// 	- Metrics, using Prometheus.
// 	- Logging, using logrus/zap.
// 	- Recoverer, using grpc_recovery.
// 	- Error Reporter, using Google Cloud Error Reporter.
//
// It also activates some auxiliaries:
// 	- Profiler, using Google Cloud Profiler.
// 	- Tracing, using Google Cloud Stackdriver Trace. The sample probability is 1% for production environment. Otherwise, it is 100%.
func NewProductionGrpc(cfg *config.Config, logger *logrus.Entry) (*Grpc, error) {
	midds := []grpc.UnaryServerInterceptor{}
	midds = append(midds, defaultUnaryServerInterceptors(logger)...)
	options := grpc_middleware.WithUnaryServerChain(midds...)

	srv := NewGrpc(cfg.Port.GRPC, grpc.StatsHandler(&ocgrpc.ServerHandler{}), options)
	grpc_prometheus.Register(srv.Server)

	return srv, nil
}

// Run runs the server.
// It basically runs grpc.Server.Serve and is a blocking.
func (g *Grpc) Run() error {
	var err error
	g.listener, err = net.Listen(connProtocol, fmt.Sprintf(":%s", g.port))
	if err != nil {
		return err
	}

	go g.serve()
	log.Printf("grpc server is running on port %s\n", g.port)
	return nil
}

// AwaitTermination blocks the server and wait for termination signal.
// The termination signal must be one of SIGINT or SIGTERM.
// Once it receives one of those signals, the gRPC server will perform graceful stop and close the listener.
func (g *Grpc) AwaitTermination() error {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	g.GracefulStop()
	return g.listener.Close()
}

func (g *Grpc) serve() {
	if err := g.Serve(g.listener); err != nil {
		panic(err)
	}
}

func defaultUnaryServerInterceptors(logger *logrus.Entry) []grpc.UnaryServerInterceptor {
	grpc_prometheus.EnableHandlingTimeHistogram()

	skipHealthCheckLog := []grpc_logrus.Option{
		grpc_logrus.WithDecider(func(methodFullName string, err error) bool {
			service := path.Dir(methodFullName)[1:]
			method := path.Base(methodFullName)
			// will not log gRPC calls if it was a call to healthcheck and no error was raised
			if err == nil && method == "Check" && service == "grpc.health.v1.Health" {
				return false
			}
			// by default you will log all calls
			return true
		}),
	}

	options := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryHandler)),
		grpc_logrus.UnaryServerInterceptor(logger, skipHealthCheckLog...),
		logging.UnaryServerInterceptor(false),
		grpc_prometheus.UnaryServerInterceptor,
	}
	return options
}

func recoveryHandler(p interface{}) error {
	return status.Errorf(codes.Unknown, "%v", p)
}
