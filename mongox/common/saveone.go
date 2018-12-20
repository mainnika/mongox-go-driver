package common

import (
	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// SaveOne saves a single source document to the database
func SaveOne(db *mongox.Database, source interface{}) error {

	collection := db.GetCollectionOf(source)
	opts := &options.FindOneAndReplaceOptions{}

	opts.SetUpsert(true)
	opts.SetReturnDocument(options.After)

	var id interface{}

	switch doc := source.(type) {
	case mongox.BaseObjectID:
		id = doc.GetID()
		if id == primitive.NilObjectID {
			id = primitive.NewObjectID()
		}
	case mongox.BaseString:
		id = doc.GetID()
		if id == "" {
			panic(errors.Malformedf("source contains malformed document, %v", source))
		}
	default:
		panic(errors.Malformedf("source contains malformed document, %v", source))
	}

	result := collection.FindOneAndReplace(db.Context(), bson.M{"_id": id}, source, opts)
	if result.Err() != nil {
		return errors.NotFoundErrorf("%s", result.Err())
	}

	return result.Decode(source)
}
