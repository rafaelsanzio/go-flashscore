package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWTToken(t *testing.T) {
	tokenString, err := GenerateJWTToken("any_email")

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}
