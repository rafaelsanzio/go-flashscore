package repo

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
)

type MockPlayerRepo struct {
	player.PlayerRepo
	InsertFunc        func(ctx context.Context, p player.Player) errs.AppError
	GetFunc           func(ctx context.Context, id string) (*player.Player, errs.AppError)
	ListFunc          func(ctx context.Context) ([]player.Player, errs.AppError)
	UpdateFunc        func(ctx context.Context, p player.Player) (*player.Player, errs.AppError)
	DeleteFunc        func(ctx context.Context, id string) errs.AppError
	GetTeamPlayerFunc func(ctx context.Context, id, teamID string) (*player.Player, errs.AppError)
}

func (m MockPlayerRepo) Insert(ctx context.Context, p player.Player) errs.AppError {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, p)
	}
	return m.PlayerRepo.Insert(ctx, p)
}

func (m MockPlayerRepo) Get(ctx context.Context, id string) (*player.Player, errs.AppError) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return m.PlayerRepo.Get(ctx, id)
}

func (m MockPlayerRepo) List(ctx context.Context) ([]player.Player, errs.AppError) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return m.PlayerRepo.List(ctx)
}

func (m MockPlayerRepo) Update(ctx context.Context, p player.Player) (*player.Player, errs.AppError) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, p)
	}
	return m.PlayerRepo.Update(ctx, p)
}

func (m MockPlayerRepo) Delete(ctx context.Context, id string) errs.AppError {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return m.PlayerRepo.Delete(ctx, id)
}

func (m MockPlayerRepo) GetTeamPlayer(ctx context.Context, id, teamID string) (*player.Player, errs.AppError) {
	if m.GetTeamPlayerFunc != nil {
		return m.GetTeamPlayerFunc(ctx, id, teamID)
	}
	return m.PlayerRepo.GetTeamPlayer(ctx, id, teamID)
}
