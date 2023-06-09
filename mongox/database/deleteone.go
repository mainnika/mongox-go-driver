package database

import (
	"fmt"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/protection"

	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// DeleteOne removes a document from a database and then returns it into target
func (d *Database) DeleteOne(target interface{}, filters ...interface{}) (err error) {
	composed, err := query.Compose(filters...)
	if err != nil {
		return err
	}

	collection, err := d.GetCollectionOf(target)
	if err != nil {
		return err
	}

	if !reflect2.IsNil(target) {
		targetID, err := base.GetID(target)
		if err != nil {
			return err
		}

		composed.And(primitive.M{"_id": targetID})
	}

	protected := protection.Get(target)
	if protected != nil {
		_, err := query.Push(composed, protected)
		if err != nil {
			return err
		}

		protected.Restate()
	}

	ctx := query.WithContext(d.Context(), composed)
	m := composed.M()
	opts := options.FindOneAndDelete()
	opts.Sort = composed.Sorter()

	defer func() { _ = composed.OnClose().Invoke(ctx, target) }()

	result := collection.FindOneAndDelete(ctx, m, opts)
	if result.Err() != nil {
		return fmt.Errorf("can't create find one and delete result: %w", result.Err())
	}

	err = result.Decode(target)
	if err != nil {
		return fmt.Errorf("can't decode find one and delete result: %w", err)
	}

	_ = composed.OnDecode().Invoke(ctx, target)

	return nil
}
