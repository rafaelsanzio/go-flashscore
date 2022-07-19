package repo

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/player"
	"github.com/rafaelsanzio/go-flashscore/pkg/store"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
)

const (
	PlayerCollection = "player"
)

type playerRepo struct {
	store store.Store
}

var playerRepoSingleton player.PlayerRepo

func GetPlayerRepo() player.PlayerRepo {
	if playerRepoSingleton == nil {
		return getPlayerRepo()
	}
	return playerRepoSingleton
}

func getPlayerRepo() *playerRepo {
	s := store.GetStore()
	return &playerRepo{s}
}

func SetPlayerRepo(repo player.PlayerRepo) {
	playerRepoSingleton = repo
}

func (repo playerRepo) Insert(ctx context.Context, p player.Player) errs.AppError {
	p.Created = time.Now()
	_, err := repo.store.InsertOne(ctx, PlayerCollection, &p)
	return err
}

func (repo playerRepo) Get(ctx context.Context, id string) (*player.Player, errs.AppError) {
	filter := query.Filter{
		"_id": id,
	}

	opts := query.FindOneOptions{}

	mPlayer := player.Player{}
	err := repo.store.FindOne(ctx, PlayerCollection, filter, &mPlayer, opts)
	if err != nil {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", PlayerCollection, id, err)
	}

	if mPlayer.ID == "" {
		return nil, nil
	}

	return &mPlayer, nil
}

func (repo playerRepo) List(ctx context.Context) ([]player.Player, errs.AppError) {
	filter := query.Filter{}

	opts := query.FindOptions{}
	mPlayer := []player.Player{}
	players, err := repo.store.Find(ctx, PlayerCollection, filter, opts)
	if err != nil {
		return mPlayer, errs.ErrMongoFind.Throwf(applog.Log, "for collection: %s, err: [%v]", PlayerCollection, err)
	}

	defer func() {
		_ = players.Close(ctx)
	}()

	for {
		if players.Err() != nil {
			return mPlayer, err
		}

		if ok := players.Next(ctx); !ok {
			break
		}

		var p player.Player
		if err_ := players.Decode(&p); err_ != nil {
			return mPlayer, err
		}

		mPlayer = append(mPlayer, p)
	}

	return mPlayer, nil
}

func (repo playerRepo) Update(ctx context.Context, p player.Player) (*player.Player, errs.AppError) {
	res := player.Player{}
	filter := query.Filter{
		"_id": p.GetID(),
	}

	err := repo.store.FindOne(ctx, PlayerCollection, filter, &res)
	if err != nil {
		return nil, err
	}

	if res.ID == "" {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", PlayerCollection, p.GetID(), err)
	}

	p.ID = res.ID
	err = repo.store.UpdateOne(ctx, PlayerCollection, &p)
	if err != nil {
		return nil, errs.ErrMongoUpdateOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", PlayerCollection, p.GetID(), err)
	}

	return &p, nil
}

func (repo playerRepo) Delete(ctx context.Context, id string) errs.AppError {
	err := repo.store.DeleteOne(ctx, PlayerCollection, id)
	return err
}

func (repo playerRepo) GetTeamPlayer(ctx context.Context, id, teamID string) (*player.Player, errs.AppError) {
	filter := query.Filter{
		"_id":      id,
		"team._id": teamID,
	}

	opts := query.FindOneOptions{}

	mPlayer := player.Player{}
	err := repo.store.FindOne(ctx, PlayerCollection, filter, &mPlayer, opts)
	if err != nil {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", PlayerCollection, id, err)
	}

	if mPlayer.ID == "" {
		return nil, nil
	}

	return &mPlayer, nil
}
