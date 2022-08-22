package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyCache(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ok", nil)

	VerifyCache(w, r)
}
