package common

import (
	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/base"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// SaveOne saves a single source document to the database
func SaveOne(db *mongox.Database, source interface{}) error {

	collection := db.GetCollectionOf(source)
	opts := options.FindOneAndReplace()
	id := base.GetID(source)

	opts.SetUpsert(true)
	opts.SetReturnDocument(options.After)

	result := collection.FindOneAndReplace(db.Context(), bson.M{"_id": id}, source, opts)
	if result.Err() != nil {
		return errors.NotFoundErrorf("%s", result.Err())
	}

	return result.Decode(source)
}
