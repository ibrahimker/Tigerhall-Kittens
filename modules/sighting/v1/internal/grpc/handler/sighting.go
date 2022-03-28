// Package handler provides the functionality of HTTP/2 gRPC handler.
// It receives request and returns response.
package handler

import (
	"context"

	"github.com/sirupsen/logrus"

	tigerv1 "github.com/ibrahimker/tigerhall-kittens/api/proto"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/service"
)

// Sighting handles HTTP/2 gRPC request for tiger sighting services.
type Sighting struct {
	tigerv1.UnimplementedTigerSightingServiceServer
	logger      *logrus.Entry
	sightingSvc service.TigerSighting
}

// NewSighting creates an instance of Sighting.
func NewSighting(logger *logrus.Entry, sightingSvc service.TigerSighting) *Sighting {
	return &Sighting{
		logger:      logger,
		sightingSvc: sightingSvc,
	}
}

// GetTigers handles HTTP/2 gRPC request similar to GET in HTTP/1.1.
func (s *Sighting) GetTigers(ctx context.Context, req *tigerv1.GetTigersRequest) (*tigerv1.GetTigersResponse, error) {
	return &tigerv1.GetTigersResponse{}, nil
}

// CreateTiger handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (s *Sighting) CreateTiger(ctx context.Context, req *tigerv1.CreateTigerRequest) (*tigerv1.CreateTigerResponse, error) {
	return &tigerv1.CreateTigerResponse{}, nil
}

// GetSightings handles HTTP/2 gRPC request similar to GET in HTTP/1.1.
func (s *Sighting) GetSightings(ctx context.Context, req *tigerv1.GetSightingsRequest) (*tigerv1.GetSightingsResponse, error) {
	return &tigerv1.GetSightingsResponse{}, nil
}

// CreateSighting handles HTTP/2 gRPC request similar to POST in HTTP/1.1.
func (s *Sighting) CreateSighting(ctx context.Context, req *tigerv1.CreateSightingRequest) (*tigerv1.CreateSightingResponse, error) {
	return &tigerv1.CreateSightingResponse{}, nil
}
