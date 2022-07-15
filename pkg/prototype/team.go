package prototype

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

func PrototypeTeam() team.Team {
	return team.Team{
		ID:        "1",
		Name:      "Real Madrid Club",
		ShortCode: "RMC",
		Country:   "Spain",
		City:      "Madrid",
	}
}
