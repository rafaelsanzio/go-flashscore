package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/prototype"
	"github.com/rafaelsanzio/go-flashscore/pkg/transfer"
)

func TestTransferRepoInsert(t *testing.T) {
	ctx := context.Background()

	SetTransferRepo(MockTransferRepo{
		InsertFunc: func(ctx context.Context, p transfer.Transfer) errs.AppError {
			return nil
		},
	})
	defer SetTransferRepo(nil)

	newTransfer := prototype.PrototypeTransfer()

	err := GetTransferRepo().Insert(ctx, newTransfer)
	assert.NoError(t, err)
}

func TestTransferRepoGet(t *testing.T) {
	ctx := context.Background()

	SetTransferRepo(MockTransferRepo{
		GetFunc: func(ctx context.Context, id string) (*transfer.Transfer, errs.AppError) {
			player := prototype.PrototypeTransfer()
			return &player, nil
		},
	})
	defer SetTransferRepo(nil)

	newTransfer := prototype.PrototypeTransfer()

	result, err := GetTransferRepo().Get(ctx, "new-transfer-id")
	assert.NoError(t, err)

	assert.Equal(t, newTransfer, *result)
}

func TestTransferRepoList(t *testing.T) {
	ctx := context.Background()

	SetTransferRepo(MockTransferRepo{
		ListFunc: func(ctx context.Context) ([]transfer.Transfer, errs.AppError) {
			transferMock := prototype.PrototypeTransfer()
			transferMock2 := prototype.PrototypeTransfer()

			return []transfer.Transfer{transferMock, transferMock2}, nil
		},
	})
	defer SetTransferRepo(nil)

	players, err := GetTransferRepo().List(ctx)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(players))
}
