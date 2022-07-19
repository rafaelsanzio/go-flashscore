package repo

import (
	"context"
	"time"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/store"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
	"github.com/rafaelsanzio/go-flashscore/pkg/team"
)

const (
	TeamCollection = "team"
)

type teamRepo struct {
	store store.Store
}

var teamRepoSingleton team.TeamRepo

func GetTeamRepo() team.TeamRepo {
	if teamRepoSingleton == nil {
		return getTeamRepo()
	}
	return teamRepoSingleton
}

func getTeamRepo() *teamRepo {
	s := store.GetStore()
	return &teamRepo{s}
}

func SetTeamRepo(repo team.TeamRepo) {
	teamRepoSingleton = repo
}

func (repo teamRepo) Insert(ctx context.Context, t team.Team) errs.AppError {
	t.Created = time.Now()
	_, err := repo.store.InsertOne(ctx, TeamCollection, &t)
	return err
}

func (repo teamRepo) Get(ctx context.Context, id string) (*team.Team, errs.AppError) {
	filter := query.Filter{
		"_id": id,
	}

	opts := query.FindOneOptions{}

	mTeam := team.Team{}
	err := repo.store.FindOne(ctx, TeamCollection, filter, &mTeam, opts)
	if err != nil {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", TeamCollection, id, err)
	}

	if mTeam.ID == "" {
		return nil, nil
	}

	return &mTeam, nil
}

func (repo teamRepo) List(ctx context.Context) ([]team.Team, errs.AppError) {
	filter := query.Filter{}

	opts := query.FindOptions{}
	mTeam := []team.Team{}
	teams, err := repo.store.Find(ctx, TeamCollection, filter, opts)
	if err != nil {
		return mTeam, errs.ErrMongoFind.Throwf(applog.Log, "for collection: %s, err: [%v]", TeamCollection, err)
	}

	defer func() {
		_ = teams.Close(ctx)
	}()

	for {
		if teams.Err() != nil {
			return mTeam, err
		}

		if ok := teams.Next(ctx); !ok {
			break
		}

		var u team.Team
		if err_ := teams.Decode(&u); err_ != nil {
			return mTeam, err
		}

		mTeam = append(mTeam, u)
	}

	return mTeam, nil
}

func (repo teamRepo) Update(ctx context.Context, t team.Team) (*team.Team, errs.AppError) {
	res := team.Team{}
	filter := query.Filter{
		"_id": t.GetID(),
	}

	err := repo.store.FindOne(ctx, TeamCollection, filter, &res)
	if err != nil {
		return nil, err
	}

	if res.ID == "" {
		return nil, errs.ErrMongoFindOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", TeamCollection, t.GetID(), err)
	}

	t.ID = res.ID
	t.Created = res.Created
	err = repo.store.UpdateOne(ctx, TeamCollection, &t)
	if err != nil {
		return nil, errs.ErrMongoUpdateOne.Throwf(applog.Log, "for collection: %s, and ID: %s, err: [%v]", TeamCollection, t.GetID(), err)
	}

	return &t, nil
}

func (repo teamRepo) Delete(ctx context.Context, id string) errs.AppError {
	err := repo.store.DeleteOne(ctx, TeamCollection, id)
	return err
}
