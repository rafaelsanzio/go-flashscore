package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/config"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

type Store struct {
	client *redis.Client
}

func NewStore() (*Store, errs.AppError) {
	client, err := Client()
	if err != nil {
		return nil, err
	}

	return &Store{
		client: client,
	}, nil
}

func Client() (*redis.Client, errs.AppError) {
	redisPort, err_ := config.Value(key.RedisPort)
	if err_ != nil {
		errApp := err_.Annotatef(applog.Log, "unable to get redis port config: %v", err_)
		return nil, errApp
	}

	redisServer := fmt.Sprintf("redis:%s", redisPort)

	var cache = redis.NewClient(&redis.Options{
		Addr: redisServer,
	})

	return cache, nil
}
