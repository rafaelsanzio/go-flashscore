package match

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

type MatchRepo interface {
	Insert(ctx context.Context, m Match) errs.AppError
	Get(ctx context.Context, id string) (*Match, errs.AppError)
	List(ctx context.Context) ([]Match, errs.AppError)
	Update(ctx context.Context, m Match) (*Match, errs.AppError)
	Delete(ctx context.Context, id string) errs.AppError

	FindMatchForTournament(ctx context.Context, id, tournamentID string) (*Match, errs.AppError)
}

type Match struct {
	ID          string `bson:"_id"`
	Tournament  tournament.Tournament
	HomeTeam    team.Team
	AwayTeam    team.Team
	DateOfMatch string
	TimeOfMatch string
	Status      model.MatchStatus
	Events      []interface{}
	Created     time.Time
}

func (mt Match) GetID() string {
	return mt.ID
}

func (mt *Match) SetID(id string) {
	mt.ID = id
}

func (mt *Match) FindTeamInMatch(teamID string) bool {
	if teamID == mt.HomeTeam.ID || teamID == mt.AwayTeam.ID {
		return true
	}
	return false
}

func (mt *Match) IsTheMatchForTournament(tournamentID string) bool {
	return mt.Tournament.ID == tournamentID
}
