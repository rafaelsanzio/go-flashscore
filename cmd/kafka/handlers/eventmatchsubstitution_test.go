package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func mockGetTeamPlayerInThrowFunc(ctx context.Context, id, teamID string) (*player.Player, errs.AppError) {
	if id == "any-player-in-id" {
		return nil, errs.ErrRepoMockAction
	}
	playerMock := prototype.PrototypePlayer()
	return &playerMock, nil
}

func TestHandleEventMatchSubstitution(t *testing.T) {
	ctx := context.Background()

	data := map[string]string{
		"matchEventType":     "Substitution",
		"tournamentID":       "any-tournament-id",
		"matchID":            "any-match-id",
		"teamID":             "any-team-id",
		"playerOutID":        "any-player-out-id",
		"playerInID":         "any-player-in-id",
		"substitutionMinute": "10",
	}

	testCases := []struct {
		Name                             string
		HandleGetTournamentFunc          func(ctx context.Context, id string) (*tournament.Tournament, errs.AppError)
		HandleUpdateMatchFunc            func(ctx context.Context, mt match.Match) (*match.Match, errs.AppError)
		HandleFindMatchForTournamentFunc func(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError)
		HandleGetTeamFunc                func(ctx context.Context, id string) (*team.Team, errs.AppError)
		HandleGetTeamPlayerFunc          func(ctx context.Context, id, teamID string) (*player.Player, errs.AppError)
		StrconvAtoiFunc                  func(s string) (int, error)
		ExpectedError                    bool
	}{
		{
			Name:                             "Handle event match goal correct",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    false,
		}, {
			Name:                             "Handle event match goal throw error on get tournament function",
			HandleGetTournamentFunc:          mockGetTournamentThrowFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle event match goal throw error on update match function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchThrowFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle event match goal throw error on find match fot tournament function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentThrowFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle event match goal throw error on get team function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:                mockGetTeamThrowFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle event match goal throw error on get team player function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerThrowFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle event match goal throw error on get team player function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerInThrowFunc,
			StrconvAtoiFunc:                  strconvAtoi,
			ExpectedError:                    true,
		}, {
			Name:                             "Handle event match goal throw error on strconv Atoi function",
			HandleGetTournamentFunc:          mockGetTournamentFunc,
			HandleUpdateMatchFunc:            mockUpdateMatchFunc,
			HandleFindMatchForTournamentFunc: mockFindMatchForTournamentFunc,
			HandleGetTeamFunc:                mockGetTeamFunc,
			HandleGetTeamPlayerFunc:          mockGetTeamPlayerFunc,
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

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetTeamPlayerFunc: tc.HandleGetTeamPlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		strconvAtoi = tc.StrconvAtoiFunc
		defer restoreStrconvAtoi(strconvAtoi)

		err := HandleEventMatchSubstitution(ctx, data)
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}
