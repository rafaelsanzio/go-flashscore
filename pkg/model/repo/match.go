package repo

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/match"
	"github.com/rafaelsanzio/go-flashscore/pkg/store"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
)

const (
	MatchCollection = "match"
)

type matchRepo struct {
	store store.Store
}

var matchRepoSingleton match.MatchRepo

func GetMatchRepo() match.MatchRepo {
	if matchRepoSingleton == nil {
		return getMatchRepo()
	}
	return matchRepoSingleton
}

func getMatchRepo() *matchRepo {
	s := store.GetStore()
	return &matchRepo{s}
}

func SetMatchRepo(repo match.MatchRepo) {
	matchRepoSingleton = repo
}

func (repo matchRepo) Insert(ctx context.Context, mt match.Match) errs.AppError {
	mt.Created = time.Now()
	_, err := repo.store.InsertOne(ctx, MatchCollection, &mt)
	return err
}

func (repo matchRepo) Get(ctx context.Context, id string) (*match.Match, errs.AppError) {
	filter := query.Filter{
		"_id": id,
	}

	opts := query.FindOneOptions{}

	mMtach := match.Match{}
	err := repo.store.FindOne(ctx, MatchCollection, filter, &mMtach, opts)
	if err != nil {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", MatchCollection, id, err)
	}

	if mMtach.ID == "" {
		return nil, nil
	}

	return &mMtach, nil
}

func (repo matchRepo) List(ctx context.Context) ([]match.Match, errs.AppError) {
	filter := query.Filter{}

	opts := query.FindOptions{}
	mMtach := []match.Match{}
	matches, err := repo.store.Find(ctx, MatchCollection, filter, opts)
	if err != nil {
		return mMtach, errs.ErrMongoFind.Throwf(applog.Log, "for collection: %s, err: [%v]", MatchCollection, err)
	}

	defer func() {
		_ = matches.Close(ctx)
	}()

	for {
		if matches.Err() != nil {
			return mMtach, err
		}

		if ok := matches.Next(ctx); !ok {
			break
		}

		var p match.Match
		if err_ := matches.Decode(&p); err_ != nil {
			return mMtach, err
		}

		mMtach = append(mMtach, p)
	}

	return mMtach, nil
}

func (repo matchRepo) Update(ctx context.Context, p match.Match) (*match.Match, errs.AppError) {
	res := match.Match{}
	filter := query.Filter{
		"_id": p.GetID(),
	}

	err := repo.store.FindOne(ctx, MatchCollection, filter, &res)
	if err != nil {
		return nil, err
	}

	if res.ID == "" {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", MatchCollection, p.GetID(), err)
	}

	p.ID = res.ID
	err = repo.store.UpdateOne(ctx, MatchCollection, &p)
	if err != nil {
		return nil, errs.ErrMongoUpdateOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", MatchCollection, p.GetID(), err)
	}

	return &p, nil
}

func (repo matchRepo) Delete(ctx context.Context, id string) errs.AppError {
	err := repo.store.DeleteOne(ctx, MatchCollection, id)
	return err
}

func (repo matchRepo) FindMatchForTournament(ctx context.Context, id, tournamentID string) (*match.Match, errs.AppError) {
	filter := query.Filter{
		"_id":            id,
		"tournament._id": tournamentID,
	}

	opts := query.FindOneOptions{}

	mMtach := match.Match{}
	err := repo.store.FindOne(ctx, MatchCollection, filter, &mMtach, opts)
	if err != nil {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, and tournamentid: %s, err: [%v]", MatchCollection, id, tournamentID, err)
	}

	if mMtach.ID == "" {
		return nil, nil
	}

	return &mMtach, nil
}
