package repo

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/store"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
	"github.com/rafaelsanzio/go-flashscore/pkg/transfer"
)

const (
	TransferCollection = "transfer"
)

type transferRepo struct {
	store store.Store
}

var transferRepoSingleton transfer.TransferRepo

func GetTransferRepo() transfer.TransferRepo {
	if transferRepoSingleton == nil {
		return getTransferRepo()
	}
	return transferRepoSingleton
}

func getTransferRepo() *transferRepo {
	s := store.GetStore()
	return &transferRepo{s}
}

func SetTransferRepo(repo transfer.TransferRepo) {
	transferRepoSingleton = repo
}

func (repo transferRepo) Insert(ctx context.Context, p transfer.Transfer) errs.AppError {
	p.Created = time.Now()
	_, err := repo.store.InsertOne(ctx, TransferCollection, &p)
	return err
}

func (repo transferRepo) Get(ctx context.Context, id string) (*transfer.Transfer, errs.AppError) {
	filter := query.Filter{
		"_id": id,
	}

	opts := query.FindOneOptions{}

	mTransfer := transfer.Transfer{}
	err := repo.store.FindOne(ctx, TransferCollection, filter, &mTransfer, opts)
	if err != nil {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", TransferCollection, id, err)
	}

	if mTransfer.ID == "" {
		return nil, nil
	}

	return &mTransfer, nil
}

func (repo transferRepo) List(ctx context.Context) ([]transfer.Transfer, errs.AppError) {
	filter := query.Filter{}

	opts := query.FindOptions{}
	mTransfer := []transfer.Transfer{}
	transfers, err := repo.store.Find(ctx, TransferCollection, filter, opts)
	if err != nil {
		return mTransfer, errs.ErrMongoFind.Throwf(applog.Log, "for collection: %s, err: [%v]", TransferCollection, err)
	}

	defer func() {
		_ = transfers.Close(ctx)
	}()

	for {
		if transfers.Err() != nil {
			return mTransfer, err
		}

		if ok := transfers.Next(ctx); !ok {
			break
		}

		var p transfer.Transfer
		if err_ := transfers.Decode(&p); err_ != nil {
			return mTransfer, err
		}

		mTransfer = append(mTransfer, p)
	}

	return mTransfer, nil
}
