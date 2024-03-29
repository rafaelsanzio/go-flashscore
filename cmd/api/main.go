package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"

	"github.com/rafaelsanzio/go-flashscore/pkg/api"
	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/cache"
	"github.com/rafaelsanzio/go-flashscore/pkg/config"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		_ = errs.ErrGettingEnv.Throwf(applog.Log, errs.ErrFmt, err)
	}

	store.GetStore() // mongo
	cache.GetStore() // redis

	mongoURI, err_ := config.Value(key.MongoURI)
	if err_ != nil {
		_ = err_.Annotatef(applog.Log, "unable to get mongo config: %v", err_)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		_ = errs.ErrMongoConnect.Throwf(applog.Log, errs.ErrFmt, err)
	}

	ctx, cancFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancFunc()

	err = client.Connect(ctx)
	if err != nil {
		_ = errs.ErrMongoConnect.Throwf(applog.Log, errs.ErrFmt, err)
	}

	defer func() {
		err = client.Disconnect(ctx)
		if err != nil {
			_ = errs.ErrMongoConnect.Throwf(applog.Log, errs.ErrFmt, err)
		}
	}()

	log.Println("MongoDB server is healthy.")

	appPort, err_ := config.Value(key.AppPort)
	if err_ != nil {
		_ = err_.Annotatef(applog.Log, "unable to get app port config: %v", err_)
	}

	log.Println("starting up... on PORT", appPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appPort), api.NewRouter()))
}
