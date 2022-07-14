package prototype

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
)

func PrototypeMatch() match.Match {
	return match.Match{
		ID:          "1",
		Tournament:  PrototypeTournament(),
		HomeTeam:    PrototypeTeam(),
		AwayTeam:    PrototypeTeam(),
		DateOfMatch: "2022-02-01",
		TimeOfMatch: "16:00",
		Status:      model.MatchStatusNotStart,
		Events:      nil,
	}
}
