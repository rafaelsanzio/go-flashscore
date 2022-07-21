package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func mockGetPlayerFunc(ctx context.Context, id string) (*player.Player, errs.AppError) {
	playerMock := prototype.PrototypePlayer()
	return &playerMock, nil
}

func mockGetPlayerThrowFunc(ctx context.Context, id string) (*player.Player, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockUpdatePlayerFunc(ctx context.Context, p player.Player) (*player.Player, errs.AppError) {
	return &p, nil
}

func mockUpdatePlayerThrowFunc(ctx context.Context, p player.Player) (*player.Player, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func mockGetTeamFunc(ctx context.Context, id string) (*team.Team, errs.AppError) {
	teamMock := prototype.PrototypeTeam()
	return &teamMock, nil
}

func mockGetTeamThrowFunc(ctx context.Context, id string) (*team.Team, errs.AppError) {
	return nil, errs.ErrRepoMockAction
}

func TestHandleEventUpdateTeamPlayer(t *testing.T) {
	ctx := context.Background()

	data := map[string]string{
		"playerID":      "any-player-id",
		"teamDestinyID": "any-team-id",
	}

	testCases := []struct {
		Name                   string
		HandleGetPlayerFunc    func(ctx context.Context, id string) (*player.Player, errs.AppError)
		HandleUpdatePlayerFunc func(ctx context.Context, t player.Player) (*player.Player, errs.AppError)
		HandleGetTeamFunc      func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ExpectedError          bool
	}{
		{
			Name:                   "Handle update team player correct",
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			ExpectedError:          false,
		}, {
			Name:                   "Handle update team player throw error on get player function",
			HandleGetPlayerFunc:    mockGetPlayerThrowFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			ExpectedError:          true,
		}, {
			Name:                   "Handle update team player throw error on update player function",
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerThrowFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			ExpectedError:          true,
		}, {
			Name:                   "Handle update team player throw error on get team function",
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			HandleGetTeamFunc:      mockGetTeamThrowFunc,
			ExpectedError:          true,
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.Name)

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetFunc:    tc.HandleGetPlayerFunc,
			UpdateFunc: tc.HandleUpdatePlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		err := HandleEventUpdateTeamPlayer(ctx, data)
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}
