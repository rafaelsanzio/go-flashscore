package repo

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/apikey"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/store"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
)

const (
	ApiKeyCollection = "apikey"
)

type apiKeyRepo struct {
	store store.Store
}

var apiKeySingleton apikey.ApiKeyRepo

func GetApiKeyRepo() apikey.ApiKeyRepo {
	if apiKeySingleton == nil {
		return getApiKeyRepo()
	}
	return apiKeySingleton
}

func getApiKeyRepo() *apiKeyRepo {
	s := store.GetStore()
	return &apiKeyRepo{s}
}

func SetApiKeyRepo(repo apikey.ApiKeyRepo) {
	apiKeySingleton = repo
}

func (repo apiKeyRepo) Insert(ctx context.Context, a apikey.ApiKey) errs.AppError {
	a.Created = time.Now()
	_, err := repo.store.InsertOne(ctx, ApiKeyCollection, &a)
	return err
}

func (repo apiKeyRepo) Get(ctx context.Context, id string) (*apikey.ApiKey, errs.AppError) {
	filter := query.Filter{
		"_id": id,
	}

	opts := query.FindOneOptions{}

	mApiKey := apikey.ApiKey{}
	err := repo.store.FindOne(ctx, ApiKeyCollection, filter, &mApiKey, opts)
	if err != nil {
		return &mApiKey, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", ApiKeyCollection, id, err)
	}

	return &mApiKey, nil
}

func (repo apiKeyRepo) Update(ctx context.Context, a apikey.ApiKey) (*apikey.ApiKey, errs.AppError) {
	res := apikey.ApiKey{}
	filter := query.Filter{
		"_id": a.GetID(),
	}

	err := repo.store.FindOne(ctx, ApiKeyCollection, filter, &res)
	if err != nil {
		return nil, err
	}

	if res.ID == "" {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", ApiKeyCollection, a.GetID(), err)
	}

	a.ID = res.ID
	err = repo.store.UpdateOne(ctx, ApiKeyCollection, &a)
	if err != nil {
		return nil, errs.ErrMongoUpdateOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", ApiKeyCollection, a.GetID(), err)
	}

	return &a, nil
}

func (repo apiKeyRepo) Delete(ctx context.Context, id string) errs.AppError {
	err := repo.store.DeleteOne(ctx, ApiKeyCollection, id)
	return err
}

func (repo apiKeyRepo) FindByEmail(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	filter := query.Filter{
		"email": email,
	}

	opts := query.FindOneOptions{}

	mApiKey := apikey.ApiKey{}
	err := repo.store.FindOne(ctx, ApiKeyCollection, filter, &mApiKey, opts)
	if err != nil {
		return &mApiKey, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", ApiKeyCollection, email, err)
	}

	if mApiKey.ID == "" {
		return nil, nil
	}

	return &mApiKey, nil
}
