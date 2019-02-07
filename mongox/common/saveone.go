package common

import (
	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/base"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"
)

// SaveOne saves a single source document to the database
func SaveOne(db *mongox.Database, source interface{}) error {

	collection := db.GetCollectionOf(source)
	opts := options.FindOneAndReplace()
	id := base.GetID(source)
	protected := base.GetProtection(source)
	composed := query.Compose(bson.M{"_id": id})

	opts.SetUpsert(true)
	opts.SetReturnDocument(options.After)

	if protected != nil {
		query.Push(composed, protected)
		protected.X = primitive.NewObjectID()
		protected.V = time.Now().Unix()
	}

	result := collection.FindOneAndReplace(db.Context(), composed.M(), source, opts)
	if result.Err() != nil {
		return errors.NotFoundErrorf("%s", result.Err())
	}

	return result.Decode(source)
}
