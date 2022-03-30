package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"

	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/entity"
	"github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/repository/postgres"
)

type SightingTestSuite struct {
	repo *postgres.TigerSightingRepo
	pgx  pgxmock.PgxPoolIface
}

type RepositoryTestCases struct {
	testcaseName     string
	testcaseFunction func(t *testing.T)
}

func SightingRepositoryTestSuite() *SightingTestSuite {
	mock, _ := pgxmock.NewPool(pgxmock.MonitorPingsOption(true))

	return &SightingTestSuite{
		repo: postgres.NewTigerSightingRepo(mock),
		pgx:  mock,
	}
}

func TestRepositorySuite_TestSuite(t *testing.T) {
	t.Parallel()
	t.Run("successfully create an instance of repository", func(t *testing.T) {
		t.Parallel()
		repositorySuite := SightingRepositoryTestSuite()
		require.NotNil(t, repositorySuite.repo)
		require.NotNil(t, repositorySuite.pgx)
	})
}

func TestGetTigers(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	queryString := `SELECT id,name,date_of_birth,last_seen_timestamp,last_seen_latitude,last_seen_longitude,created_at,updated_at
FROM sighting.tiger WHERE deleted_at IS NULL ORDER BY last_seen_timestamp desc`
	queryStringRow := []string{"id", "name", "date_of_birth", "last_seen_timestamp", "last_seen_latitude", "last_seen_longitude", "created_at", "updated_at"}
	expQueryStringRes := []interface{}{int32(1), "tiger-1", time.Now(), time.Now(), -6.19, 108.0, sql.NullTime{Time: time.Now()}, sql.NullTime{Time: time.Now()}}

	testCases := []RepositoryTestCases{
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnError(pgx.ErrNoRows)

				resData, err := repositorySuite.repo.GetTigers(context.Background())
				require.Error(t, err)
				require.Equal(t, 0, len(resData))
			},
		},
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnError(pgx.ErrNoRows)

				resData, err := repositorySuite.repo.GetTigers(context.Background())
				require.Error(t, err)
				require.Equal(t, 0, len(resData))
			},
		},
		{
			testcaseName: "Error when scanning rows",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnRows(pgxmock.
						NewRows([]string{"id"}).
						AddRow("test-id"),
					)

				resData, err := repositorySuite.repo.GetTigers(context.Background())
				require.NoError(t, err)
				require.Equal(t, 0, len(resData))
			},
		},
		{
			testcaseName: "Error when check rows",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnRows(pgxmock.
						NewRows(queryStringRow).
						AddRow(expQueryStringRes...).RowError(1, pgx.ErrNoRows),
					)

				resData, err := repositorySuite.repo.GetTigers(context.Background())
				require.Error(t, err)
				require.Equal(t, 0, len(resData))
			},
		},
		{
			testcaseName: "sucessfullly retrieve tigers data",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnRows(pgxmock.
						NewRows(queryStringRow).
						AddRow(expQueryStringRes...),
					)

				resData, err := repositorySuite.repo.GetTigers(context.Background())
				require.NoError(t, err)
				require.Equal(t, 1, len(resData))
				require.Equal(t, expQueryStringRes, expQueryStringRes)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestGetTigerByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	queryString := `SELECT id,name,date_of_birth,last_seen_timestamp,last_seen_latitude,last_seen_longitude,created_at,updated_at
FROM sighting.tiger WHERE id = \$1 and deleted_at IS NULL`
	queryStringRow := []string{"id", "name", "date_of_birth", "last_seen_timestamp", "last_seen_latitude", "last_seen_longitude", "created_at", "updated_at"}
	expQueryStringRes := []interface{}{int32(1), "tiger-1", time.Now(), time.Now(), -6.19, 108.0, sql.NullTime{Time: time.Now()}, sql.NullTime{Time: time.Now()}}
	tigerID := int32(1)

	testCases := []RepositoryTestCases{
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnError(pgx.ErrNoRows)

				resData, err := repositorySuite.repo.GetTigerByID(context.Background(), tigerID)
				require.Error(t, err)
				require.Nil(t, resData)
			},
		},
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnError(pgx.ErrNoRows)

				resData, err := repositorySuite.repo.GetTigerByID(context.Background(), tigerID)
				require.Error(t, err)
				require.Nil(t, resData)
			},
		},
		{
			testcaseName: "Error when scanning rows",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnRows(pgxmock.
						NewRows([]string{"id"}).
						AddRow("test-id"),
					)

				resData, err := repositorySuite.repo.GetTigerByID(context.Background(), tigerID)
				require.NoError(t, err)
				require.Equal(t, &entity.Tiger{}, resData)
			},
		},
		{
			testcaseName: "Error when check rows",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnRows(pgxmock.
						NewRows(queryStringRow).
						AddRow(expQueryStringRes...).RowError(1, pgx.ErrNoRows),
					)

				resData, err := repositorySuite.repo.GetTigerByID(context.Background(), tigerID)
				require.Error(t, err)
				require.Nil(t, resData)
			},
		},
		{
			testcaseName: "sucessfullly retrieve tigers data",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnRows(pgxmock.
						NewRows(queryStringRow).
						AddRow(expQueryStringRes...),
					)

				resData, err := repositorySuite.repo.GetTigerByID(context.Background(), tigerID)
				require.NoError(t, err)
				require.Equal(t, int32(1), resData.ID)
				require.Equal(t, expQueryStringRes, expQueryStringRes)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestCreateTiger(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	queryString := `INSERT INTO sighting.tiger \(name,date_of_birth,last_seen_timestamp,last_seen_latitude,last_seen_longitude,created_at,updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`
	tiger := &entity.Tiger{
		Name:              "tiger 1",
		DateOfBirth:       time.Now(),
		LastSeenTimestamp: time.Now(),
		LastSeenLatitude:  -6.18,
		LastSeenLongitude: 108.00,
	}

	testCases := []RepositoryTestCases{
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectExec(queryString).
					WillReturnError(pgx.ErrNoRows)

				err := repositorySuite.repo.CreateTiger(context.Background(), tiger)
				require.Error(t, err)
			},
		},
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectExec(queryString).
					WillReturnError(pgx.ErrNoRows)

				err := repositorySuite.repo.CreateTiger(context.Background(), tiger)
				require.Error(t, err)
			},
		},
		{
			testcaseName: "sucessfullly create tiger data",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectExec(queryString).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				err := repositorySuite.repo.CreateTiger(context.Background(), tiger)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestUpdateTiger(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	queryString := `UPDATE sighting.tiger SET last_seen_timestamp = \$2, last_seen_latitude = \$3, last_seen_longitude = \$4, updated_at = \$5 WHERE id = \$1`
	tiger := &entity.Tiger{
		Name:              "tiger 1",
		DateOfBirth:       time.Now(),
		LastSeenTimestamp: time.Now(),
		LastSeenLatitude:  -6.18,
		LastSeenLongitude: 108.00,
	}

	testCases := []RepositoryTestCases{
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectExec(queryString).
					WillReturnError(pgx.ErrNoRows)

				err := repositorySuite.repo.UpdateTiger(context.Background(), tiger)
				require.Error(t, err)
			},
		},
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectExec(queryString).
					WillReturnError(pgx.ErrNoRows)

				err := repositorySuite.repo.UpdateTiger(context.Background(), tiger)
				require.Error(t, err)
			},
		},
		{
			testcaseName: "sucessfullly update tiger data",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectExec(queryString).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				err := repositorySuite.repo.UpdateTiger(context.Background(), tiger)
				require.NoError(t, err)
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
	queryString := `SELECT id,tiger_id,seen_at,latitude,longitude,image_data,created_at,updated_at
FROM sighting.sighting WHERE tiger_id = \$1 and deleted_at IS NULL ORDER BY seen_at desc`
	queryStringRow := []string{"id", "tiger_id", "seen_at", "latitude", "longitude", "image_data", "created_at", "updated_at"}
	expQueryStringRes := []interface{}{int32(1), int32(1), time.Now(), -6.19, 108.0, "https://test.com/dummy.jpeg", sql.NullTime{Time: time.Now()}, sql.NullTime{Time: time.Now()}}
	tigerID := int32(1)

	testCases := []RepositoryTestCases{
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnError(pgx.ErrNoRows)

				resData, err := repositorySuite.repo.GetSightingsByTigerID(context.Background(), tigerID)
				require.Error(t, err)
				require.Equal(t, 0, len(resData))
			},
		},
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnError(pgx.ErrNoRows)

				resData, err := repositorySuite.repo.GetSightingsByTigerID(context.Background(), tigerID)
				require.Error(t, err)
				require.Equal(t, 0, len(resData))
			},
		},
		{
			testcaseName: "Error when scanning rows",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnRows(pgxmock.
						NewRows([]string{"id"}).
						AddRow("test-id"),
					)

				resData, err := repositorySuite.repo.GetSightingsByTigerID(context.Background(), tigerID)
				require.NoError(t, err)
				require.Equal(t, 0, len(resData))
			},
		},
		{
			testcaseName: "Error when check rows",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnRows(pgxmock.
						NewRows(queryStringRow).
						AddRow(expQueryStringRes...).RowError(1, pgx.ErrNoRows),
					)

				resData, err := repositorySuite.repo.GetSightingsByTigerID(context.Background(), tigerID)
				require.Error(t, err)
				require.Equal(t, 0, len(resData))
			},
		},
		{
			testcaseName: "sucessfullly retrieve sighting data",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectQuery(queryString).
					WillReturnRows(pgxmock.
						NewRows(queryStringRow).
						AddRow(expQueryStringRes...),
					)

				resData, err := repositorySuite.repo.GetSightingsByTigerID(context.Background(), tigerID)
				require.NoError(t, err)
				require.Equal(t, 1, len(resData))
				require.Equal(t, expQueryStringRes, expQueryStringRes)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func TestCreateSighting(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	queryString := `INSERT INTO sighting.sighting \(tiger_id,seen_at,latitude,longitude,image_data,created_at,updated_at\) 
VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`
	sighting := &entity.Sighting{
		TigerID:   1,
		SeenAt:    time.Now(),
		Latitude:  -6.18,
		Longitude: 107.00,
		ImageData: "data:image/png;base64,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aUIec8wxU46imoBad9ZZZ+nfqYlgjz32ICV5FennzJlDS5cu1b9T9tVU8MQTT+gJ4Pjjj6etttpK/+7aa6+lH3/8kZ566indWJptfvLJJ3TggQfSr7/+Snfffbdep6aGxrPtttvSs88+S4cffrgZ0Xn1OEpbVsEz5oPRdAwM1Vi9++6706effko77rgj9ff30wEHHEBffPGFHqtPOukkuuuuu2jGjBm0xRZb6LF+n3320aGNjY3RZpttptdus802+tZ+4YUX6O+//9ZNQ/1ckb/5UZfrfvvtp6WDahidbL744oukxvlFixbRVVddNWFCyYVTTz1VN4fK0R16PEYJRfQhnGydw0sb+MjICB177LG0ePFiuvTSS3XCFIH23XdfTcKff/6ZTj75ZFKjtxqZlaZWN6669RvP7Nmz6Z133tF73njjDX1bKzKrMfuII46YVARr166lc889l9avX0+KyOrpZlO9O1DxqRG/8Zxwwgn6NleNoivRo+rxiLUOV0DABoGPP/5Y63D1oqt5jG+18fLLL9M111xD3377LW255Zb06quvai2vHkVYdaMvX76cdtllF/0zNR0o3a60f/OjJMMpp5yiX6Tdd999+qWcerrZVATfeeed6YILLqCBDf+/XTUd1YSWLFnSnujQ4zZlgLVREEh0oX/44Yf6ZZjSuoo4jeevv/4i9Ts1Ljee5557jm677TZNZrX27LPP1hOAej777DOtsX///XdN+IcffnhidFc390UXXaTXqclg0X//S0sWL6bLLrtsErTdbN5xxx26wSjdPzw8rJuJ0uivvPKK9jvlRufV41FKQoCTRFUp4ORTQsgYCvVme3BwUH9Mpd5+Nz/qd9tttx3deuutepT/4Ycf9Es09ZHbnXfeqd+4qxd1SqerN+PqZlXaXBH8u+++0y/R3n//fVIjuiLwRx99RDvssIPW1OqjMWWj9elmUzUX9b7gtdde09Lg+uuv12/8ly1bpl/0TSK6jx7POJ8S6dE5JgAdLV8vvfSSvs2VNm5+lDZXJP7ggw/oyiuv1De40t9nnHEG3X777Xq9evl20003acKrCVnpZfUx3cyZM/UUoF7cqb3qUXvefvttrcW33nprfeP39Gz8IwtK5w8NDbW3ufQh2nz25trO888/TzfccIN++7733nvTzTffrOWBemk4QXTVoVQQ6of2T52qL5ez5hKnfbVhRxUCG3PfuLxnzZpFPUNDQ2Pq1X1r56oyx/p71dBUvCEfDpsh44MtIMCAgJIH6rP2/wN/XBtlM2lCsQAAAABJRU5ErkJggg==",
	}

	testCases := []RepositoryTestCases{
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectExec(queryString).
					WillReturnError(pgx.ErrNoRows)

				err := repositorySuite.repo.CreateSighting(context.Background(), sighting)
				require.Error(t, err)
			},
		},
		{
			testcaseName: "database returns no rows when scanning",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()

				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectExec(queryString).
					WillReturnError(pgx.ErrNoRows)

				err := repositorySuite.repo.CreateSighting(context.Background(), sighting)
				require.Error(t, err)
			},
		},
		{
			testcaseName: "sucessfullly create tiger data",
			testcaseFunction: func(t *testing.T) {
				t.Parallel()
				repositorySuite := SightingRepositoryTestSuite()
				repositorySuite.pgx.
					ExpectExec(queryString).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				err := repositorySuite.repo.CreateSighting(context.Background(), sighting)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}
