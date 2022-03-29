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

func TestGetSightingsByTigerID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []RepositoryTestCases{}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}

func (s *SightingTestSuite) TestCreateSighting(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []RepositoryTestCases{}

	for _, tc := range testCases {
		t.Run(tc.testcaseName, tc.testcaseFunction)
	}
}
