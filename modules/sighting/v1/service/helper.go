package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
	"time"

	"github.com/nfnt/resize"

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

func isValidSighting(sighting *entity.Sighting) error {
	if sighting.TigerID == 0 {
		return errors.New("tiger id cannot be 0")
	}
	if validateTime(sighting.SeenAt) {
		return errors.New("seen_at cannot be zero")
	}
	if sighting.Latitude < -90.0 || sighting.Latitude > 90.0 {
		return errors.New("not a valid latitude")
	}
	if sighting.Longitude < -180.0 || sighting.Longitude > 180.0 {
		return errors.New("not a valid longitude")
	}
	if sighting.ImageData == "" {
		return errors.New("image data should contain valid base64 image format")
	}
	return nil
}

func resizeBase64Image(in string) (resizedImageBase64 string, err error) {
	const (
		jpegPrefix = "data:image/jpeg;base64,"
		pngPrefix  = "data:image/png;base64,"
	)
	coI := strings.Index(in, ",")
	var base64Image string
	if strings.HasPrefix(in, pngPrefix) || strings.HasPrefix(in, jpegPrefix) {
		base64Image = in[coI+1:]
	} else { // case no prefix
		base64Image = in
	}

	unbased, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", err
	}

	r := bytes.NewReader(unbased)
	var im image.Image
	if strings.HasPrefix(in, jpegPrefix) { // case jpeg
		im, err = jpeg.Decode(r)
		if err != nil {
			return "", err
		}
	} else { // default case we treat as png
		im, err = png.Decode(r)
		if err != nil {
			return "", err
		}
	}

	resizedImage := resize.Resize(250, 200, im, resize.Lanczos3)
	var resizedImageBuf bytes.Buffer

	if strings.HasPrefix(in, jpegPrefix) { // case jpeg
		if err = jpeg.Encode(&resizedImageBuf, resizedImage, &jpeg.Options{Quality: 80}); err != nil {
			return "", err
		}
		resizedImageBase64 = jpegPrefix
		resizedImageBase64 += base64.StdEncoding.EncodeToString(resizedImageBuf.Bytes())
	} else { // default case we treat as png
		if err = png.Encode(&resizedImageBuf, resizedImage); err != nil {
			return "", err
		}
		resizedImageBase64 = pngPrefix
		resizedImageBase64 += base64.StdEncoding.EncodeToString(resizedImageBuf.Bytes())
	}
	return resizedImageBase64, nil
}
