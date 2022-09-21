package handlers

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/middleware"
)

var rateLimit = rate.NewLimiter(rate.Every(10*time.Second), 10) // 10 request every 10 seconds

func HandleAdapter(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		applog.Log.Infof("Requesting - Method: %s, URL %s", r.Method, r.URL)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if !rateLimitAllow() {
			applog.Log.Warnf("Rate limit to requests was exceed")
			errs.HttpToManyRequests(w)
			return
		}

		err := middleware.IsAuthorized(r.Context(), r.Header.Get("Authorization"))
		if err != nil {
			applog.Log.Warnf("API Key Token is not valid: %s", err.Error())
			errs.HttpUnauthorized(w)
			return
		}

		middleware.VerifyCache(w, r)

		hf(w, r)
	}
}
