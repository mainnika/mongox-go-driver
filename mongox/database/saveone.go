package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// SaveOne saves a single source document to the database
func (d *Database) SaveOne(source interface{}, filters ...interface{}) (err error) {

	collection := d.GetCollectionOf(source)
	opts := options.FindOneAndReplace()
	id := base.GetID(source)
	protected := base.GetProtection(source)
	composed := query.Compose(filters...)
	ctx := query.WithContext(d.Context(), composed)

	composed.And(primitive.M{"_id": id})

	opts.SetUpsert(true)
	opts.SetReturnDocument(options.After)

	if protected != nil {
		query.Push(composed, protected)
		protected.X = primitive.NewObjectID()
		protected.V = time.Now().Unix()
	}

	defer func() {
		invokerr := composed.OnClose().Invoke(ctx, source)
		if err == nil {
			err = invokerr
		}

		return
	}()

	result := collection.FindOneAndReplace(ctx, composed.M(), source, opts)
	if result.Err() != nil {
		return result.Err()
	}

	err = result.Decode(source)
	if err != nil {
		return
	}

	err = composed.OnDecode().Invoke(ctx, source)
	if err != nil {
		return
	}

	return
}
