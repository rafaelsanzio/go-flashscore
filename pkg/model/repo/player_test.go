package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
)

func TestPlayerRepoInsert(t *testing.T) {
	ctx := context.Background()

	SetPlayerRepo(MockPlayerRepo{
		InsertFunc: func(ctx context.Context, p player.Player) errs.AppError {
			return nil
		},
	})
	defer SetPlayerRepo(nil)

	newPlayer := prototype.PrototypePlayer()

	err := GetPlayerRepo().Insert(ctx, newPlayer)
	assert.NoError(t, err)
}

func TestPlayerRepoGet(t *testing.T) {
	ctx := context.Background()

	SetPlayerRepo(MockPlayerRepo{
		GetFunc: func(ctx context.Context, id string) (*player.Player, errs.AppError) {
			player := prototype.PrototypePlayer()
			return &player, nil
		},
	})
	defer SetPlayerRepo(nil)

	newPlayer := prototype.PrototypePlayer()

	result, err := GetPlayerRepo().Get(ctx, "new-player-id")
	assert.NoError(t, err)

	assert.Equal(t, newPlayer, *result)
}

func TestPlayerRepoList(t *testing.T) {
	ctx := context.Background()

	SetPlayerRepo(MockPlayerRepo{
		ListFunc: func(ctx context.Context) ([]player.Player, errs.AppError) {
			playerMock := prototype.PrototypePlayer()
			playerMock2 := prototype.PrototypePlayer()

			return []player.Player{playerMock, playerMock2}, nil
		},
	})
	defer SetPlayerRepo(nil)

	players, err := GetPlayerRepo().List(ctx)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(players))
}

func TestPlayerRepoUpdate(t *testing.T) {
	ctx := context.Background()

	SetPlayerRepo(MockPlayerRepo{
		UpdateFunc: func(ctx context.Context, p player.Player) (*player.Player, errs.AppError) {
			return &p, nil
		},
	})
	defer SetPlayerRepo(nil)

	newPlayer := prototype.PrototypePlayer()

	playerUpdated, err := GetPlayerRepo().Update(ctx, newPlayer)
	assert.NoError(t, err)

	assert.Equal(t, newPlayer, *playerUpdated)

}

func TestPlayerRepoDelete(t *testing.T) {
	ctx := context.Background()

	SetPlayerRepo(MockPlayerRepo{
		DeleteFunc: func(ctx context.Context, id string) errs.AppError {
			return nil
		},
	})
	defer SetPlayerRepo(nil)

	newPlayer := prototype.PrototypePlayer()

	err := GetPlayerRepo().Delete(ctx, newPlayer.GetID())
	assert.NoError(t, err)
}

func TestPlayerRepoGetTeamPlayer(t *testing.T) {
	ctx := context.Background()

	SetPlayerRepo(MockPlayerRepo{
		GetTeamPlayerFunc: func(ctx context.Context, id, teamID string) (*player.Player, errs.AppError) {
			player := prototype.PrototypePlayer()
			return &player, nil
		},
	})
	defer SetPlayerRepo(nil)

	newPlayer := prototype.PrototypePlayer()

	result, err := GetPlayerRepo().GetTeamPlayer(ctx, "new-player-id", "new-team-id")
	assert.NoError(t, err)

	assert.Equal(t, newPlayer, *result)
}
