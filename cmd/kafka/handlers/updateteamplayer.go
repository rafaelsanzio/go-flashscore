package handlers

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model/repo"
)

func HandleEventUpdateTeamPlayer(ctx context.Context, data map[string]string) errs.AppError {
	playerID := data["playerID"]
	teamDestinyID := data["teamDestinyID"]

	player, err := repo.GetPlayerRepo().Get(ctx, playerID)
	if err != nil || player == nil {
		return errs.ErrPlayerIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	teamDestiny, err := repo.GetTeamRepo().Get(ctx, teamDestinyID)
	if err != nil || teamDestiny == nil {
		return errs.ErrTeamIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	player.Team = *teamDestiny
	playerUpdated, err := repo.GetPlayerRepo().Update(ctx, *player)
	if err != nil || playerUpdated == nil {
		return errs.ErrPlayerIsNotFound.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	return nil
}
