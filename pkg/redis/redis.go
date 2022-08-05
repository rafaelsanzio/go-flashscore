package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/config"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
)

func RedisConnect() *redis.Client {
	redisPort, err_ := config.Value(key.RedisPort)
	if err_ != nil {
		_ = err_.Annotatef(applog.Log, "unable to get redis port config: %v", err_)
	}

	redisServer := fmt.Sprintf("redis:%s", redisPort)

	var cache = redis.NewClient(&redis.Options{
		Addr: redisServer,
	})

	return cache
}

func Set(ctx context.Context, key string, value []byte) error {
	conn := RedisConnect()
	defer conn.Close()

	err := conn.Set(ctx, key, value, 10*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func Get(ctx context.Context, key string) (interface{}, error) {
	conn := RedisConnect()
	defer conn.Close()

	var thing interface{}
	thingCached, err := conn.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	if thingCached != nil {
		err = json.Unmarshal(thingCached, &thing)
		if err != nil {
			return nil, err
		}
	}

	return thing, nil
}
