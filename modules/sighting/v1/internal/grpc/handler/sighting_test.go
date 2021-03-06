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

func TestHelpCenterService_CreateSightings(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	now := time.Now().UTC()
	tigerID := int32(1)
	imageData := "data:image/png;base64,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lwHEKdtiJTl7FFFwhhIDBSWQ+wE2GwHZtHEjtaamqNVqqQQZGxuj0dFRGsztW1RB0x+TgQAIuCCQaws9eQERp/A4LLtYbxgDBLIl0K9DkuLWkaSAuGxRZZulcAwEQEA8gagH4BmimQGiAQtKgQWkvsbmWgJaxA6PggAIZEJA/P7WZ+sOLCDVox1VoaubiydAAARAwIpASh0WkQKSEkCrTCGi+jWZ7cxSngcBKZGAHfII1D5AB1pWYgREfAknL7eysShQrmfDC44kTqBGwkvdH+0EpAaIhaE3VlgHcyWedv7MB1t/bDEyCDgmIKlDYycgNcFIAlDTBTwmnQBEUXqEDOxDEHWQjA/guoFqXg8mIFJLsJrc8BgIgIBXAmWIxx9//EG77rorHXPMMXT//fcrovzD0Oeffz699dZbtPnmm9O5555Lt9xyCw0MDNCmTZvommuuoccee4z+/fdfOvzww+nhhx+mFStW0JTmB6nfe+89OuOMM2jvvfem5557rh292267ja699lpaunRp++9OPvlkevzxx9X/P/HEE3TDDTcou3baaSd66KGHaI899lDXvAtIbIX0muMYPCyBMvaUsEwxW1QCZ599Nr377rt0xBFHtAXkpJNOouXLl9N9991HvH8edthhdOWVV9J5551H99xzDz344IP0+uuv01ZbbUXnnHOO+m0bTz75ZNuPXh2eV155he68807aeeedad26dfME5Prrr6fffvuN7r333i4WX3zxBR188MH06quv0r777qvmv/322+nbb7+lzTbbzI+AoEUVNScxOQiAQAIEXn75Zbr55pvp6KOPpp9//lkJyF9//UUjIyP0ww8/0DbbbKO84BP/008/rSqSAw88kC666CI6/fTT1bXvvvuOdtttN+JKhsWEKxEWHv6ceeaZtGHDBrrxxhvp448/pn322UdVFV999dU8Abn44ouVGHGVs/DD4vLLL7+oKmf2w3Y99dRTdMghh7gTELSoDDIWJ2gDSLhljgASJtds4BM/b+hcGTz77LP0008/KQH5/PPPlUisXbu27TpXKCeeeKLayFeuXElvvPGGakPxZ2ZmRrW5+Lmtt96adt99d3rmmWdoenqauLrhv2dRmd2fb7rpJvrmm29UxTL7uwJPO+00ajabqgrhNth+++2nqpVtt92WTjjhBGUPV0CzH26bsT0sZNYtLLSock1x+AUCIOCLAFcQe+21F11xxRXEm/qsgHzwwQd03HHH0Zo1a9pTf/LJJ+q0/+eff9IWW2xBn376Ke2yyy7t68PDw6o6YUF67bXX6Oqrr1bvR7jdxO2vzg/PxePdfffdSiz4WRYw7hpdfvnlqi116aWXqjn4z1FHHaXs4Spl9nPssccqe6666qp6AoIWla+0wrggAAK5E3j++efpjjvuUJv+kiVL5gkIv3Pg6uKff/5pY3jzzTdVy4pP+1w1vPjii3TAAQeo69yi4grk66+/ph133FH9Hb+U5/cin332WRdKFpDx8fF2C6tXAcD7O8/z448/0mWXXabsue6669pjHXTQQcqeCy+80FxA5rWoZojG/pPhr0tPumOQtPG57xnwDwTaBE455RR655131MbPH35/wULAL6lfeukl9Q6EhWSHHXZQ17mdxO9LuHXFmzd/k+qCCy5Q11gMuMXE7Sce74EHHmi3sM466yz1ba6FFUingPC1Dz/8kLbffnvVDuNvWvEfrjC4rcXvU77//nv1DoY/rAOrVq2iF154Qc2rbWGhRYXMBwEQAAF/BDpbWDwLv5Pgr+s+8sgjqpXF7xz4nlNPPVV9A4urFxYTfvHNlQC/+2Dh4Bfv+++/P3Eb7O+//1Zi89FHH9F2223XNn5hBcIX+CU+//tJPDZ/VZjFid+T8Mt7fu/C3wpjAeOx+Su//LXeL7/8UlVPPQUELSp/yYKRQQAEQGBhVTD7DoT/nqsJ/soui8TQ0BBdcsklxN+G4g9XCfzfvNnzew5+H8Ev37fccks69NBD6fjjj1fvVfjD377i9he3yo488kh6//33VQXB4sTvOrjlxW2uyclJ9Y7j7bffVq0v/truXXfdRY1GQ32NmL+5deutt9Kvv/5Ke+65Jz366KPtdllbQLgs0f0gCsKeKAF0txINHMwGARkE+n3LdvXq1TQwPj4+wzewivEPsHT+RKKZ+fxPkvAuhQ8IgAAIgEDOBPh9Df8wIv/MClcs/wXIZT3M35g24AAAAABJRU5ErkJggg=="
	sightingProtoData := &tigerv1.CreateSightingRequest{
		Id:        tigerID,
		SeenAt:    timestamppb.New(now),
		Latitude:  wrapperspb.Double(-6.18),
		Longitude: wrapperspb.Double(108.00),
		ImageData: imageData,
	}
	sightingData := &entity.Sighting{
		TigerID:   sightingProtoData.GetId(),
		SeenAt:    sightingProtoData.GetSeenAt().AsTime(),
		Latitude:  sightingProtoData.GetLatitude().GetValue(),
		Longitude: sightingProtoData.GetLongitude().GetValue(),
		ImageData: sightingProtoData.GetImageData(),
	}
	mockCtx := context.Background()
	testCases := []HandlerTestCase{
		{
			testcaseName: "Error when hit service",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				serviceSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceSuite.sightingSvc.EXPECT().CreateSighting(gomock.Any(), sightingData).Return(errors.New("db error"))

				resData, resErr := serviceSuite.sightingHandler.CreateSighting(mockCtx, sightingProtoData)
				require.Error(t, resErr)
				require.Nil(t, resData)
			},
		},
		{
			testcaseName: "Successfully hit service",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				serviceSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceSuite.sightingSvc.EXPECT().CreateSighting(gomock.Any(), sightingData).Return(nil)

				_, resErr := serviceSuite.sightingHandler.CreateSighting(mockCtx, sightingProtoData)
				require.Nil(t, resErr)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}
