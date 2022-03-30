package handler

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	tigerv1 "github.com/ibrahimker/tigerhall-kittens/api/proto"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/entity"
)

func composeTigersProto(req []*entity.Tiger) (res []*tigerv1.Tiger) {
	for _, v := range req {
		res = append(res, &tigerv1.Tiger{
			Id:                v.ID,
			Name:              v.Name,
			DateOfBirth:       timestamppb.New(v.DateOfBirth),
			LastSeenTimestamp: timestamppb.New(v.LastSeenTimestamp),
			LastSeenLatitude:  wrapperspb.Double(v.LastSeenLatitude),
			LastSeenLongitude: wrapperspb.Double(v.LastSeenLongitude),
			CreatedAt:         timestamppb.New(v.CreatedAt.Time),
			UpdatedAt:         timestamppb.New(v.UpdatedAt.Time),
		})
	}
	return res
}

func composeSightingsProto(req []*entity.Sighting) (res []*tigerv1.Sighting) {
	for _, v := range req {
		res = append(res, &tigerv1.Sighting{
			Id:        v.ID,
			SeenAt:    timestamppb.New(v.SeenAt),
			Latitude:  wrapperspb.Double(v.Latitude),
			Longitude: wrapperspb.Double(v.Longitude),
			ImageData: v.ImageData,
		})
	}
	return res
}
