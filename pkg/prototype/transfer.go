package prototype

import (
	"github.com/rafaelsanzio/go-flashscore/pkg/money"
	"github.com/rafaelsanzio/go-flashscore/pkg/transfer"
)

func PrototypeTransfer() transfer.Transfer {
	return transfer.Transfer{
		ID:             "1",
		Player:         PrototypePlayer(),
		TeamDestiny:    PrototypeTeam(),
		Amount:         money.Money{Cents: 1000, Currency: money.USD},
		DateOfTransfer: "1990-01-01",
	}
}
