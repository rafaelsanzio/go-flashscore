package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockGetTournamentFunc(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
	tournamentMock := prototype.PrototypeTournament()
	return &tournamentMock, nil
}

func mockGetTournamentThrowFunc(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockFindMatchForTournamentFunc(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError) {
	matchMock := prototype.PrototypeMatch()
	return &matchMock, nil
}

func mockFindMatchForTournamentThrowFunc(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockUpdateMatchFunc(ctx context.Context, mt match.Match) (*match.Match, errs.AppError) {
	return &mt, nil
}

func mockUpdateMatchThrowFunc(ctx context.Context, mt match.Match) (*match.Match, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleEventMatchStart(t *testing.T) {
	ctx := context.Background()

	data := map[string]string{
		"tournamentID": "any-player-id",
		"matchID":      "any-match-id",
		"timeStarted":  "16:00",
	}

	testCases := []struct {
		Name                             string
		HandleGetTournamentFunc          func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		HandleUpdateMatchFunc            func(ctx context.Context, mt match.Match) (*match.Match, errs.AppError)
		HandleFindMatchForTournamentFunc func(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError)
		ExpectedError                    bool
	}{
		{
			Name:                             "Handle post match start correct",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			ExpectedError:                    false,
		}, {
			Name:                             "Handle post match start throw error on get tournament function",
			HandleGetTournamentFunc:          mockGetTournamentThrowFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle post match start throw error on update match function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchThrowFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle post match start throw error on find match fot tournament function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentThrowFunc,
			ExpectedError:                    true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetTournamentRepo(repo.MockTournamentRepo{
			GetFunc: tc.HandleGetTournamentFunc,
		})
		defer repo.SetTournamentRepo(nil)

		repo.SetMatchRepo(repo.MockMatchRepo{
			UpdateFunc:                 tc.HandleUpdateMatchFunc,
			FindMatchForTournamentFunc: tc.HandleFindMatchForTournamentFunc,
		})
		defer repo.SetMatchRepo(nil)

		err := HandleEventMatchStart(ctx, data)
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}
