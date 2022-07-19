package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func TestHandleEventMatchExtratime(t *testing.T) {
	ctx := context.Background()

	data := map[string]string{
		"tournamentID": "any-player-id",
		"matchID":      "any-match-id",
		"extratime":    "5",
	}

	testCases := []struct {
		Name                             string
		HandleGetTournamentFunc          func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		HandleUpdateMatchFunc            func(ctx context.Context, mt match.Match) (*match.Match, errs.AppError)
		HandleFindMatchForTournamentFunc func(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError)
		StrconvAtoiFunc                  func(s string) (int, error)
		ExpectedError                    bool
	}{
		{
			Name:                             "Handle event match extratime correct",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    false,
		}, {
			Name:                             "Handle event match extratime throw error on get tournament function",
			HandleGetTournamentFunc:          mockGetTournamentThrowFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle event match extratime throw error on update match function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchThrowFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle event match extratime throw error on find match fot tournament function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentThrowFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle event match extratime throw error on strconv Atoi function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			StrconvAtoiFunc:                  fakeStrconvAtoi,
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

		strconvAtoi = tc.StrconvAtoiFunc
		defer restoreStrconvAtoi(strconvAtoi)

		err := HandleEventMatchExtratime(ctx, data)
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}
