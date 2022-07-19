package repo

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/transfer"
)

type MockTransferRepo struct {
	transfer.TransferRepo
	InsertFunc func(ctx context.Context, t transfer.Transfer) errs.AppError
	GetFunc    func(ctx context.Context, id string) (*transfer.Transfer, errs.AppError)
	ListFunc   func(ctx context.Context) ([]transfer.Transfer, errs.AppError)
}

func (m MockTransferRepo) Insert(ctx context.Context, p transfer.Transfer) errs.AppError {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, p)
	}
	return m.TransferRepo.Insert(ctx, p)
}

func (m MockTransferRepo) Get(ctx context.Context, id string) (*transfer.Transfer, errs.AppError) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return m.TransferRepo.Get(ctx, id)
}

func (m MockTransferRepo) List(ctx context.Context) ([]transfer.Transfer, errs.AppError) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return m.TransferRepo.List(ctx)
}
