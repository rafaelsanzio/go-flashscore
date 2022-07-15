package transfer

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/money"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

type TransferRepo interface {
	Insert(ctx context.Context, t Transfer) errs.AppError
	Get(ctx context.Context, id string) (*Transfer, errs.AppError)
	List(ctx context.Context) ([]Transfer, errs.AppError)
}

type Transfer struct {
	ID             string `bson:"_id"`
	Player         player.Player
	TeamDestiny    team.Team
	DateOfTransfer string
	Amount         money.Money
	Created        time.Time
}

func (t Transfer) GetID() string {
	return t.ID
}

func (t *Transfer) SetID(id string) {
	t.ID = id
}
