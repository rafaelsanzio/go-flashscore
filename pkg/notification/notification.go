package notification

import (
	"context"
	"encoding/json"

	"github.com/rafaelsanzio/go-flashscore/cmd/kafka/handlers"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
)

type Notification struct {
	Context context.Context `json:"-"`
	Action  model.ActionType
	Data    map[string]string
}

func Handler(ctx context.Context, body, messageId string) errs.AppError {
	var pn Notification
	err := json.Unmarshal([]byte(body), &pn)
	if err != nil {
		return errs.ErrUnmarshalingJson.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	switch pn.Action {
	case model.ActionUpdateTeamPlayer:
		err := handlers.HandleEventUpdateTeamPlayer(ctx, pn.Data)
		if err != nil {
			return errs.ErrHandlingUpdateTeamPlayer.Throwf(applog.Log, errs.ErrFmt, err.Error())
		}
	case model.ActionGameEvents:
		switch model.EventsMatchType(pn.Data["matchEventType"]) {
		case model.EventStart:
			err := handlers.HandleEventMatchStart(ctx, pn.Data)
			if err != nil {
				return errs.ErrHandlingGameEventStarted.Throwf(applog.Log, errs.ErrFmt, err.Error())
			}
		case model.EventGoal:
			err := handlers.HandleEventMatchGoal(ctx, pn.Data)
			if err != nil {
				return errs.ErrHandlingGameEventGoal.Throwf(applog.Log, errs.ErrFmt, err.Error())
			}
		case model.EventHalftime:
			err := handlers.HandleEventMatchHalfTime(ctx, pn.Data)
			if err != nil {
				return errs.ErrHandlingGameEventHalftime.Throwf(applog.Log, errs.ErrFmt, err.Error())
			}
		case model.EventExtratime:
			err := handlers.HandleEventMatchExtratime(ctx, pn.Data)
			if err != nil {
				return errs.ErrHandlingGameEventExtratime.Throwf(applog.Log, errs.ErrFmt, err.Error())
			}
		case model.EventSubstitution:
			err := handlers.HandleEventMatchSubstitution(ctx, pn.Data)
			if err != nil {
				return errs.ErrHandlingGameEventFinish.Throwf(applog.Log, errs.ErrFmt, err.Error())
			}
		case model.EventWarning:
			err := handlers.HandleEventMatchWarning(ctx, pn.Data)
			if err != nil {
				return errs.ErrHandlingGameEventWarning.Throwf(applog.Log, errs.ErrFmt, err.Error())
			}
		case model.EventFinish:
			err := handlers.HandleEventMatchFinish(ctx, pn.Data)
			if err != nil {
				return errs.ErrHandlingGameEventFinish.Throwf(applog.Log, errs.ErrFmt, err.Error())
			}
		default:
			return errs.ErrActionNotImplemented.Throwf(applog.Log, errs.ErrFmt, pn.Data["matchEventType"])
		}
	default:
		return errs.ErrActionNotImplemented.Throwf(applog.Log, errs.ErrFmt, pn.Action)
	}

	return nil
}
