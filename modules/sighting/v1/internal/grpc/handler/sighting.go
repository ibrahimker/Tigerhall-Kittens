// Package handler provides the functionality of HTTP/2 gRPC handler.
// It receives request and returns response.
package handler

import (
	"context"

	"github.com/sirupsen/logrus"

	tigerv1 "github.com/ibrahimker/tigerhall-kittens/api/proto"
	"github.com/ibrahimker/tigerhall-kittens/common/logging"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/entity"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/service"
)

// TigerSighting handles HTTP/2 gRPC request for tiger sighting services.
type TigerSighting struct {
	tigerv1.UnimplementedTigerSightingServiceServer
	logger      *logrus.Entry
	sightingSvc service.TigerSighting
}

// NewTigerSighting creates an instance of TigerSighting.
func NewTigerSighting(logger *logrus.Entry, sightingSvc service.TigerSighting) *TigerSighting {
	return &TigerSighting{
		logger:      logger,
		sightingSvc: sightingSvc,
	}
}

// GetTigers handles HTTP/2 gRPC request similar to GET in HTTP/1.1.
func (s *TigerSighting) GetTigers(ctx context.Context, req *tigerv1.GetTigersRequest) (*tigerv1.GetTigersResponse, error) {
	logger, ctx := logging.NewHandlerLogger(ctx, s.logger, "GetTigers", req)

	data, err := s.sightingSvc.GetTigers(ctx)
	if err != nil {
		logging.WithError(err, logger).Error("Error when call s.sightingSvc.GetTigers")
		return nil, err
	}

	res := &tigerv1.GetTigersResponse{
		Data: composeTigersProto(data),
	}
	return res, nil
}

// CreateTiger handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (s *TigerSighting) CreateTiger(ctx context.Context, req *tigerv1.CreateTigerRequest) (*tigerv1.CreateTigerResponse, error) {
	logger, ctx := logging.NewHandlerLogger(ctx, s.logger, "GetTigers", req)

	if err := s.sightingSvc.CreateTiger(ctx, &entity.Tiger{
		Name:              req.GetName(),
		DateOfBirth:       req.GetDateOfBirth().AsTime(),
		LastSeenTimestamp: req.GetLastSeenTimestamp().AsTime(),
		LastSeenLatitude:  req.GetLastSeenLatitude().GetValue(),
		LastSeenLongitude: req.GetLastSeenLongitude().GetValue(),
	}); err != nil {
		logging.WithError(err, logger).Error("Error when call s.sightingSvc.CreateTiger")
		return nil, err
	}

	res := &tigerv1.CreateTigerResponse{
		Message: "Successfully create new tiger",
	}
	return res, nil
}

// GetSightings handles HTTP/2 gRPC request similar to GET in HTTP/1.1.
func (s *TigerSighting) GetSightings(ctx context.Context, req *tigerv1.GetSightingsRequest) (*tigerv1.GetSightingsResponse, error) {
	logger, ctx := logging.NewHandlerLogger(ctx, s.logger, "GetTigers", req)

	data, err := s.sightingSvc.GetSightingsByTigerID(ctx, req.GetId())
	if err != nil {
		logging.WithError(err, logger).Error("Error when call s.sightingSvc.GetSightingsByTigerID")
		return nil, err
	}

	res := &tigerv1.GetSightingsResponse{
		Data: composeSightingsProto(data),
	}
	return res, nil
}

// CreateSighting handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (s *TigerSighting) CreateSighting(ctx context.Context, req *tigerv1.CreateSightingRequest) (*tigerv1.CreateSightingResponse, error) {
	logger, ctx := logging.NewHandlerLogger(ctx, s.logger, "CreateSighting", req)

	if err := s.sightingSvc.CreateSighting(ctx, &entity.Sighting{
		TigerID:   req.GetId(),
		SeenAt:    req.GetSeenAt().AsTime(),
		Latitude:  req.GetLatitude().GetValue(),
		Longitude: req.GetLongitude().GetValue(),
		ImageData: req.GetImageData(),
	}); err != nil {
		logging.WithError(err, logger).Error("Error when call s.sightingSvc.CreateSighting")
		return nil, err
	}

	res := &tigerv1.CreateSightingResponse{
		Message: "Successfully create new sighting",
	}
	return res, nil
}
