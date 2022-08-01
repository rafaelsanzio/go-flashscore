package notification

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

func mockGetTeamFunc(ctx context.Context, id string) (*team.Team, errs.AppError) {
	teamMock := prototype.PrototypeTeam()
	return &teamMock, nil
}

func TestHandler(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		Name                   string
		Body                   string
		HandleGetPlayerFunc    func(ctx context.Context, id string) (*player.Player, errs.AppError)
		HandleUpdatePlayerFunc func(ctx context.Context, t player.Player) (*player.Player, errs.AppError)
		HandleGetTeamFunc      func(ctx context.Context, id string) (*team.Team, errs.AppError)
		ExpectedError          bool
	}{
		{
			Name:                   "Handle update team player correct",
			Body:                   `{"Action":"UpdateTeamPlayer","Data":{"playerID":"any-player-id","teamDestinyID":"any-team-id"}}`,
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			ExpectedError:          false,
		}, {
			Name:                   "Handle update team player throw error on get player function",
			Body:                   `{"Action":"UpdateTeamPlayer","Data":{"playerID":"any-player-id","teamDestinyID":"any-team-id"}}`,
			HandleGetPlayerFunc:    mockGetPlayerThrowFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			ExpectedError:          true,
		}, {
			Name:                   "Handle update team player throw error on unmarshal function",
			Body:                   `{"Action":"UpdateTeamPlayer","Data":{"playerID":"any-player-id","teamDestinyID":"any-team-id"}`,
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			ExpectedError:          true,
		}, {
			Name:                   "Handle update team player throw error on non exist action",
			Body:                   `{"Action":"Unknown","Data":{"playerID":"any-player-id","teamDestinyID":"any-team-id"}}`,
			HandleGetPlayerFunc:    mockGetPlayerFunc,
			HandleUpdatePlayerFunc: mockUpdatePlayerFunc,
			HandleGetTeamFunc:      mockGetTeamFunc,
			ExpectedError:          true,
		},
	}

	for _, tc := range testCases {

		repo.SetPlayerRepo(repo.MockPlayerRepo{
			GetFunc:    tc.HandleGetPlayerFunc,
			UpdateFunc: tc.HandleUpdatePlayerFunc,
		})
		defer repo.SetPlayerRepo(nil)

		repo.SetTeamRepo(repo.MockTeamRepo{
			GetFunc: tc.HandleGetTeamFunc,
		})
		defer repo.SetTeamRepo(nil)

		err := Handler(ctx, tc.Body, "any-key")
		if tc.ExpectedError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}
