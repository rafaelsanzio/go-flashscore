package mongo

import (
	"context"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/config"
	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

const connectTimeout = 5000 // ms

var dbName string

type Store struct {
	client *mongo.Client
}

func NewStore() (*Store, errs.AppError) {
	v, err := config.Value(key.MongoDBName)
	if err != nil {
		return nil, err.Annotatef(applog.Log, "unable to get mongo config: %v", err)
	}

	dbName = v
	client, err := Client()
	if err != nil {
		return nil, err
	}

	return &Store{
		client: client,
	}, nil
}

func Client() (*mongo.Client, errs.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Millisecond)
	defer cancel()

	opts := options.Client()

	uri, err := config.Value(key.MongoURI)
	if err != nil {
		return nil, err
	}

	user, err := config.Value(key.MongoDBUsername)
	if err != nil {
		return nil, err
	}

	pass, err := config.Value(key.MongoDBPassword)
	if err != nil {
		return nil, err
	}

	parsedUri, err_ := url.ParseRequestURI(uri)

	if err_ != nil {
		return nil, errs.ErrParseRequestURI.Throwf(applog.Log, errs.ErrFmt, err_)
	}

	if user != "" {
		parsedUri.User = url.UserPassword(user, pass)
	}

	opts.ApplyURI(parsedUri.String())

	cl, err_ := mongo.Connect(ctx, opts)

	if err_ != nil {
		return cl, errs.ErrMongoConnect.Throwf(applog.Log, errs.ErrFmt, err_)
	}

	return cl, nil
}
