package service_test

import (
	"context"
	"errors"
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
				callbackFunc := func(ctx context.Context, key string, anySchoolProfile *[]*entity.Tiger, ttl time.Duration, callback func() (interface{}, error)) {
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
				callbackFunc := func(ctx context.Context, key string, anySchoolProfile *[]*entity.Tiger, ttl time.Duration, callback func() (interface{}, error)) {
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []ServiceTestCase{}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestCreateSighting(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []ServiceTestCase{}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}
