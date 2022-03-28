package redis_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	redismock "github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/suite"

	"github.com/ibrahimker/tigerhall-kittens/driver/redis"
)

const authRedisBaseKey = "dummy-name" + ":" + "dummy-version"

type RedisTestStruct struct {
	Val string
}

type RedisTestSuite struct {
	suite.Suite
	client redis.Redis
	mock   redismock.ClientMock
}

func TestRedis(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

func (s *RedisTestSuite) SetupTest() {
	rds, mock := redismock.NewClientMock()

	s.client = redis.NewRedisClient(rds)
	s.mock = mock
}

func (s *RedisTestSuite) TestDel() {
	ctx := context.Background()

	s.Run("Del returns error", func() {
		key := "del-cache-key"

		s.mock.ExpectDel(cacheKey(key)).SetErr(errors.New("redis del error"))

		err := s.client.Del(ctx, cacheKey(key))

		s.Error(err)
	})

	s.Run("Success delete single key", func() {
		key := "del-cache-key"

		s.mock.ExpectDel(cacheKey(key)).SetVal(1)

		err := s.client.Del(ctx, cacheKey(key))

		s.NoError(err)
	})

	s.Run("Success delete multiple key", func() {
		keys := []string{cacheKey("key-1"), cacheKey("key-2"), cacheKey("key-3")}
		cacheKeys := []string{authRedisBaseKey + ":key-1", authRedisBaseKey + ":key-2", authRedisBaseKey + ":key-3"}

		s.mock.ExpectDel(cacheKeys...).SetVal(1)

		err := s.client.Del(ctx, keys...)

		s.NoError(err)
	})
}

func (s *RedisTestSuite) TestFetch() {
	key := "redis-fetch-key"
	ctx := context.Background()
	exp := 5 * time.Minute

	s.Run("Get redis returns error", func() {
		value := &RedisTestStruct{}

		s.mock.ExpectGet(cacheKey(key)).SetErr(errors.New("Error get redis"))

		err := s.client.Fetch(ctx, cacheKey(key), value, exp, func() (interface{}, error) {
			return nil, nil
		})

		s.Error(err)
		s.Equal(err.Error(), "Error get redis")
	})

	s.Run("Get returns value not nil", func() {
		value := &RedisTestStruct{}
		cacheVal := &RedisTestStruct{Val: "amazing-value"}
		bytes, _ := json.Marshal(cacheVal)

		s.mock.ExpectGet(cacheKey(key)).SetVal(string(bytes))

		err := s.client.Fetch(ctx, cacheKey(key), value, exp, func() (interface{}, error) {
			return nil, nil
		})

		s.NoError(err)
		s.Equal(value.Val, cacheVal.Val)
	})

	s.Run("Get returns nil, callback returns error", func() {
		value := &RedisTestStruct{}

		s.mock.ExpectGet(cacheKey(key)).RedisNil()

		err := s.client.Fetch(ctx, cacheKey(key), value, exp, func() (interface{}, error) {
			return nil, errors.New("callback error")
		})

		s.Error(err)
		s.Equal(err.Error(), "callback error")
	})

	s.Run("Get returns nil, marshall returns error", func() {
		value := &RedisTestStruct{}

		s.mock.ExpectGet(cacheKey(key)).RedisNil()

		err := s.client.Fetch(ctx, cacheKey(key), value, exp, func() (interface{}, error) {
			return make(chan int), nil
		})

		s.Error(err)
	})

	s.Run("Get returns nil, Set returns error", func() {
		value := &RedisTestStruct{}
		callValue := &RedisTestStruct{Val: "incredible-value"}
		bytes, _ := json.Marshal(callValue)

		s.mock.ExpectGet(cacheKey(key)).RedisNil()
		s.mock.ExpectSet(cacheKey(key), bytes, exp).SetErr(errors.New("Set redis error"))

		err := s.client.Fetch(ctx, cacheKey(key), value, exp, func() (interface{}, error) {
			return callValue, nil
		})

		s.Error(err)
	})

	s.Run("Get returns nil, Success set new cache", func() {
		value := &RedisTestStruct{}
		callValue := &RedisTestStruct{Val: "spectacular-value"}
		bytes, _ := json.Marshal(callValue)

		s.mock.ExpectSet(cacheKey(key), bytes, exp).SetVal("success")
		s.mock.ExpectGet(cacheKey(key)).SetVal("any-val")

		err := s.client.Fetch(ctx, cacheKey(key), value, exp, func() (interface{}, error) {
			return callValue, nil
		})

		s.Error(err)
	})
}

func cacheKey(key string) string {
	return fmt.Sprintf("%s:%s", authRedisBaseKey, key)
}
