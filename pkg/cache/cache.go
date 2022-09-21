package cache

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
)

func SetCache(ctx context.Context, cacheKey string, data []byte) {
	_err := GetStore().Set(ctx, cacheKey, data)
	if _err != nil {
		applog.Log.Warnf("Cache could not be set for this key: %s with error: %s", cacheKey, _err.Error())
	}
}
