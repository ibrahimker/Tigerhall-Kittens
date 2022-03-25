package redis

import (
	"context"

	goredis "github.com/go-redis/redis/v8"
	nrredis "github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/ibrahimker/tigerhall-kittens/common/config"
)

// NewClient creates an instance of redis client.
func NewClient(cfg *config.Redis) (*goredis.Client, error) {
	opt := &goredis.Options{
		Addr: cfg.Address,
	}

	client := goredis.NewClient(opt)
	client.AddHook(nrredis.NewHook(opt))

	txn := newrelicTxn()
	ctx := newrelic.NewContext(context.Background(), txn)

	_, err := client.WithContext(ctx).Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newrelicTxn() *newrelic.Transaction { return nil }
