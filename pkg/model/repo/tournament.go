package repo

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/store"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
	"github.com/rafaelsanzio/go-flashscore/pkg/tournament"
)

const (
	TournamentCollection = "tournament"
)

type tournamentRepo struct {
	store store.Store
}

var tournamentRepoSingleton tournament.TournamentRepo

func GetTournamentRepo() tournament.TournamentRepo {
	if tournamentRepoSingleton == nil {
		return getTournamentRepo()
	}
	return tournamentRepoSingleton
}

func getTournamentRepo() *tournamentRepo {
	s := store.GetStore()
	return &tournamentRepo{s}
}

func SetTournamentRepo(repo tournament.TournamentRepo) {
	tournamentRepoSingleton = repo
}

func (repo tournamentRepo) Insert(ctx context.Context, t tournament.Tournament) errs.AppError {
	t.Created = time.Now()
	_, err := repo.store.InsertOne(ctx, TournamentCollection, &t)
	return err
}

func (repo tournamentRepo) Get(ctx context.Context, id string) (*tournament.Tournament, errs.AppError) {
	filter := query.Filter{
		"_id": id,
	}

	opts := query.FindOneOptions{}

	mTournament := tournament.Tournament{}
	err := repo.store.FindOne(ctx, TournamentCollection, filter, &mTournament, opts)
	if err != nil {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", TournamentCollection, id, err)
	}

	if mTournament.ID == "" {
		return nil, nil
	}

	return &mTournament, nil
}

func (repo tournamentRepo) List(ctx context.Context) ([]tournament.Tournament, errs.AppError) {
	filter := query.Filter{}

	opts := query.FindOptions{}
	mTournament := []tournament.Tournament{}
	tournaments, err := repo.store.Find(ctx, TournamentCollection, filter, opts)
	if err != nil {
		return mTournament, errs.ErrMongoFind.Throwf(applog.Log, "for collection: %s, err: [%v]", TournamentCollection, err)
	}

	defer func() {
		_ = tournaments.Close(ctx)
	}()

	for {
		if tournaments.Err() != nil {
			return mTournament, err
		}

		if ok := tournaments.Next(ctx); !ok {
			break
		}

		var t tournament.Tournament
		if err_ := tournaments.Decode(&t); err_ != nil {
			return mTournament, err
		}

		mTournament = append(mTournament, t)
	}

	return mTournament, nil
}

func (repo tournamentRepo) Update(ctx context.Context, t tournament.Tournament) (*tournament.Tournament, errs.AppError) {
	res := tournament.Tournament{}
	filter := query.Filter{
		"_id": t.GetID(),
	}

	err := repo.store.FindOne(ctx, TournamentCollection, filter, &res)
	if err != nil {
		return nil, err
	}

	if res.ID == "" {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", TournamentCollection, t.GetID(), err)
	}

	t.ID = res.ID
	t.Created = res.Created
	err = repo.store.UpdateOne(ctx, TournamentCollection, &t)
	if err != nil {
		return nil, errs.ErrMongoUpdateOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", TournamentCollection, t.GetID(), err)
	}

	return &t, nil
}

func (repo tournamentRepo) Delete(ctx context.Context, id string) errs.AppError {
	err := repo.store.DeleteOne(ctx, TournamentCollection, id)
	return err
}
