package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

func (s *Store) InsertOne(ctx context.Context, collection string, data interface{}) (string, errs.AppError) {
	col := s.client.Database(dbName).Collection(collection)
	doc, ok := data.(Document)

	if !ok {
		return "", errs.ErrNotDocumentInterface.Throw(applog.Log)
	}

	id := doc.GetID()

	if id == "" {
		doc.SetID(primitive.NewObjectID().Hex())
	}

	b, err := bson.Marshal(doc)
	if err != nil {
		return "", errs.ErrMarshalingBson.Throwf(applog.Log, errs.ErrFmt, err)
	}

	res, err := col.InsertOne(ctx, b)
	if err != nil {
		return "", errs.ErrMongoInsertOne.Throwf(applog.Log, errs.ErrFmt, err)
	}

	if res == nil {
		return "", nil
	}

	id, ok = res.InsertedID.(string)

	if !ok {
		return "", errs.ErrDecodingInsertedId.Throw(applog.Log)
	}

	return id, nil
}
