package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/apikey"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
)

func TestApiKeyRepoInsert(t *testing.T) {
	ctx := context.Background()

	SetApiKeyRepo(MockAPIKeyRepo{
		InsertFunc: func(ctx context.Context, a apikey.ApiKey) errs.AppError {
			return nil
		},
	})
	defer SetApiKeyRepo(nil)

	newApiKey := prototype.PrototypeApiKey()

	err := GetApiKeyRepo().Insert(ctx, newApiKey)
	assert.NoError(t, err)
}

func TestApiKeyRepoGet(t *testing.T) {
	ctx := context.Background()

	SetApiKeyRepo(MockAPIKeyRepo{
		GetFunc: func(ctx context.Context, id string) (*apikey.ApiKey, errs.AppError) {
			apikey := prototype.PrototypeApiKey()
			return &apikey, nil
		},
	})
	defer SetApiKeyRepo(nil)

	newApiKey := prototype.PrototypeApiKey()

	result, err := GetApiKeyRepo().Get(ctx, "new-apikey-id")
	assert.NoError(t, err)

	assert.Equal(t, newApiKey, *result)
}

func TestApiKeyRepoUpdate(t *testing.T) {
	ctx := context.Background()

	SetApiKeyRepo(MockAPIKeyRepo{
		UpdateFunc: func(ctx context.Context, t apikey.ApiKey) (*apikey.ApiKey, errs.AppError) {
			return &t, nil
		},
	})
	defer SetApiKeyRepo(nil)

	newApiKey := prototype.PrototypeApiKey()

	apiKeyUpdated, err := GetApiKeyRepo().Update(ctx, newApiKey)
	assert.NoError(t, err)

	assert.Equal(t, newApiKey, *apiKeyUpdated)

}

func TestApiKeyRepoDelete(t *testing.T) {
	ctx := context.Background()

	SetApiKeyRepo(MockAPIKeyRepo{
		DeleteFunc: func(ctx context.Context, id string) errs.AppError {
			return nil
		},
	})
	defer SetApiKeyRepo(nil)

	newApiKey := prototype.PrototypeApiKey()

	err := GetApiKeyRepo().Delete(ctx, newApiKey.GetID())
	assert.NoError(t, err)
}

func TestApiKeyRepoFindByEmail(t *testing.T) {
	ctx := context.Background()

	SetApiKeyRepo(MockAPIKeyRepo{
		FindByEmailFunc: func(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
			apikey := prototype.PrototypeApiKey()
			return &apikey, nil
		},
	})
	defer SetApiKeyRepo(nil)

	newApiKey := prototype.PrototypeApiKey()

	result, err := GetApiKeyRepo().FindByEmail(ctx, "any_email")
	assert.NoError(t, err)

	assert.Equal(t, newApiKey, *result)
}
