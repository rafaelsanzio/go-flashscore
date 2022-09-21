package middleware

import (
	"fmt"
	"net/http"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/cache"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

func VerifyCache(w http.ResponseWriter, r *http.Request) {
	cacheKey := fmt.Sprintf("%s%s", r.Method, r.URL)

	value, err_ := cache.GetStore().Get(r.Context(), cacheKey)
	if value != nil && err_ == nil {
		data, err_ := jsonMarshal(value)
		if err_ != nil {
			_ = errs.ErrMarshalingJson.Throwf(applog.Log, errs.ErrFmt, err_)
			errs.HttpInternalServerError(w)
			return
		}

		_, err_ = write(w, data)
		if err_ != nil {
			_ = errs.ErrResponseWriter.Throwf(applog.Log, errs.ErrFmt, err_)
			errs.HttpInternalServerError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	if value == nil {
		applog.Log.Warnf("Cache does not exist for this key: %s", cacheKey)
	}
}
