package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/ibrahimker/tigerhall-kittens/common/logging"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/entity"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/service"
	mockRedisRepo "github.com/ibrahimker/tigerhall-kittens/test/mock/common/redis"
	mockRepo "github.com/ibrahimker/tigerhall-kittens/test/mock/modules/sighting/v1/service"
)

type SightingTestSuite struct {
	logger       *logrus.Entry
	sightingSvc  service.TigerSighting
	redisRepo    *mockRedisRepo.MockRedis
	sightingRepo *mockRepo.MockTigerSightingRepository
}

type ServiceTestCase struct {
	testcaseName     string
	testcaseFunction func(t *testing.T)
}

func TigerSightingServiceTestSuite(ctrl *gomock.Controller) *SightingTestSuite {
	logger := logging.NewTestLogger()
	mockRedisRepo := mockRedisRepo.NewMockRedis(ctrl)
	mockSightingRepo := mockRepo.NewMockTigerSightingRepository(ctrl)

	tigerSightingService := service.NewTigerSightingService(mockSightingRepo, mockRedisRepo)

	return &SightingTestSuite{
		logger:       logger,
		sightingSvc:  tigerSightingService,
		redisRepo:    mockRedisRepo,
		sightingRepo: mockSightingRepo,
	}
}

func TestHelpCenterService_TestService(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Success create service instance", func(t *testing.T) {
		t.Parallel()
		serviceSuite := TigerSightingServiceTestSuite(mockCtrl)

		require.NotNil(t, serviceSuite.sightingSvc)
		require.NotNil(t, serviceSuite.sightingRepo)
		require.NotNil(t, serviceSuite.redisRepo)
	})
}

func TestGetTigers(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCtx := context.Background()
	var emptyTiger []*entity.Tiger
	tigerData := []*entity.Tiger{{ID: 1}}
	mockTTL := 1 * time.Minute

	testCases := []ServiceTestCase{
		{
			testcaseName: "successfully get the data from redis",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceTestSuite.redisRepo.EXPECT().Fetch(mockCtx, service.GetTigersKey, &emptyTiger, mockTTL, gomock.Any()).
					SetArg(2, tigerData).Return(nil)

				resData, resErr := serviceTestSuite.sightingSvc.GetTigers(mockCtx)

				require.NoError(t, resErr)
				require.NotNil(t, resData)
				require.Equal(t, tigerData, resData)
			},
		},
		{
			testcaseName: "Error when retrieve from database",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				callbackFunc := func(ctx context.Context, key string, anyTigers *[]*entity.Tiger, ttl time.Duration, callback func() (interface{}, error)) {
					serviceTestSuite.sightingRepo.EXPECT().GetTigers(mockCtx).Return(nil, errors.New("db error"))
					_, _ = callback()
				}

				serviceTestSuite.redisRepo.EXPECT().Fetch(mockCtx, service.GetTigersKey, &emptyTiger, mockTTL, gomock.Any()).Do(callbackFunc).Return(errors.New("db error"))

				resData, resErr := serviceTestSuite.sightingSvc.GetTigers(mockCtx)
				require.Equal(t, errors.New("db error"), resErr)
				require.Nil(t, resData)
			},
		},
		{
			testcaseName: "successfully get all the data from database",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				callbackFunc := func(ctx context.Context, key string, anyTigers *[]*entity.Tiger, ttl time.Duration, callback func() (interface{}, error)) {
					serviceTestSuite.sightingRepo.EXPECT().GetTigers(mockCtx).Return(tigerData, nil)
					_, _ = callback()
				}

				serviceTestSuite.redisRepo.EXPECT().Fetch(mockCtx, service.GetTigersKey, &emptyTiger, mockTTL, gomock.Any()).Do(callbackFunc).Return(nil)

				resData, resErr := serviceTestSuite.sightingSvc.GetTigers(mockCtx)
				require.NoError(t, resErr)
				require.NotNil(t, resData)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestCreateTiger(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCtx := context.Background()
	tigerData := &entity.Tiger{Name: "tiger-1",
		DateOfBirth: time.Now(), LastSeenTimestamp: time.Now(), LastSeenLatitude: -6.18, LastSeenLongitude: 106.0}

	testCases := []ServiceTestCase{
		{
			testcaseName: "Error invalid name",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				tigerData2 := *tigerData
				tigerData2.Name = ""
				resErr := serviceTestSuite.sightingSvc.CreateTiger(mockCtx, &tigerData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error invalid date of birth",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				tigerData2 := *tigerData
				tigerData2.DateOfBirth = time.Unix(0, 0)
				resErr := serviceTestSuite.sightingSvc.CreateTiger(mockCtx, &tigerData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error invalid last seen timestamp",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				tigerData2 := *tigerData
				tigerData2.LastSeenTimestamp = time.Unix(0, 0)
				resErr := serviceTestSuite.sightingSvc.CreateTiger(mockCtx, &tigerData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error invalid latitude",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				tigerData2 := *tigerData
				tigerData2.LastSeenLatitude = 200.0
				resErr := serviceTestSuite.sightingSvc.CreateTiger(mockCtx, &tigerData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error invalid longitude",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				tigerData2 := *tigerData
				tigerData2.LastSeenLongitude = 200.0
				resErr := serviceTestSuite.sightingSvc.CreateTiger(mockCtx, &tigerData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error when insert to database",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceTestSuite.sightingRepo.EXPECT().CreateTiger(mockCtx, tigerData).Return(errors.New("db error"))

				resErr := serviceTestSuite.sightingSvc.CreateTiger(mockCtx, tigerData)
				require.Equal(t, errors.New("db error"), resErr)
			},
		},
		{
			testcaseName: "successfully insert to database",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceTestSuite.sightingRepo.EXPECT().CreateTiger(mockCtx, tigerData).Return(nil)
				serviceTestSuite.redisRepo.EXPECT().Del(mockCtx, service.GetTigersKey)

				resErr := serviceTestSuite.sightingSvc.CreateTiger(mockCtx, tigerData)
				require.NoError(t, resErr)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestGetSightingsByTigerID(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCtx := context.Background()
	var emptySighting []*entity.Sighting
	sightingsData := []*entity.Sighting{{ID: 1}}
	mockTTL := 1 * time.Minute
	tigerID := int32(1)

	testCases := []ServiceTestCase{
		{
			testcaseName: "successfully get the data from redis",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceTestSuite.redisRepo.EXPECT().Fetch(mockCtx, fmt.Sprintf(service.GetSightingsByTigerIDKey, tigerID), &emptySighting, mockTTL, gomock.Any()).
					SetArg(2, sightingsData).Return(nil)

				resData, resErr := serviceTestSuite.sightingSvc.GetSightingsByTigerID(mockCtx, tigerID)

				require.NoError(t, resErr)
				require.NotNil(t, resData)
				require.Equal(t, sightingsData, resData)
			},
		},
		{
			testcaseName: "Error when retrieve from database",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				callbackFunc := func(ctx context.Context, key string, anySightings *[]*entity.Sighting, ttl time.Duration, callback func() (interface{}, error)) {
					serviceTestSuite.sightingRepo.EXPECT().GetSightingsByTigerID(mockCtx, tigerID).Return(nil, errors.New("db error"))
					_, _ = callback()
				}

				serviceTestSuite.redisRepo.EXPECT().Fetch(mockCtx, fmt.Sprintf(service.GetSightingsByTigerIDKey, tigerID), &emptySighting, mockTTL, gomock.Any()).Do(callbackFunc).Return(errors.New("db error"))

				resData, resErr := serviceTestSuite.sightingSvc.GetSightingsByTigerID(mockCtx, tigerID)
				require.Equal(t, errors.New("db error"), resErr)
				require.Nil(t, resData)
			},
		},
		{
			testcaseName: "successfully get all the data from database",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				callbackFunc := func(ctx context.Context, key string, anySightings *[]*entity.Sighting, ttl time.Duration, callback func() (interface{}, error)) {
					serviceTestSuite.sightingRepo.EXPECT().GetSightingsByTigerID(mockCtx, tigerID).Return(sightingsData, nil)
					_, _ = callback()
				}

				serviceTestSuite.redisRepo.EXPECT().Fetch(mockCtx, fmt.Sprintf(service.GetSightingsByTigerIDKey, tigerID), &emptySighting, mockTTL, gomock.Any()).Do(callbackFunc).Return(nil)

				resData, resErr := serviceTestSuite.sightingSvc.GetSightingsByTigerID(mockCtx, tigerID)
				require.NoError(t, resErr)
				require.NotNil(t, resData)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestCreateSighting(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCtx := context.Background()
	tigerID := int32(1)
	tigerData := &entity.Tiger{Name: "tiger-1",
		DateOfBirth: time.Now(), LastSeenTimestamp: time.Now(), LastSeenLatitude: -6.18, LastSeenLongitude: 106.0}
	imageData := "data:image/png;base64,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lwHEKdtiJTl7FFFwhhIDBSWQ+wE2GwHZtHEjtaamqNVqqQQZGxuj0dFRGsztW1RB0x+TgQAIuCCQaws9eQERp/A4LLtYbxgDBLIl0K9DkuLWkaSAuGxRZZulcAwEQEA8gagH4BmimQGiAQtKgQWkvsbmWgJaxA6PggAIZEJA/P7WZ+sOLCDVox1VoaubiydAAARAwIpASh0WkQKSEkCrTCGi+jWZ7cxSngcBKZGAHfII1D5AB1pWYgREfAknL7eysShQrmfDC44kTqBGwkvdH+0EpAaIhaE3VlgHcyWedv7MB1t/bDEyCDgmIKlDYycgNcFIAlDTBTwmnQBEUXqEDOxDEHWQjA/guoFqXg8mIFJLsJrc8BgIgIBXAmWIxx9//EG77rorHXPMMXT//fcrovzD0Oeffz699dZbtPnmm9O5555Lt9xyCw0MDNCmTZvommuuoccee4z+/fdfOvzww+nhhx+mFStW0JTmB6nfe+89OuOMM2jvvfem5557rh292267ja699lpaunRp++9OPvlkevzxx9X/P/HEE3TDDTcou3baaSd66KGHaI899lDXvAtIbIX0muMYPCyBMvaUsEwxW1QCZ599Nr377rt0xBFHtAXkpJNOouXLl9N9991HvH8edthhdOWVV9J5551H99xzDz344IP0+uuv01ZbbUXnnHOO+m0bTz75ZNuPXh2eV155he68807aeeedad26dfME5Prrr6fffvuN7r333i4WX3zxBR188MH06quv0r777qvmv/322+nbb7+lzTbbzI+AoEUVNScxOQiAQAIEXn75Zbr55pvp6KOPpp9//lkJyF9//UUjIyP0ww8/0DbbbKO84BP/008/rSqSAw88kC666CI6/fTT1bXvvvuOdtttN+JKhsWEKxEWHv6ceeaZtGHDBrrxxhvp448/pn322UdVFV999dU8Abn44ouVGHGVs/DD4vLLL7+oKmf2w3Y99dRTdMghh7gTELSoDDIWJ2gDSLhljgASJtds4BM/b+hcGTz77LP0008/KQH5/PPPlUisXbu27TpXKCeeeKLayFeuXElvvPGGakPxZ2ZmRrW5+Lmtt96adt99d3rmmWdoenqauLrhv2dRmd2fb7rpJvrmm29UxTL7uwJPO+00ajabqgrhNth+++2nqpVtt92WTjjhBGUPV0CzH26bsT0sZNYtLLSock1x+AUCIOCLAFcQe+21F11xxRXEm/qsgHzwwQd03HHH0Zo1a9pTf/LJJ+q0/+eff9IWW2xBn376Ke2yyy7t68PDw6o6YUF67bXX6Oqrr1bvR7jdxO2vzg/PxePdfffdSiz4WRYw7hpdfvnlqi116aWXqjn4z1FHHaXs4Spl9nPssccqe6666qp6AoIWla+0wrggAAK5E3j++efpjjvuUJv+kiVL5gkIv3Pg6uKff/5pY3jzzTdVy4pP+1w1vPjii3TAAQeo69yi4grk66+/ph133FH9Hb+U5/cin332WRdKFpDx8fF2C6tXAcD7O8/z448/0mWXXabsue6669pjHXTQQcqeCy+80FxA5rWoZojG/pPhr0tPumOQtPG57xnwDwTaBE455RR655131MbPH35/wULAL6lfeukl9Q6EhWSHHXZQ17mdxO9LuHXFmzd/k+qCCy5Q11gMuMXE7Sce74EHHmi3sM466yz1ba6FFUingPC1Dz/8kLbffnvVDuNvWvEfrjC4rcXvU77//nv1DoY/rAOrVq2iF154Qc2rbWGhRYXMBwEQAAF/BDpbWDwLv5Pgr+s+8sgjqpXF7xz4nlNPPVV9A4urFxYTfvHNlQC/+2Dh4Bfv+++/P3Eb7O+//1Zi89FHH9F2223XNn5hBcIX+CU+//tJPDZ/VZjFid+T8Mt7fu/C3wpjAeOx+Su//LXeL7/8UlVPPQUELSp/yYKRQQAEQGBhVTD7DoT/nqsJ/soui8TQ0BBdcsklxN+G4g9XCfzfvNnzew5+H8Ev37fccks69NBD6fjjj1fvVfjD377i9he3yo488kh6//33VQXB4sTvOrjlxW2uyclJ9Y7j7bffVq0v/truXXfdRY1GQ32NmL+5deutt9Kvv/5Ke+65Jz366KPtdllbQLgs0f0gCsKeKAF0txINHMwGARkE+n3LdvXq1TQwPj4+wzewivEPsHT+RKKZ+fxPkvAuhQ8IgAAIgEDOBPh9Df8wIv/MClcs/wXIZT3M35g24AAAAABJRU5ErkJggg=="
	sightingData := &entity.Sighting{TigerID: tigerID, SeenAt: time.Now(), Latitude: -6.18, Longitude: 106.0,
		ImageData: imageData}

	testCases := []ServiceTestCase{
		{
			testcaseName: "Error invalid tigerID",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				sightingData2 := *sightingData
				sightingData2.TigerID = 0
				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, &sightingData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error invalid seen at timestamp",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				sightingData2 := *sightingData
				sightingData2.SeenAt = time.Unix(0, 0)
				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, &sightingData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error invalid latitude",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				sightingData2 := *sightingData
				sightingData2.Latitude = 200.0
				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, &sightingData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error invalid longitude",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				sightingData2 := *sightingData
				sightingData2.Longitude = 200.0
				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, &sightingData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error empty image data",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				sightingData2 := *sightingData
				sightingData2.ImageData = ""
				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, &sightingData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error no tiger in database",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				serviceTestSuite.sightingRepo.EXPECT().GetTigerByID(mockCtx, tigerID).Return(nil, errors.New("db error"))

				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, sightingData)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error distance exceed 5 km",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				sightingData2 := *sightingData
				sightingData2.Latitude = -8.10

				serviceTestSuite.sightingRepo.EXPECT().GetTigerByID(mockCtx, tigerID).Return(tigerData, nil)

				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, &sightingData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error invalid base64 image",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)
				sightingData2 := *sightingData
				sightingData2.ImageData = "this-is-invalid-base64-image"

				serviceTestSuite.sightingRepo.EXPECT().GetTigerByID(mockCtx, tigerID).Return(tigerData, nil)

				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, &sightingData2)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error when insert to database",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)

				serviceTestSuite.sightingRepo.EXPECT().GetTigerByID(mockCtx, tigerID).Return(tigerData, nil)
				serviceTestSuite.sightingRepo.EXPECT().CreateSighting(mockCtx, sightingData).Return(errors.New("db error"))

				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, sightingData)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "Error when update tiger data",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)

				serviceTestSuite.sightingRepo.EXPECT().GetTigerByID(mockCtx, tigerID).Return(tigerData, nil)
				serviceTestSuite.sightingRepo.EXPECT().CreateSighting(mockCtx, sightingData).Return(nil)
				serviceTestSuite.sightingRepo.EXPECT().UpdateTiger(mockCtx, tigerData).Return(errors.New("error db"))

				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, sightingData)
				require.Error(t, resErr)
			},
		},
		{
			testcaseName: "successfully insert to database",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)

				serviceTestSuite.sightingRepo.EXPECT().GetTigerByID(mockCtx, tigerID).Return(tigerData, nil)
				serviceTestSuite.sightingRepo.EXPECT().CreateSighting(mockCtx, sightingData).Return(nil)
				serviceTestSuite.sightingRepo.EXPECT().UpdateTiger(mockCtx, tigerData).Return(nil)
				serviceTestSuite.redisRepo.EXPECT().Del(mockCtx, service.GetTigersKey).Return(nil)
				serviceTestSuite.redisRepo.EXPECT().Del(mockCtx, fmt.Sprintf(service.GetSightingsByTigerIDKey, tigerID)).Return(nil)

				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, sightingData)
				require.NoError(t, resErr)
			},
		},
		{
			testcaseName: "successfully insert to database using jpeg image",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				serviceTestSuite := TigerSightingServiceTestSuite(mockCtrl)

				sightingData2 := *sightingData
				sightingData2.ImageData = "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/4QAiRXhpZgAATU0AKgAAAAgAAQESAAMAAAABAAEAAAAAAAD/2wBDAAIBAQIBAQICAgICAgICAwUDAwMDAwYEBAMFBwYHBwcGBwcICQsJCAgKCAcHCg0KCgsMDAwMBwkODw0MDgsMDAz/2wBDAQICAgMDAwYDAwYMCAcIDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAz/wAARCAAyAEsDASIAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD9EKKKKACipLKzm1K/gtbeMy3F1KkMSAgF3ZgqrzxySBzxzXXXnwXm0+6kt7jxZ8P4LiFzHLFJrRV43U4ZWHlcEEEEdiKAONorrf8AhUf/AFOPw7/8HZ/+NUf8Kj/6nH4d/wDg7P8A8aoA5Kitrxh4Eu/BcdjLNdaVqFpqSSNbXWnXX2iCUxtsdQ2AcqxAPGMnGcggYtABRRRQAUUUUAa3w+/5KH4d/wCwtZ/+j0r1LwH4h8DaP4q8dR+KoLF719dvHjkurM3CmASt8qfK2GDbyQACcr1xx5b8Pv8Akofh3/sLWf8A6PSvX/hB8ILDxr8S/FmvakwuIdN8RXtvDaMvyPKspfe/qBvXC+oOc8CgDx3RvDVx448Wf2foNnNIbqVzbxO2TDFu4MjcgBVIy3r0ySAev+Lv7PmofDDTLfUI5v7SsNirdSqm37NL3yP+eZPRj06HsT9GeEvhzongS4vpdJsobOTUpfNmK/oq/wB1ByQowoLHAFbVzDHeQtFIqSRyAq6MAyuDwQR3BoA+RfGP/JHPh/8A9xb/ANKxXI11ni0/8WX+Hv01X/0rWuToAKKKKACiiigCbTdQm0fVLW8t2VbizmS4iLLuAdGDKSO/IHFddqXxP0PWb+a7vPAehXN3dSNNNKbmUeY7EszY7ZYk4964uigDrf8AhPPDP/RPNB/8CpaP+E88M/8ARPNB/wDAqWuSooA3vG3jr/hL7XTbWDTbPSdP0lJVtra3LNtMrh5CWY5OWAPbHPrWDRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQB//Z"

				serviceTestSuite.sightingRepo.EXPECT().GetTigerByID(mockCtx, tigerID).Return(tigerData, nil)
				serviceTestSuite.sightingRepo.EXPECT().CreateSighting(mockCtx, sightingData2).Return(nil)
				serviceTestSuite.sightingRepo.EXPECT().UpdateTiger(mockCtx, tigerData).Return(nil)
				serviceTestSuite.redisRepo.EXPECT().Del(mockCtx, service.GetTigersKey).Return(nil)
				serviceTestSuite.redisRepo.EXPECT().Del(mockCtx, fmt.Sprintf(service.GetSightingsByTigerIDKey, tigerID)).Return(nil)

				resErr := serviceTestSuite.sightingSvc.CreateSighting(mockCtx, &sightingData2)
				require.NoError(t, resErr)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}
