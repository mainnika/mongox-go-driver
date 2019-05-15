package common

import (
	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LoadOne function loads a first single target document by a query
func LoadOne(db *mongox.Database, target interface{}, filters ...interface{}) error {

	collection := db.GetCollectionOf(target)
	opts := options.FindOne()
	composed := query.Compose(filters...)

	opts.Sort = composed.Sorter()

	result := collection.FindOne(db.Context(), composed.M(), opts)
	if result.Err() != nil {
		return errors.InternalErrorf("can't create find one result: %s", result.Err())
	}

	err := result.Decode(target)
	if err == mongo.ErrNoDocuments {
		return errors.NotFoundErrorf("%s", err)
	}
	if err != nil {
		return errors.InternalErrorf("can't decode desult: %s", err)
	}

	return nil
}
