package main

import (
	"context"

	"github.com/joho/godotenv"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/config"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/kafka"
	"github.com/rafaelsanzio/go-flashscore/pkg/model"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		_ = errs.ErrGettingEnv.Throwf(applog.Log, errs.ErrFmt, err)
	}

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

	conTransfer := kafka.Consumer{
		Broker:      []string{broker1Address, broker2Address, broker3Address},
		GroupTopics: model.Topics,
		GroupID:     "flashscore",
	}

	conTransfer.Init(ctx)
}
