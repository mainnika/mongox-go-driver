package database

import (
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/protection"
	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// UpdateOne updates a single document in the database and loads it into target
func (d *Database) UpdateOne(target interface{}, filters ...interface{}) (err error) {
	composed, err := query.Compose(filters...)
	if err != nil {
		return err
	}

	update, err := composed.Updater()
	if err != nil {
		return err
	}

	protected := protection.Get(target)
	if protected != nil {
		if !protected.X.IsZero() {
			query.Push(composed, protected)
		}
		protected.Restate()

		setCmd, _ := update["$set"].(primitive.M)
		if reflect2.IsNil(setCmd) {
			setCmd = primitive.M{}
		}
		protected.Inject(setCmd)
		update["$set"] = setCmd
	}

	collection, err := d.GetCollectionOf(target)
	if err != nil {
		return err
	}

	ctx := query.WithContext(d.Context(), composed)
	m := composed.M()
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)

	defer func() { _ = composed.OnClose().Invoke(ctx, target) }()

	result := collection.FindOneAndUpdate(ctx, m, update, opts)
	if result.Err() != nil {
		return result.Err()
	}

	err = result.Decode(target)
	if err != nil {
		return err
	}

	_ = composed.OnDecode().Invoke(ctx, target)

	return nil
}
