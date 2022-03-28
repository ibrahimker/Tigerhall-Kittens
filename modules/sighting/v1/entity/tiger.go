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
type Tiger struct {
	ID                int32        `json:"id"`
	Name              string       `json:"name"`
	LastSeenTimestamp time.Time    `json:"last_seen_timestamp"`
	LastSeenLatitude  float64      `json:"last_seen_latitude"`
	LastSeenLongitude float64      `json:"last_seen_ongitude"`
	CreatedAt         sql.NullTime `json:"created_at"`
	UpdatedAt         sql.NullTime `json:"updated_at"`
}

// Sighting is a struct to model sighting of tiger data
type Sighting struct {
	ID        int32     `json:"id"`
	TigerID   int32     `json:"tiger_id"`
	SeenAt    time.Time `json:"seen_at"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	ImageURL  string    `json:"image_url"`
}
