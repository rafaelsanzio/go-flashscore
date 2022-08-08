package cache

import (
	"context"
)

type MockCacheStore struct {
	Store
	SetFunc func(ctx context.Context, key string, value []byte) error
	GetFunc func(ctx context.Context, key string) (interface{}, error)
}

func (m MockCacheStore) Get(ctx context.Context, key string) (interface{}, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	return m.Store.Get(ctx, key)
}

func (m MockCacheStore) Set(ctx context.Context, key string, value []byte) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value)
	}
	return m.Store.Set(ctx, key, value)
}
