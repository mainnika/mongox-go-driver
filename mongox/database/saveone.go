package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// SaveOne saves a single source document to the database
func (d *Database) SaveOne(source interface{}, filters ...interface{}) (err error) {

	composed, err := query.Compose(filters...)
	if err != nil {
		return
	}

	collection := d.GetCollectionOf(source)
	id := base.GetID(source)
	protected := base.GetProtection(source)
	ctx := query.WithContext(d.Context(), composed)

	composed.And(primitive.M{"_id": id})

	opts := options.FindOneAndReplace()
	opts.SetUpsert(true)
	opts.SetReturnDocument(options.After)

	if protected != nil {
		query.Push(composed, protected)
		protected.Restate()
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
