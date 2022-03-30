// Package service defines the business logic of the requirement.
// The general flow of the requirements are explicitly stated in the code.
package service

import (
	"context"
	"fmt"
	"time"

	geo "github.com/kellydunn/golang-geo"
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
	// GetSightingsByTigerIDKey is a base key for caching GetSightingsByTigerID service
	GetSightingsByTigerIDKey = BaseKey + "sighting:get-sightings-by-tiger:%d"
)

var (
	// GetTigersRedisTTL set time needed for cache to expire. This cache is to prevent db overload
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
	CreateSighting(ctx context.Context, sighting *entity.Sighting) error
}

// TigerSightingRepository defines the interface to tiger sighting repository.
type TigerSightingRepository interface {
	// GetTigers get list of tigers from database order by last seen timestamp
	GetTigers(ctx context.Context) ([]*entity.Tiger, error)
	// GetTigerByID get tiger by ID from database
	GetTigerByID(ctx context.Context, tigerID int32) (*entity.Tiger, error)
	// CreateTiger store a new tiger in database
	CreateTiger(ctx context.Context, tiger *entity.Tiger) error
	// UpdateTiger update tiger data in database
	UpdateTiger(ctx context.Context, tiger *entity.Tiger) error

	// GetSightingsByTigerID get list of sightings for given tiger ID order by latest sighting
	GetSightingsByTigerID(ctx context.Context, tigerID int32) ([]*entity.Sighting, error)
	// CreateSighting store a new sighting for given tiger ID in database
	CreateSighting(ctx context.Context, sighting *entity.Sighting) error
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

	// invalidate cache
	_ = t.redisRepo.Del(ctx, GetTigersKey)

	return nil
}

// GetSightingsByTigerID get list of sightings for given tiger ID order by latest sighting
func (t *TigerSightingService) GetSightingsByTigerID(ctx context.Context, tigerID int32) (sightings []*entity.Sighting, err error) {
	logger := logging.NewServiceLogger(ctx, "GetSightingsByTigerID", logrus.Fields{})

	// Get data cache from Redis, if data empty or not found then get tiger data from Database
	if err = t.redisRepo.Fetch(ctx, fmt.Sprintf(GetSightingsByTigerIDKey, tigerID), &sightings, GetTigersRedisTTL, func() (interface{}, error) {
		sightings, err = t.repo.GetSightingsByTigerID(ctx, tigerID)
		if err != nil {
			logging.WithError(err, logger).Warn("Error when get from repo.GetSightingsByTigerID")
			return nil, err
		}
		return sightings, nil
	}); err != nil {
		logging.WithError(err, logger).Warn("Error when get from redisRepo.Fetch")
		return nil, err
	}

	return sightings, nil
}

// CreateSighting store a new sighting for given tiger ID in database if not within 5 km of previous sighting
// It will also resize sighting image into 250x200
func (t *TigerSightingService) CreateSighting(ctx context.Context, sighting *entity.Sighting) error {
	logger := logging.NewServiceLogger(ctx, "CreateTiger", logrus.Fields{})

	// validate input
	if err := isValidSighting(sighting); err != nil {
		logging.WithError(err, logger).Warn("Error when get from validate sighting")
		return err
	}

	// validate is new lat/long in 5km radius
	tiger, err := t.repo.GetTigerByID(ctx, sighting.TigerID)
	if err != nil {
		logging.WithError(err, logger).Warn("Error when get repo.GetSightingsByTigerID")
		return err
	}
	dist := geo.NewPoint(tiger.LastSeenLatitude, tiger.LastSeenLongitude).GreatCircleDistance(geo.NewPoint(sighting.Latitude, sighting.Longitude))
	if dist > 5.00 {
		err = fmt.Errorf("distance exceed 5000. Distance: %.2f", dist)
		logging.WithError(err, logger).Warn("Error when get validate distance")
		return err
	}

	// resize image into 250x200
	resizedBase64, err := resizeBase64Image(sighting.ImageData)
	if err != nil {
		logging.WithError(err, logger).Warn("Error when get resizeBase64Image")
		return err
	}
	sighting.ImageData = resizedBase64

	// insert to repo
	if err = t.repo.CreateSighting(ctx, sighting); err != nil {
		logging.WithError(err, logger).Warn("Error when get from repo.CreateSighting")
		return err
	}

	// update tiger data
	tiger.LastSeenTimestamp = sighting.SeenAt
	tiger.LastSeenLatitude = sighting.Latitude
	tiger.LastSeenLongitude = sighting.Longitude
	if err = t.repo.UpdateTiger(ctx, tiger); err != nil {
		logging.WithError(err, logger).Warn("Error when get from repo.UpdateTiger")
		return err
	}

	// invalidate cache
	_ = t.redisRepo.Del(ctx, fmt.Sprintf(GetSightingsByTigerIDKey, sighting.TigerID))
	_ = t.redisRepo.Del(ctx, GetTigersKey)

	return nil
}
