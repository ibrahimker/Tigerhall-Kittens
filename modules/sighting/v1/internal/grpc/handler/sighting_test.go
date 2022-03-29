package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	tigerv1 "github.com/ibrahimker/tigerhall-kittens/api/proto"
	"github.com/ibrahimker/tigerhall-kittens/common/logging"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/entity"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/internal/grpc/handler"
	mock_service "github.com/ibrahimker/tigerhall-kittens/test/mock/modules/sighting/v1/service"
)

type SightingTestSuite struct {
	logger          *logrus.Entry
	sightingSvc     *mock_service.MockTigerSighting
	sightingHandler *handler.TigerSighting
}

type HandlerTestCase struct {
	testcaseName     string
	testcaseFunction func(t *testing.T)
}

func TigerSightingServiceTestSuite(ctrl *gomock.Controller) *SightingTestSuite {
	logger := logging.NewTestLogger()
	mockSightingSvc := mock_service.NewMockTigerSighting(ctrl)
	sightingHandler := handler.NewTigerSighting(logger, mockSightingSvc)

	return &SightingTestSuite{
		logger:          logger,
		sightingSvc:     mockSightingSvc,
		sightingHandler: sightingHandler,
	}
}

func TestHelpCenterService_TestService(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	testCases := []HandlerTestCase{
		{
			testcaseName: "successfully get the data from redis",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				serviceSuite := TigerSightingServiceTestSuite(mockCtrl)

				require.NotNil(t, serviceSuite.sightingSvc)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestHelpCenterService_GetTigers(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCtx := context.Background()
	testCases := []HandlerTestCase{
		{
			testcaseName: "Error when hit service",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				serviceSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceSuite.sightingSvc.EXPECT().GetTigers(gomock.Any()).Return(nil, errors.New("db error"))

				resData, resErr := serviceSuite.sightingHandler.GetTigers(mockCtx, &tigerv1.GetTigersRequest{})
				require.Error(t, resErr)
				require.Nil(t, resData)
			},
		},
		{
			testcaseName: "Successfully hit service",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				serviceSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceSuite.sightingSvc.EXPECT().GetTigers(gomock.Any()).Return([]*entity.Tiger{{ID: 1}}, nil)

				resData, resErr := serviceSuite.sightingHandler.GetTigers(mockCtx, &tigerv1.GetTigersRequest{})
				require.Nil(t, resErr)
				require.Equal(t, 1, len(resData.Data))
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestHelpCenterService_CreateTiger(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	now := time.Now().UTC()
	tigerProtoData := &tigerv1.CreateTigerRequest{
		Name:              "tiger-1",
		DateOfBirth:       timestamppb.New(now),
		LastSeenTimestamp: timestamppb.New(now),
		LastSeenLatitude:  wrapperspb.Double(-6.18),
		LastSeenLongitude: wrapperspb.Double(108.00),
	}
	tigerData := &entity.Tiger{
		Name:              "tiger-1",
		DateOfBirth:       now,
		LastSeenTimestamp: now,
		LastSeenLatitude:  -6.18,
		LastSeenLongitude: 108.00,
	}
	mockCtx := context.Background()
	testCases := []HandlerTestCase{
		{
			testcaseName: "Error when hit service",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				serviceSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceSuite.sightingSvc.EXPECT().CreateTiger(gomock.Any(), tigerData).Return(errors.New("db error"))

				resData, resErr := serviceSuite.sightingHandler.CreateTiger(mockCtx, tigerProtoData)
				require.Error(t, resErr)
				require.Nil(t, resData)
			},
		},
		{
			testcaseName: "Successfully hit service",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				serviceSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceSuite.sightingSvc.EXPECT().CreateTiger(gomock.Any(), tigerData).Return(nil)

				_, resErr := serviceSuite.sightingHandler.CreateTiger(mockCtx, tigerProtoData)
				require.Nil(t, resErr)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestHelpCenterService_GetSightings(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tigerID := int32(1)

	mockCtx := context.Background()
	testCases := []HandlerTestCase{
		{
			testcaseName: "Error when hit service",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				serviceSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceSuite.sightingSvc.EXPECT().GetSightingsByTigerID(gomock.Any(), tigerID).Return(nil, errors.New("db error"))

				resData, resErr := serviceSuite.sightingHandler.GetSightings(mockCtx, &tigerv1.GetSightingsRequest{Id: tigerID})
				require.Error(t, resErr)
				require.Nil(t, resData)
			},
		},
		{
			testcaseName: "Successfully hit service",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				serviceSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceSuite.sightingSvc.EXPECT().GetSightingsByTigerID(gomock.Any(), tigerID).Return([]*entity.Sighting{{ID: 1}}, nil)

				resData, resErr := serviceSuite.sightingHandler.GetSightings(mockCtx, &tigerv1.GetSightingsRequest{Id: tigerID})
				require.Nil(t, resErr)
				require.Equal(t, 1, len(resData.Data))
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}
