package builder_test

import (
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/ibrahimker/tigerhall-kittens/common/config"
	"github.com/ibrahimker/tigerhall-kittens/common/logging"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/internal/builder"
)

func TestBuildTigerSightingandler(t *testing.T) {
	t.Run("successfully build Tiger Sighting handler", func(t *testing.T) {
		pool := &pgxpool.Pool{}
		cfg, _ := config.NewConfig("../../../../../test/fixture/env.valid")
		rds, _ := redismock.NewClientMock()

		hdr := builder.BuildTigerSightingHandler(cfg, pool, rds, logging.NewTestLogger())

		assert.NotNil(t, hdr)
	})
}
