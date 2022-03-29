package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"

	"github.com/ibrahimker/tigerhall-kittens/common/logging"
)

func queryWrapper(ctx context.Context, pool PgxPoolIface, queryString string, args ...interface{}) (pgx.Rows, error) {
	logger := logging.FromContext(ctx).WithFields(logrus.Fields{
		"sub-repo-name": "repo.queryWrapper",
		"queryString":   queryString,
		"args":          args,
	})

	rows, err := pool.Query(ctx, queryString, args...)
	if err != nil {
		logging.WithError(err, logger).Warn("Error when execute query")
		return nil, err
	}
	if err = rows.Err(); err != nil {
		logging.WithError(err, logger).Warn("Error when check rows err")
		return nil, err
	}

	return rows, nil
}
