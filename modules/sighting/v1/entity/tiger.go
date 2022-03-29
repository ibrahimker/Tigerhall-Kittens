// Package entity defines all structs that flow through all application layers.
package entity

import (
	"database/sql"
	"time"
)

const (
	// ModuleVersion defines current sighting module version
	ModuleVersion = "v1"
	// ModuleName defines current sighting module name
	ModuleName = "sighting"
)

// Tiger is a struct to model tiger data
// we use float64 in lat/long because we don't need to calculate the distance so precise
type Tiger struct {
	ID                int32
	Name              string
	DateOfBirth       time.Time
	LastSeenTimestamp time.Time
	LastSeenLatitude  float64
	LastSeenLongitude float64
	CreatedAt         sql.NullTime
	UpdatedAt         sql.NullTime
}

// Sighting is a struct to model sighting of tiger data
// we use float64 in lat/long because we don't need to calculate the distance so precise
type Sighting struct {
	ID        int32
	TigerID   int32
	SeenAt    time.Time
	Latitude  float64
	Longitude float64
	ImageURL  string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
