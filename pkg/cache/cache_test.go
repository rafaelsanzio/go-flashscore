package cache

import (
	"context"
	"testing"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

func mockCacheSetFunc(ctx context.Context, key string, value []byte) error {
	return nil
}

func mockCacheSetThrowFunc(ctx context.Context, key string, value []byte) error {
	return errs.ErrRepoMockAction
}

func mockCacheGetFunc(ctx context.Context, key string) (interface{}, error) {
	return []byte("any_data"), nil
}

func TestSetCache(t *testing.T) {
	testCases := []struct {
		Name         string
		CacheKey     string
		Data         []byte
		CacheSetFunc func(ctx context.Context, key string, value []byte) error
		CacheGetFunc func(ctx context.Context, key string) (interface{}, error)
	}{
		{
			Name:         "Success setting cache",
			CacheKey:     "any_key",
			Data:         []byte("any_data"),
			CacheSetFunc: mockCacheSetFunc,
			CacheGetFunc: mockCacheGetFunc,
		}, {
			Name:         "Throwing setting cache",
			CacheKey:     "any_key",
			Data:         []byte("any_data"),
			CacheSetFunc: mockCacheSetThrowFunc,
			CacheGetFunc: nil,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		SetStore(MockCacheStore{
			SetFunc: tc.CacheSetFunc,
			GetFunc: tc.CacheGetFunc,
		})
		defer SetStore(nil)

		SetCache(context.Background(), tc.CacheKey, tc.Data)
	}
}
