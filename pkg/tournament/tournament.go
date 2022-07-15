package tournament

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

type TournamentRepo interface {
	Insert(ctx context.Context, t Tournament) errs.AppError
	Get(ctx context.Context, id string) (*Tournament, errs.AppError)
	List(ctx context.Context) ([]Tournament, errs.AppError)
	Update(ctx context.Context, t Tournament) (*Tournament, errs.AppError)
	Delete(ctx context.Context, id string) errs.AppError
}

type Tournament struct {
	ID      string `bson:"_id"`
	Name    string
	Teams   []team.Team
	Created time.Time
}

func (t Tournament) GetID() string {
	return t.ID
}

func (t *Tournament) SetID(id string) {
	t.ID = id
}

func (t *Tournament) GetTeams() ([]team.Team, errs.AppError) {
	return t.Teams, nil
}

func (t *Tournament) FindTeam(teamID string) bool {
	for _, team := range t.Teams {
		if team.ID == teamID {
			return true
		}
	}
	return false
}
