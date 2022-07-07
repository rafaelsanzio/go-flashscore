package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

func (s *Store) UpdateOne(ctx context.Context, collection string, data interface{}) errs.AppError {
	col := s.client.Database(dbName).Collection(collection)
	doc, ok := data.(Document)

	if !ok {
		return errs.ErrNotDocumentInterface.Throw(applog.Log)
	}

	_, err := col.UpdateOne(ctx, bson.M{"_id": doc.GetID()}, bson.M{"$set": data})
	if err != nil {
		return errs.ErrMongoUpdateOne.Throwf(applog.Log, errs.ErrFmt, err)
	}

	return nil
}
