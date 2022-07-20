package kafka

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/config"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
	"github.com/rafaelsanzio/go-flashscore/pkg/notification"
)

func Notify(ctx context.Context, data map[string]string, action model.ActionType, keyKafka, topic string) {

	broker1Address, err_ := config.Value(key.KafkaAddress1)
	if err_ != nil {
		_ = errs.ErrGettingEnv.Throwf(applog.Log, errs.ErrFmt, err_)
	}

	broker2Address, err_ := config.Value(key.KafkaAddress3)
	if err_ != nil {
		_ = errs.ErrGettingEnv.Throwf(applog.Log, errs.ErrFmt, err_)
	}

	broker3Address, err_ := config.Value(key.KafkaAddress3)
	if err_ != nil {
		_ = errs.ErrGettingEnv.Throwf(applog.Log, errs.ErrFmt, err_)
	}

	bodyNotify := notification.Notification{
		Action: action,
		Data:   data,
	}

	pro := Producer{
		Broker: []string{broker1Address, broker2Address, broker3Address},
		Topic:  topic,
		Key:    keyKafka,
		Body:   bodyNotify,
	}

	pro.Create(context.Background())
}
