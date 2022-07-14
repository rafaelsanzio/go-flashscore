package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
)

func TestMatchRepoInsert(t *testing.T) {
	ctx := context.Background()

	SetMatchRepo(MockMatchRepo{
		InsertFunc: func(ctx context.Context, mt match.Match) errs.AppError {
			return nil
		},
	})
	defer SetMatchRepo(nil)

	newMatch := prototype.PrototypeMatch()

	err := GetMatchRepo().Insert(ctx, newMatch)
	assert.NoError(t, err)
}

func TestMatchRepoGet(t *testing.T) {
	ctx := context.Background()

	SetMatchRepo(MockMatchRepo{
		GetFunc: func(ctx context.Context, id string) (*match.Match, errs.AppError) {
			match := prototype.PrototypeMatch()
			return &match, nil
		},
	})
	defer SetMatchRepo(nil)

	newMatch := prototype.PrototypeMatch()

	result, err := GetMatchRepo().Get(ctx, "new-match-id")
	assert.NoError(t, err)

	assert.Equal(t, newMatch, *result)
}

func TestMatchRepoList(t *testing.T) {
	ctx := context.Background()

	SetMatchRepo(MockMatchRepo{
		ListFunc: func(ctx context.Context) ([]match.Match, errs.AppError) {
			matchMock := prototype.PrototypeMatch()
			matchMock2 := prototype.PrototypeMatch()

			return []match.Match{matchMock, matchMock2}, nil
		},
	})
	defer SetMatchRepo(nil)

	matches, err := GetMatchRepo().List(ctx)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(matches))
}

func TestMatchRepoUpdate(t *testing.T) {
	ctx := context.Background()

	SetMatchRepo(MockMatchRepo{
		UpdateFunc: func(ctx context.Context, p match.Match) (*match.Match, errs.AppError) {
			return &p, nil
		},
	})
	defer SetMatchRepo(nil)

	newMatch := prototype.PrototypeMatch()

	matchUpdated, err := GetMatchRepo().Update(ctx, newMatch)
	assert.NoError(t, err)

	assert.Equal(t, newMatch, *matchUpdated)

}

func TestMatchRepoDelete(t *testing.T) {
	ctx := context.Background()

	SetMatchRepo(MockMatchRepo{
		DeleteFunc: func(ctx context.Context, id string) errs.AppError {
			return nil
		},
	})
	defer SetMatchRepo(nil)

	newMatch := prototype.PrototypeMatch()

	err := GetMatchRepo().Delete(ctx, newMatch.GetID())
	assert.NoError(t, err)
}

func TestMatchRepoFindMatchForTournament(t *testing.T) {
	ctx := context.Background()

	SetMatchRepo(MockMatchRepo{
		FindMatchForTournamentFunc: func(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError) {
			match := prototype.PrototypeMatch()
			return &match, nil
		},
	})
	defer SetMatchRepo(nil)

	newMatch := prototype.PrototypeMatch()

	result, err := GetMatchRepo().FindMatchForTournament(ctx, "match-id", "tournament-id")
	assert.NoError(t, err)

	assert.Equal(t, newMatch, *result)
}
