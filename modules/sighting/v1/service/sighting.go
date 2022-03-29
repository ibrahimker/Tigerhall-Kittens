// Package service defines the business logic of the requirement.
// The general flow of the requirements are explicitly stated in the code.
package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ibrahimker/tigerhall-kittens/common/logging"
	"github.com/ibrahimker/tigerhall-kittens/driver/redis"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/entity"
)

const (
	// BaseKey is the Base key for Cache
	BaseKey = entity.ModuleName + ":" + entity.ModuleVersion + ":"
	// GetTigersKey is a base key for caching GetTigers service
	GetTigersKey = BaseKey + "sighting:get-tigers"
)

var (
	// GetTigersRedisTTL set time needed for cache to expire
	GetTigersRedisTTL = 1 * time.Minute
)

// TigerSighting defines the interface to tiger sighting services.
type TigerSighting interface {
	// GetTigers get list of tigers from database order by last seen timestamp
	GetTigers(ctx context.Context) ([]*entity.Tiger, error)
	// CreateTiger store a new tiger in database
	CreateTiger(ctx context.Context, tiger *entity.Tiger) error
	// GetSightingsByTigerID get list of sightings for given tiger ID order by latest sighting
	GetSightingsByTigerID(ctx context.Context, tigerID int32) ([]*entity.Sighting, error)
	// CreateSighting store a new sighting for given tiger ID in database if not within 5 km of previous sighting
	CreateSighting(ctx context.Context, tigerID int32, sighting *entity.Sighting) error
}

// TigerSightingRepository defines the interface to tiger sighting repository.
type TigerSightingRepository interface {
	// GetTigers get list of tigers from database order by last seen timestamp
	GetTigers(ctx context.Context) ([]*entity.Tiger, error)
	// CreateTiger store a new tiger in database
	CreateTiger(ctx context.Context, tiger *entity.Tiger) error
	// GetSightingsByTigerID get list of sightings for given tiger ID order by latest sighting
	GetSightingsByTigerID(ctx context.Context, tigerID int32) ([]*entity.Sighting, error)
	// CreateSighting store a new sighting for given tiger ID in database
	CreateSighting(ctx context.Context, tigerID int32, sighting *entity.Sighting) error
}

// TigerSightingService is responsible for hold dependencies related to tiger sighting service.
type TigerSightingService struct {
	repo      TigerSightingRepository
	redisRepo redis.Redis
}

// NewTigerSightingService creates an instance of TigerSightingService.
func NewTigerSightingService(repo TigerSightingRepository, redisRepo redis.Redis) *TigerSightingService {
	return &TigerSightingService{
		repo:      repo,
		redisRepo: redisRepo,
	}
}

// GetTigers get list of tigers from database order by last seen timestamp
func (t *TigerSightingService) GetTigers(ctx context.Context) (tigers []*entity.Tiger, err error) {
	logger := logging.NewServiceLogger(ctx, "GetTigers", logrus.Fields{})

	// Get data cache from Redis, if data empty or not found then get tiger data from Database
	if err = t.redisRepo.Fetch(ctx, GetTigersKey, &tigers, GetTigersRedisTTL, func() (interface{}, error) {
		tigers, err = t.repo.GetTigers(ctx)
		if err != nil {
			logging.WithError(err, logger).Warn("Error when get from repo.GetTigers")
			return nil, err
		}
		return tigers, nil
	}); err != nil {
		logging.WithError(err, logger).Warn("Error when get from redisRepo.Fetch")
		return nil, err
	}

	return tigers, nil
}

// CreateTiger store a new tiger in database
func (t *TigerSightingService) CreateTiger(ctx context.Context, tiger *entity.Tiger) error {
	logger := logging.NewServiceLogger(ctx, "CreateTiger", logrus.Fields{})

	// validate input
	if err := isValidTiger(tiger); err != nil {
		logging.WithError(err, logger).Warn("Error when get from validate tiger")
		return err
	}

	// insert to repo
	if err := t.repo.CreateTiger(ctx, tiger); err != nil {
		logging.WithError(err, logger).Warn("Error when get from repo.CreateTiger")
		return err
	}
	return nil
}

// GetSightingsByTigerID get list of sightings for given tiger ID order by latest sighting
func (t *TigerSightingService) GetSightingsByTigerID(ctx context.Context, tigerID int32) ([]*entity.Sighting, error) {
	// TODO: implement me
	panic("implement me")
}

// CreateSighting store a new sighting for given tiger ID in database if not within 5 km of previous sighting
func (t *TigerSightingService) CreateSighting(ctx context.Context, tigerID int32, sighting *entity.Sighting) error {
	// TODO: implement me
	panic("implement me")
}
