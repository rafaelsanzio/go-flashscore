package repo

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

type MockTournamentRepo struct {
	tournament.TournamentRepo
	InsertFunc func(ctx context.Context, t tournament.Tournament) errs.AppError
	GetFunc    func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
	ListFunc   func(ctx context.Context) ([]tournament.Tournament, errs.AppError)
	UpdateFunc func(ctx context.Context, t tournament.Tournament) (*tournament.Tournament, errs.AppError)
	DeleteFunc func(ctx context.Context, id string) errs.AppError
}

func (m MockTournamentRepo) Insert(ctx context.Context, t tournament.Tournament) errs.AppError {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, t)
	}
	return m.TournamentRepo.Insert(ctx, t)
}

func (m MockTournamentRepo) Get(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return m.TournamentRepo.Get(ctx, id)
}

func (m MockTournamentRepo) List(ctx context.Context) ([]tournament.Tournament, errs.AppError) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return m.TournamentRepo.List(ctx)
}

func (m MockTournamentRepo) Update(ctx context.Context, t tournament.Tournament) (*tournament.Tournament, errs.AppError) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, t)
	}
	return m.TournamentRepo.Update(ctx, t)
}

func (m MockTournamentRepo) Delete(ctx context.Context, id string) errs.AppError {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return m.TournamentRepo.Delete(ctx, id)
}
