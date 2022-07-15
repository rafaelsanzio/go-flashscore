package prototype

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/event"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
)

func PrototypeEvent() event.Event {

	return event.Event{
		ID:           "1",
		TournamentID: PrototypeTournament().ID,
		MatchID:      PrototypeMatch().ID,
		Type:         model.EventStart,
		Value:        nil,
	}
}
