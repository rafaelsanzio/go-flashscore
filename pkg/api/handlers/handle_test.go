package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/apikey"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
)

func mockFindByEmailApiKeyFunc(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	apikey := prototype.PrototypeApiKey()
	apikey.ValidUntil = time.Now().AddDate(1, 0, 0)
	return &apikey, nil
}

func TestHandleAdapter(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ok", nil)

	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodGet, "/teams", nil)
	r2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJhZmFlbHNhbnppb0BnbWFpbC5jb20ifQ.qtgF2cdlKxQQ-LV5ej43tarZHMXcylxmmSrvQZvtotc")

	handleFunc2 := func(w http.ResponseWriter, r2 *http.Request) {
		w2.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}

	testCases := []struct {
		Name           string
		ResponseWriter http.ResponseWriter
		Request        *http.Request
		HandleFunc     func(w http.ResponseWriter, r *http.Request)
		RateLimitFunc  func() bool
	}{
		{
			Name:           "Handle function",
			ResponseWriter: w,
			Request:        r,
			HandleFunc:     handleFunc,
			RateLimitFunc:  rateLimitAllow,
		}, {
			Name:           "Handle function with Rate Limit",
			ResponseWriter: w,
			Request:        r,
			HandleFunc:     handleFunc,
			RateLimitFunc:  fakeRateLimitAllow,
		}, {
			Name:           "Handle function with auth",
			ResponseWriter: w2,
			Request:        r2,
			HandleFunc:     handleFunc2,
			RateLimitFunc:  rateLimitAllow,
		},
	}

	for _, tc := range testCases {

		repo.SetApiKeyRepo(repo.MockAPIKeyRepo{
			FindByEmailFunc: mockFindByEmailApiKeyFunc,
		})
		defer repo.SetApiKeyRepo(nil)

		rateLimitAllow = tc.RateLimitFunc
		defer restoreRateLimitAllow(rateLimitAllow)

		handle := HandleAdapter(tc.HandleFunc)

		handle.ServeHTTP(tc.ResponseWriter, tc.Request)
	}
}
