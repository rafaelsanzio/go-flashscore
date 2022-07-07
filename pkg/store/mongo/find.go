package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/cursor"
	"github.com/rafaelsanzio/go-flashscore/pkg/store/query"
)

func (s *Store) FindOne(ctx context.Context, collection string, filter query.Filter, v interface{}, opts ...query.FindOneOptions) errs.AppError {
	col := s.client.Database(dbName).Collection(collection)

	f, err := bson.Marshal(filter)
	if err != nil {
		return errs.ErrMarshalingBson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	var mongoOpts []*options.FindOneOptions

	if opts != nil {
		mongoOpts = make([]*options.FindOneOptions, len(opts))

		for i, o := range opts {
			mongoOpts[i] = &options.FindOneOptions{Sort: o.Sort}
		}
	}

	err = col.FindOne(ctx, f, mongoOpts...).Decode(v)

	if err == mongo.ErrNoDocuments {
		return nil
	}

	if err != nil {
		return errs.ErrMongoFindOne.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return nil
}

func (s *Store) Find(ctx context.Context, collection string, filter query.Filter, opts ...query.FindOptions) (cursor.Cursor, errs.AppError) {
	col := s.client.Database(dbName).Collection(collection)

	var c cursor.Cursor
	var err error

	var mongoOpts []*options.FindOptions

	if opts != nil {
		mongoOpts = make([]*options.FindOptions, len(opts))

		for i, o := range opts {
			mongoOpts[i] = &options.FindOptions{Sort: o.Sort}
		}
	}

	if filter != nil {
		var f []byte
		f, err = bson.Marshal(filter)
		if err != nil {
			return nil, errs.ErrMarshalingBson.Throwf(applog.Log, errs.ErrFmt, err)
		}

		c, err = col.Find(ctx, f, mongoOpts...)
	} else {
		c, err = col.Find(ctx, bson.M{}, mongoOpts...)
	}

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, errs.ErrMongoFind.Throwf(applog.Log, "for collection: %s, "+errs.ErrFmt, collection, err)
	}

	return c, nil
}
