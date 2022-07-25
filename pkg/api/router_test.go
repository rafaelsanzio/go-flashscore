package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	server := httptest.NewServer(router)

	healthCheckURL := fmt.Sprintf("%s/ok", server.URL)
	notFoundURL := fmt.Sprintf("%s/a-non-existent-path", server.URL)

	res, err := http.DefaultClient.Get(healthCheckURL)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	res, err = http.DefaultClient.Get(notFoundURL)
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
}
