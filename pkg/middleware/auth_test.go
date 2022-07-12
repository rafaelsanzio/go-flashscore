package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/apikey"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
)

func mockFindByEmailApiKeyFunc(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	apikey := prototype.PrototypeApiKey()
	apikey.ValidUntil = time.Now().AddDate(1, 0, 0)
	return &apikey, nil
}

func mockFindByEmailApiKeyWithTokenInvalidFunc(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	apikey := prototype.PrototypeApiKey()
	apikey.ValidUntil = time.Now().AddDate(-1, 0, 0)
	return &apikey, nil
}

func mockFindByEmailApiKeyEmptyFunc(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	return nil, nil
}

func mockFindByEmailApiKeyThrowFunc(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestIsAuthorized(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		Name                  string
		Token                 string
		HandleFindByEmailFunc func(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError)
		ConfigValueFunc       func(k key.Key) (string, errs.AppError)
		ExpectedError         bool
	}{
		{
			Name:                  "Should authorized",
			Token:                 "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJhZmFlbEBnbWFpbC5jb20iLCJleHAiOjE2NTc5OTI1NTF9.vr2-W08K3OxTqtuhm2C344JVbwDGLB1sniU9cmNfABc",
			HandleFindByEmailFunc: mockFindByEmailApiKeyFunc,
			ConfigValueFunc:       configValue,
			ExpectedError:         false,
		}, {
			Name:                  "Should throw an error on find by email function",
			Token:                 "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJhZmFlbEBnbWFpbC5jb20iLCJleHAiOjE2NTc5OTI1NTF9.vr2-W08K3OxTqtuhm2C344JVbwDGLB1sniU9cmNfABc",
			HandleFindByEmailFunc: mockFindByEmailApiKeyThrowFunc,
			ConfigValueFunc:       configValue,
			ExpectedError:         true,
		}, {
			Name:                  "Should unauthorized with a non exist email",
			Token:                 "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJhZmFlbEBnbWFpbC5jb20iLCJleHAiOjE2NTc5OTI1NTF9.vr2-W08K3OxTqtuhm2C344JVbwDGLB1sniU9cmNfABc",
			HandleFindByEmailFunc: mockFindByEmailApiKeyEmptyFunc,
			ConfigValueFunc:       configValue,
			ExpectedError:         true,
		}, {
			Name:                  "Should unauthorized with valid until invalid",
			Token:                 "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJhZmFlbEBnbWFpbC5jb20iLCJleHAiOjE2NTc5OTI1NTF9.vr2-W08K3OxTqtuhm2C344JVbwDGLB1sniU9cmNfABc",
			HandleFindByEmailFunc: mockFindByEmailApiKeyWithTokenInvalidFunc,
			ConfigValueFunc:       configValue,
			ExpectedError:         true,
		}, {
			Name:                  "Should unauthorized with token empty",
			Token:                 "",
			HandleFindByEmailFunc: mockFindByEmailApiKeyFunc,
			ConfigValueFunc:       configValue,
			ExpectedError:         true,
		}, {
			Name:                  "Should unauthorized with wrong token",
			Token:                 "Bearer any_token",
			HandleFindByEmailFunc: mockFindByEmailApiKeyFunc,
			ConfigValueFunc:       configValue,
			ExpectedError:         true,
		}, {
			Name:                  "Should return an error on config value function",
			Token:                 "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJhZmFlbEBnbWFpbC5jb20iLCJleHAiOjE2NTc5OTI1NTF9.vr2-W08K3OxTqtuhm2C344JVbwDGLB1sniU9cmNfABc",
			HandleFindByEmailFunc: mockFindByEmailApiKeyFunc,
			ConfigValueFunc:       fakeConfigValue,
			ExpectedError:         true,
		}, {
			Name:                  "Should return an error on invalid key of signature token",
			Token:                 "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJhZmFlbEBnbWFpbC5jb20iLCJleHAiOjE2NTc5OTI1NTF9.vr2-W08K3OxTqtuhm2C344JVbwDGLB1sniU9cmNfABc",
			HandleFindByEmailFunc: mockFindByEmailApiKeyFunc,
			ConfigValueFunc:       fakeConfigValueWithAnyKey,
			ExpectedError:         true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetApiKeyRepo(repo.MockAPIKeyRepo{
			FindByEmailFunc: tc.HandleFindByEmailFunc,
		})
		defer repo.SetApiKeyRepo(nil)

		configValue = tc.ConfigValueFunc
		defer restoreConfigValue(configValue)

		err := IsAuthorized(ctx, tc.Token)
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
