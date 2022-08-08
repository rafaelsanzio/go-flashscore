package cache

import (
	"context"
	"log"

	"github.com/rafaelsanzio/go-flashscore/pkg/cache/redis"
)

type Store interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) (interface{}, error)
}

var store Store

func init() {
	s, err := redis.NewStore()
	if err != nil {
		log.Fatalf("error creating redis store: %v", err)
	}

	store = s
}

func GetStore() Store {
	return store
}

func SetStore(s Store) {
	store = s
}
