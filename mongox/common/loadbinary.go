package common

import (
	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoadBinary function loads a first single target document by a query
func LoadBinary(db *mongox.Database, target interface{}, filters ...interface{}) (bson.Raw, error) {

	composed := query.Compose(append(filters, query.Limit(1))...)
	hasPreloader, _ := composed.Preloader()

	var result *mongo.Cursor
	var err error

	if hasPreloader {
		result, err = createAggregateLoad(db, target, composed)
	} else {
		result, err = createSimpleLoad(db, target, composed)
	}
	if err != nil {
		return nil, errors.InternalErrorf("can't create find result: %s", err)
	}

	hasNext := result.Next(db.Context())
	if !hasNext {
		return nil, errors.NotFoundErrorf("can't find result: %s", result.Err())
	}

	return result.Current, nil
}
