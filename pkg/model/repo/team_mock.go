package repo

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

type MockTeamRepo struct {
	team.TeamRepo
	InsertFunc func(ctx context.Context, t team.Team) errs.AppError
	GetFunc    func(ctx context.Context, id string) (*team.Team, errs.AppError)
	ListFunc   func(ctx context.Context) ([]team.Team, errs.AppError)
	UpdateFunc func(ctx context.Context, t team.Team) (*team.Team, errs.AppError)
	DeleteFunc func(ctx context.Context, id string) errs.AppError
}

func (m MockTeamRepo) Insert(ctx context.Context, t team.Team) errs.AppError {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, t)
	}
	return m.TeamRepo.Insert(ctx, t)
}

func (m MockTeamRepo) Get(ctx context.Context, id string) (*team.Team, errs.AppError) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return m.TeamRepo.Get(ctx, id)
}

func (m MockTeamRepo) List(ctx context.Context) ([]team.Team, errs.AppError) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return m.TeamRepo.List(ctx)
}

func (m MockTeamRepo) Update(ctx context.Context, t team.Team) (*team.Team, errs.AppError) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, t)
	}
	return m.TeamRepo.Update(ctx, t)
}

func (m MockTeamRepo) Delete(ctx context.Context, id string) errs.AppError {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return m.TeamRepo.Delete(ctx, id)
}
