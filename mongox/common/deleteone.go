package common

import (
	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// DeleteOne removes a document from a database and then returns it into target
func DeleteOne(db *mongox.Database, target interface{}, filters ...interface{}) error {

	collection := db.GetCollectionOf(target)
	opts := &options.FindOneAndDeleteOptions{}
	composed := query.Compose(filters...)

	opts.Sort = composed.Sorter()

	result := collection.FindOneAndDelete(db.Context(), composed.M(), opts)
	if result.Err() != nil {
		return errors.InternalErrorf("can't create find one and delete result: %s", result.Err())
	}

	err := result.Decode(target)
	if err == mongo.ErrNoDocuments {
		return errors.NotFoundErrorf("%s", err)
	}
	if err != nil {
		return errors.InternalErrorf("can't decode result: %s", err)
	}

	return nil
}
