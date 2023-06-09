package database

import (
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/protection"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// SaveOne saves a single source document to the database
func (d *Database) SaveOne(source interface{}, filters ...interface{}) (err error) {
	collection, err := d.GetCollectionOf(source)
	if err != nil {
		return err
	}

	composed, err := query.Compose(filters...)
	if err != nil {
		return err
	}

	id, err := base.GetID(source)
	if err != nil {
		return err
	}
	composed.And(primitive.M{"_id": id})

	protected := protection.Get(source)
	if protected != nil {
		query.Push(composed, protected)
		protected.Restate()
	}

	ctx := query.WithContext(d.Context(), composed)
	m := composed.M()
	opts := options.FindOneAndReplace()
	opts.SetUpsert(true)
	opts.SetReturnDocument(options.After)

	defer func() { _ = composed.OnClose().Invoke(ctx, source) }()

	result := collection.FindOneAndReplace(ctx, m, source, opts)
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
