package service

import (
	"errors"
	"time"

	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/entity"
)

func validateTime(in time.Time) bool {
	return in.IsZero() || in.Equal(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
}
func isValidTiger(tiger *entity.Tiger) error {
	if tiger.Name == "" {
		return errors.New("name cannot be empty")
	}
	if validateTime(tiger.DateOfBirth) {
		return errors.New("time cannot be zero")
	}
	if validateTime(tiger.LastSeenTimestamp) {
		return errors.New("last seen time cannot be null")
	}
	if tiger.LastSeenLatitude < -90.0 || tiger.LastSeenLatitude > 90.0 {
		return errors.New("not a valid latitude")
	}
	if tiger.LastSeenLongitude < -180.0 || tiger.LastSeenLongitude > 180.0 {
		return errors.New("not a valid longitude")
	}
	return nil
}
