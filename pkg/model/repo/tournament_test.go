package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func TestTournamentRepoInsert(t *testing.T) {
	ctx := context.Background()

	SetTournamentRepo(MockTournamentRepo{
		InsertFunc: func(ctx context.Context, t tournament.Tournament) errs.AppError {
			return nil
		},
	})
	defer SetTournamentRepo(nil)

	newTournament := prototype.PrototypeTournament()

	err := GetTournamentRepo().Insert(ctx, newTournament)
	assert.NoError(t, err)
}

func TestTournamentRepoGet(t *testing.T) {
	ctx := context.Background()

	SetTournamentRepo(MockTournamentRepo{
		GetFunc: func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
			team := prototype.PrototypeTournament()
			return &team, nil
		},
	})
	defer SetTournamentRepo(nil)

	newTournament := prototype.PrototypeTournament()

	result, err := GetTournamentRepo().Get(ctx, "new-tournament-id")
	assert.NoError(t, err)

	assert.Equal(t, newTournament, *result)
}

func TestTournamentRepoList(t *testing.T) {
	ctx := context.Background()

	SetTournamentRepo(MockTournamentRepo{
		ListFunc: func(ctx context.Context) ([]tournament.Tournament, errs.AppError) {
			tournamentMock := prototype.PrototypeTournament()
			tournamentMock2 := prototype.PrototypeTournament()

			return []tournament.Tournament{tournamentMock, tournamentMock2}, nil
		},
	})
	defer SetTournamentRepo(nil)

	teams, err := GetTournamentRepo().List(ctx)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(teams))
}

func TestTournamentRepoUpdate(t *testing.T) {
	ctx := context.Background()

	SetTournamentRepo(MockTournamentRepo{
		UpdateFunc: func(ctx context.Context, t tournament.Tournament) (*tournament.Tournament, errs.AppError) {
			return &t, nil
		},
	})
	defer SetTournamentRepo(nil)

	newTournament := prototype.PrototypeTournament()

	teamUpdated, err := GetTournamentRepo().Update(ctx, newTournament)
	assert.NoError(t, err)

	assert.Equal(t, newTournament, *teamUpdated)

}

func TestTournamentRepoDelete(t *testing.T) {
	ctx := context.Background()

	SetTournamentRepo(MockTournamentRepo{
		DeleteFunc: func(ctx context.Context, id string) errs.AppError {
			return nil
		},
	})
	defer SetTournamentRepo(nil)

	newTournament := prototype.PrototypeTournament()

	err := GetTournamentRepo().Delete(ctx, newTournament.GetID())
	assert.NoError(t, err)
}
