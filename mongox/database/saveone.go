package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// SaveOne saves a single source document to the database
func (d *Database) SaveOne(source interface{}) (err error) {

	collection := d.GetCollectionOf(source)
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

	result := collection.FindOneAndReplace(d.Context(), composed.M(), source, opts)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(source)
}
