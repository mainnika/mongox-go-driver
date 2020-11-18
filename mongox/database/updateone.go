package database

import (
	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// UpdateOne updates a single document in the database and loads it into target
func (d *Database) UpdateOne(target interface{}, filters ...interface{}) (err error) {

	composed, err := query.Compose(filters...)
	if err != nil {
		return
	}

	updaterDoc, err := composed.Updater()
	if err != nil {
		return
	}

	collection := d.GetCollectionOf(target)
	protected := base.GetProtection(target)
	ctx := query.WithContext(d.Context(), composed)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)

	if protected != nil {
		if !protected.X.IsZero() {
			query.Push(composed, protected)
		}

		protected.Restate()

		setCmd, _ := updaterDoc["$set"].(primitive.M)
		if reflect2.IsNil(setCmd) {
			setCmd = primitive.M{}
		}
		protected.PutToDocument(setCmd)
		updaterDoc["$set"] = setCmd
	}

	defer func() {
		invokerr := composed.OnClose().Invoke(ctx, target)
		if err == nil {
			err = invokerr
		}

		return
	}()

	result := collection.FindOneAndUpdate(ctx, composed.M(), updaterDoc, opts)
	if result.Err() != nil {
		return result.Err()
	}

	err = result.Decode(target)
	if err != nil {
		return
	}

	err = composed.OnDecode().Invoke(ctx, target)
	if err != nil {
		return
	}

	return
}
