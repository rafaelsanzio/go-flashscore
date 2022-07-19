package repo

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/apikey"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

type MockAPIKeyRepo struct {
	apikey.ApiKeyRepo
	InsertFunc      func(ctx context.Context, a apikey.ApiKey) errs.AppError
	GetFunc         func(ctx context.Context, id string) (*apikey.ApiKey, errs.AppError)
	ListFunc        func(ctx context.Context) ([]apikey.ApiKey, errs.AppError)
	UpdateFunc      func(ctx context.Context, a apikey.ApiKey) (*apikey.ApiKey, errs.AppError)
	DeleteFunc      func(ctx context.Context, id string) errs.AppError
	FindByEmailFunc func(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError)
}

func (m MockAPIKeyRepo) Insert(ctx context.Context, t apikey.ApiKey) errs.AppError {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, t)
	}
	return m.ApiKeyRepo.Insert(ctx, t)
}

func (m MockAPIKeyRepo) Get(ctx context.Context, id string) (*apikey.ApiKey, errs.AppError) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return m.ApiKeyRepo.Get(ctx, id)
}

func (m MockAPIKeyRepo) Update(ctx context.Context, t apikey.ApiKey) (*apikey.ApiKey, errs.AppError) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, t)
	}
	return m.ApiKeyRepo.Update(ctx, t)
}

func (m MockAPIKeyRepo) Delete(ctx context.Context, id string) errs.AppError {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return m.ApiKeyRepo.Delete(ctx, id)
}

func (m MockAPIKeyRepo) FindByEmail(ctx context.Context, email string) (*apikey.ApiKey, errs.AppError) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	return m.ApiKeyRepo.FindByEmail(ctx, email)
}
