package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidMailAddress(t *testing.T) {
	testCases := []struct {
		Name         string
		Email        string
		ExpectResult bool
	}{
		{
			Name:         "Valid Email",
			Email:        "rafael@gmail.com",
			ExpectResult: true,
		}, {
			Name:         "Invalid Email",
			Email:        "rafaelgmail.com",
			ExpectResult: false,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)
		isValid := validMailAddress(tc.Email)
		assert.Equal(t, tc.ExpectResult, isValid)
	}
}
