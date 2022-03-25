package server_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ibrahimker/tigerhall-kittens/common/logging"
	"github.com/ibrahimker/tigerhall-kittens/server"
)

var (
	testPort = "8080"
)

func TestNewGrpc(t *testing.T) {
	t.Run("successfully create a gRPC server", func(t *testing.T) {
		srv := server.NewGrpc(testPort)
		assert.NotNil(t, srv)
	})
}

func TestNewDevelopmentGrpc(t *testing.T) {
	t.Run("successfully create a development gRPC server", func(t *testing.T) {
		srv := server.NewDevelopmentGrpc(testPort, nil)
		defer srv.Stop()
		assert.NotNil(t, srv)
	})
	t.Run("successfully create a development gRPC server with logrus", func(t *testing.T) {
		logger := logging.NewTestLogger()
		srv := server.NewDevelopmentGrpc(testPort, logger)
		defer srv.Stop()
		assert.NotNil(t, srv)
	})
}

func TestGrpc_Run(t *testing.T) {
	t.Run("listener fails", func(t *testing.T) {
		srv := server.NewGrpc("abc")

		err := srv.Run()
		defer srv.Stop()

		assert.NotNil(t, err)
	})

	t.Run("success run", func(t *testing.T) {
		srv := server.NewGrpc("8018")

		err := srv.Run()
		defer srv.Stop()
		time.Sleep(1 * time.Second)

		assert.Nil(t, err)
	})
}
