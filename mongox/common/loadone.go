package common

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/errors"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// LoadOne function loads a first single target document by a query
func LoadOne(db mongox.Database, target interface{}, filters ...interface{}) error {

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
		return errors.InternalErrorf("can't create find result: %s", err)
	}

	hasNext := result.Next(db.Context())
	if !hasNext {
		return errors.NotFoundErrorf("can't find result: %s", result.Err())
	}

	base.Reset(target)

	return result.Decode(target)
}
