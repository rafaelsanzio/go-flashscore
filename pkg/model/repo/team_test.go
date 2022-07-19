package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func TestTeamRepoInsert(t *testing.T) {
	ctx := context.Background()

	SetTeamRepo(MockTeamRepo{
		InsertFunc: func(ctx context.Context, t team.Team) errs.AppError {
			return nil
		},
	})
	defer SetTeamRepo(nil)

	newTeam := prototype.PrototypeTeam()

	err := GetTeamRepo().Insert(ctx, newTeam)
	assert.NoError(t, err)
}

func TestTeamRepoGet(t *testing.T) {
	ctx := context.Background()

	SetTeamRepo(MockTeamRepo{
		GetFunc: func(ctx context.Context, id string) (*team.Team, errs.AppError) {
			team := prototype.PrototypeTeam()
			return &team, nil
		},
	})
	defer SetTeamRepo(nil)

	newTeam := prototype.PrototypeTeam()

	result, err := GetTeamRepo().Get(ctx, "new-team-id")
	assert.NoError(t, err)

	assert.Equal(t, newTeam, *result)
}

func TestTeamRepoList(t *testing.T) {
	ctx := context.Background()

	SetTeamRepo(MockTeamRepo{
		ListFunc: func(ctx context.Context) ([]team.Team, errs.AppError) {
			teamMock := prototype.PrototypeTeam()
			teamMock2 := prototype.PrototypeTeam()

			return []team.Team{teamMock, teamMock2}, nil
		},
	})
	defer SetTeamRepo(nil)

	teams, err := GetTeamRepo().List(ctx)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(teams))
}

func TestTeamRepoUpdate(t *testing.T) {
	ctx := context.Background()

	SetTeamRepo(MockTeamRepo{
		UpdateFunc: func(ctx context.Context, t team.Team) (*team.Team, errs.AppError) {
			return &t, nil
		},
	})
	defer SetTeamRepo(nil)

	newTeam := prototype.PrototypeTeam()

	teamUpdated, err := GetTeamRepo().Update(ctx, newTeam)
	assert.NoError(t, err)

	assert.Equal(t, newTeam, *teamUpdated)

}

func TestTeamRepoDelete(t *testing.T) {
	ctx := context.Background()

	SetTeamRepo(MockTeamRepo{
		DeleteFunc: func(ctx context.Context, id string) errs.AppError {
			return nil
		},
	})
	defer SetTeamRepo(nil)

	newTeam := prototype.PrototypeTeam()

	err := GetTeamRepo().Delete(ctx, newTeam.GetID())
	assert.NoError(t, err)
}
