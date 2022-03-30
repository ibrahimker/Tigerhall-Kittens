// Package postgres provides real connection to the PostgreSQL.
package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"

	"github.com/ibrahimker/tigerhall-kittens/common/logging"
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
	logger := logging.NewRepoLogger(ctx, "GetTigers", logrus.Fields{})

	queryString := `SELECT id,name,date_of_birth,last_seen_timestamp,last_seen_latitude,last_seen_longitude,created_at,updated_at
FROM sighting.tiger WHERE deleted_at IS NULL ORDER BY last_seen_timestamp desc`
	rows, err := queryWrapper(ctx, t.pool, queryString)
	if err != nil {
		logging.WithError(err, logger).Warn("Error when hit query wrapper")
		return []*entity.Tiger{}, err
	}
	defer rows.Close()

	var res []*entity.Tiger
	for rows.Next() {
		var tmp entity.Tiger
		if serr := rows.Scan(
			&tmp.ID, &tmp.Name, &tmp.DateOfBirth, &tmp.LastSeenTimestamp, &tmp.LastSeenLatitude, &tmp.LastSeenLongitude,
			&tmp.CreatedAt, &tmp.UpdatedAt,
		); serr != nil {
			logging.WithError(serr, logger).Warn("Error when scan rows")
			continue
		}
		res = append(res, &tmp)
	}
	if rows.Err() != nil {
		logging.WithError(rows.Err(), logger).Warn("Error when check rows")
		return []*entity.Tiger{}, rows.Err()
	}

	return res, nil
}

// GetTigerByID get tiger by ID from database
func (t *TigerSightingRepo) GetTigerByID(ctx context.Context, tigerID int32) (*entity.Tiger, error) {
	logger := logging.NewRepoLogger(ctx, "GetTigerByID", logrus.Fields{})

	queryString := `SELECT id,name,date_of_birth,last_seen_timestamp,last_seen_latitude,last_seen_longitude,created_at,updated_at
FROM sighting.tiger WHERE id = $1 and deleted_at IS NULL`
	rows, err := queryWrapper(ctx, t.pool, queryString, tigerID)
	if err != nil {
		logging.WithError(err, logger).Warn("Error when hit query wrapper")
		return nil, err
	}
	defer rows.Close()

	var res entity.Tiger
	for rows.Next() {
		if serr := rows.Scan(
			&res.ID, &res.Name, &res.DateOfBirth, &res.LastSeenTimestamp, &res.LastSeenLatitude, &res.LastSeenLongitude,
			&res.CreatedAt, &res.UpdatedAt,
		); serr != nil {
			logging.WithError(serr, logger).Warn("Error when scan rows")
			continue
		}
	}
	if rows.Err() != nil {
		logging.WithError(rows.Err(), logger).Warn("Error when check rows")
		return nil, rows.Err()
	}

	return &res, nil
}

// CreateTiger store a new tiger in database
func (t *TigerSightingRepo) CreateTiger(ctx context.Context, tiger *entity.Tiger) error {
	logger := logging.NewRepoLogger(ctx, "CreateTiger", logrus.Fields{})

	queryString := "INSERT INTO sighting.tiger" +
		" (name,date_of_birth,last_seen_timestamp,last_seen_latitude,last_seen_longitude,created_at,updated_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"

	currentTime := time.Now()
	_, err := t.pool.Exec(ctx, queryString,
		tiger.Name,
		tiger.DateOfBirth,
		tiger.LastSeenTimestamp,
		tiger.LastSeenLatitude,
		tiger.LastSeenLongitude,
		currentTime,
		currentTime,
	)
	if err != nil {
		logging.WithError(err, logger).Warnf("Error when execute query %s", queryString)
	}

	return err
}

// UpdateTiger update tiger data in database
func (t *TigerSightingRepo) UpdateTiger(ctx context.Context, tiger *entity.Tiger) error {
	logger := logging.NewRepoLogger(ctx, "UpdateTiger", logrus.Fields{})

	queryString := "UPDATE sighting.tiger " +
		"SET last_seen_timestamp = $2, last_seen_latitude = $3, last_seen_longitude = $4, updated_at = $5 " +
		"WHERE id = $1"

	_, err := t.pool.Exec(ctx, queryString, tiger.ID, tiger.LastSeenTimestamp, tiger.LastSeenLatitude, tiger.LastSeenLongitude, time.Now())
	if err != nil {
		logging.WithError(err, logger).Warnf("Error when execute query %s", queryString)
	}

	return err
}

// GetSightingsByTigerID get list of sightings for given tiger ID order by latest sighting
func (t *TigerSightingRepo) GetSightingsByTigerID(ctx context.Context, tigerID int32) ([]*entity.Sighting, error) {
	logger := logging.NewRepoLogger(ctx, "GetSightingsByTigerID", logrus.Fields{})

	queryString := `SELECT id,tiger_id,seen_at,latitude,longitude,image_data,created_at,updated_at
FROM sighting.sighting WHERE tiger_id = $1 and deleted_at IS NULL ORDER BY seen_at desc`
	rows, err := queryWrapper(ctx, t.pool, queryString, tigerID)
	if err != nil {
		logging.WithError(err, logger).Warn("Error when hit query wrapper")
		return []*entity.Sighting{}, err
	}
	defer rows.Close()

	var res []*entity.Sighting
	for rows.Next() {
		var tmp entity.Sighting
		if serr := rows.Scan(
			&tmp.ID, &tmp.TigerID, &tmp.SeenAt, &tmp.Latitude, &tmp.Longitude, &tmp.ImageData,
			&tmp.CreatedAt, &tmp.UpdatedAt,
		); serr != nil {
			logging.WithError(serr, logger).Warn("Error when scan rows")
			continue
		}
		res = append(res, &tmp)
	}
	if rows.Err() != nil {
		logging.WithError(rows.Err(), logger).Warn("Error when check rows")
		return []*entity.Sighting{}, rows.Err()
	}

	return res, nil
}

// CreateSighting store a new sighting for given tiger ID in database
func (t *TigerSightingRepo) CreateSighting(ctx context.Context, sighting *entity.Sighting) error {
	logger := logging.NewRepoLogger(ctx, "CreateTiger", logrus.Fields{})

	queryString := "INSERT INTO sighting.sighting" +
		" (tiger_id,seen_at,latitude,longitude,image_data,created_at,updated_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"

	currentTime := time.Now()
	_, err := t.pool.Exec(ctx, queryString,
		sighting.TigerID,
		sighting.SeenAt,
		sighting.Latitude,
		sighting.Longitude,
		sighting.ImageData,
		currentTime,
		currentTime,
	)
	if err != nil {
		logging.WithError(err, logger).Warnf("Error when execute query %s", queryString)
	}

	return err
}
