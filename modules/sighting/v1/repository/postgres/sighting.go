// Package postgres provides real connection to the PostgreSQL.
package postgres

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/entity"
)

// PgxPoolIface defines a little interface for pgxpool functionality.
// Since in the real implementation we can use pgxpool.Pool,
// this interface exists mostly for testing purpose.
type PgxPoolIface interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Ping(ctx context.Context) error
}

// TigerSightingRepo is responsible to connect tiger sighting entity with tiger sighting related table in PostgreSQL.
type TigerSightingRepo struct {
	pool PgxPoolIface
}

// NewTigerSightingRepo creates an instance of TigerSightingRepo.
func NewTigerSightingRepo(pool PgxPoolIface) *TigerSightingRepo {
	return &TigerSightingRepo{pool: pool}
}

// GetTigers get list of tigers from database order by last seen timestamp
func (t *TigerSightingRepo) GetTigers(ctx context.Context) ([]*entity.Tiger, error) {
	// TODO: implement me
	panic("implement me")
}

// CreateTiger store a new tiger in database
func (t *TigerSightingRepo) CreateTiger(ctx context.Context, tiger *entity.Tiger) error {
	// TODO: implement me
	panic("implement me")
}

// GetSightingsByTigerID get list of sightings for given tiger ID order by latest sighting
func (t *TigerSightingRepo) GetSightingsByTigerID(ctx context.Context, tigerID int32) ([]*entity.Sighting, error) {
	// TODO: implement me
	panic("implement me")
}

// CreateSighting store a new sighting for given tiger ID in database
func (t *TigerSightingRepo) CreateSighting(ctx context.Context, tigerID int32, sighting *entity.Sighting) error {
	// TODO: implement me
	panic("implement me")
}
