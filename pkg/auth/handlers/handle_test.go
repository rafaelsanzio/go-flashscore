package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAdapter(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ok", nil)

	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}

	testCases := []struct {
		Name          string
		HandleFunc    func(w http.ResponseWriter, r *http.Request)
		RateLimitFunc func() bool
	}{
		{
			Name:          "Handle function",
			HandleFunc:    handleFunc,
			RateLimitFunc: rateLimitAllow,
		}, {
			Name:          "Handle function with Rate Limit",
			HandleFunc:    handleFunc,
			RateLimitFunc: fakeRateLimitAllow,
		},
	}

	for _, tc := range testCases {
		rateLimitAllow = tc.RateLimitFunc
		defer restoreRateLimitAllow(rateLimitAllow)

		handle := HandleAdapter(tc.HandleFunc)

		handle.ServeHTTP(w, r)
	}
}
