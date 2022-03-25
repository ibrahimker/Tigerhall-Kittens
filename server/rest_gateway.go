package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/ibrahimker/tigerhall-kittens/common/config"
)

var loadtestHeader = "X-Loadtest-Email"
var maxGRPCErrCodeNumber = 16

// Rest is responsible to act as HTTP/1.1 REST server.
// It composes grpc-gateway runtime.ServeMux.
type Rest struct {
	*runtime.ServeMux
	port string
}

// NewRest creates an instance of Rest.
func NewRest(port string) *Rest {
	return &Rest{
		ServeMux: runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(MatcherLoadtestHeader),
		),
		port: port,
	}
}

// NewProductionRest creates an instance of Rest with default production options attached.
// The only difference between NewRest and NewProductionRest is the later enable Prometheus metrics by default.
func NewProductionRest(port string, cfg *config.Config) *Rest {
	srv := &Rest{
		ServeMux: runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(MatcherLoadtestHeader),
		),
		port: port,
	}

	_ = srv.EnableHealth() // error is impossible, hence ignored.
	return srv
}

// EnableHealth enables health endpoint.
// It can be accessed via /healthz.
func (r *Rest) EnableHealth() error {
	return r.ServeMux.HandlePath(http.MethodGet, "/healthz", healthHandler())
}

// Run runs HTTP/1.1 runtime.ServeMux.
// It runs inside a goroutine.
func (r *Rest) Run() error {
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", r.port), allowCORS(r.ServeMux)); err != nil {
			panic(err)
		}
	}()
	return nil
}

func healthHandler() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				headers := []string{"Content-Type", "Accept", "Authorization", loadtestHeader}
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
				methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// MatcherLoadtestHeader is used to matching custom header for loadtest
func MatcherLoadtestHeader(key string) (string, bool) {
	switch key {
	case loadtestHeader:
		return strings.ToLower(loadtestHeader), true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
