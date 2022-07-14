package repo

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
)

type MockMatchRepo struct {
	match.MatchRepo
	InsertFunc func(ctx context.Context, mt match.Match) errs.AppError
	GetFunc    func(ctx context.Context, id string) (*match.Match, errs.AppError)
	ListFunc   func(ctx context.Context) ([]match.Match, errs.AppError)
	UpdateFunc func(ctx context.Context, mt match.Match) (*match.Match, errs.AppError)
	DeleteFunc func(ctx context.Context, id string) errs.AppError

	FindMatchForTournamentFunc func(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError)
	FindTeamInMatchFunc        func(ctx context.Context, teamID string) (bool, errs.AppError)
}

func (m MockMatchRepo) Insert(ctx context.Context, mt match.Match) errs.AppError {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, mt)
	}
	return m.MatchRepo.Insert(ctx, mt)
}

func (m MockMatchRepo) Get(ctx context.Context, id string) (*match.Match, errs.AppError) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return m.MatchRepo.Get(ctx, id)
}

func (m MockMatchRepo) List(ctx context.Context) ([]match.Match, errs.AppError) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return m.MatchRepo.List(ctx)
}

func (m MockMatchRepo) Update(ctx context.Context, mt match.Match) (*match.Match, errs.AppError) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, mt)
	}
	return m.MatchRepo.Update(ctx, mt)
}

func (m MockMatchRepo) Delete(ctx context.Context, id string) errs.AppError {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return m.MatchRepo.Delete(ctx, id)
}

func (m MockMatchRepo) FindMatchForTournament(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError) {
	if m.FindMatchForTournamentFunc != nil {
		return m.FindMatchForTournamentFunc(ctx, id, tournamentID)
	}
	return m.MatchRepo.FindMatchForTournament(ctx, id, tournamentID)
}
