package common

import (
	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// LoadOne function loads a first single target document by a query
func LoadOne(db *mongox.Database, target interface{}, composed *query.Query) error {

	collection := db.GetCollectionOf(target)
	opts := options.FindOne()

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
