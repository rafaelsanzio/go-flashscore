package errs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpUnprocessableEntity(t *testing.T) {
	w := httptest.NewRecorder()
	HttpUnprocessableEntity(w, "any message")

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	write = fakeWrite
	defer restoreWrite(write)

	HttpUnprocessableEntity(w, "any message error")
}

func TestHttpInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	HttpInternalServerError(w)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHttpNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	HttpNotFound(w)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHttpToManyRequests(t *testing.T) {
	w := httptest.NewRecorder()
	HttpToManyRequests(w)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestHttpUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	HttpUnauthorized(w)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
