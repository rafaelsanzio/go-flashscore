package prototype

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

func PrototypeTournament() tournament.Tournament {
	return tournament.Tournament{
		ID:    "1",
		Name:  "any-tournament-name",
		Teams: []team.Team{PrototypeTeam(), PrototypeTeam(), PrototypeTeam()},
	}
}
