package prototype

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
)

func PrototypePlayer() player.Player {
	return player.Player{
		ID:           "1",
		Name:         "Cristiano Ronaldo",
		Team:         PrototypeTeam(),
		Country:      "Portugal",
		BirthdayDate: "1990-01-01",
	}
}
