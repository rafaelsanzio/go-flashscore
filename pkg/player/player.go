package player

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

type PlayerRepo interface {
	Insert(ctx context.Context, p Player) errs.AppError
	Get(ctx context.Context, id string) (*Player, errs.AppError)
	List(ctx context.Context) ([]Player, errs.AppError)
	Update(ctx context.Context, p Player) (*Player, errs.AppError)
	Delete(ctx context.Context, id string) errs.AppError

	GetTeamPlayer(ctx context.Context, id, teamID string) (*Player, errs.AppError)
}

type Player struct {
	ID           string `bson:"_id"`
	Name         string
	Country      string
	BirthdayDate string
	Team         team.Team
	Created      time.Time
}

func (p Player) GetID() string {
	return p.ID
}

func (p *Player) SetID(id string) {
	p.ID = id
}
