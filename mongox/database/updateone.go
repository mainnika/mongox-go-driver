package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// UpdateOne updates a single document in the database and loads it into target
func (d *Database) UpdateOne(target interface{}, filters ...interface{}) (err error) {

	collection := d.GetCollectionOf(target)
	opts := options.FindOneAndUpdate()
	protected := base.GetProtection(target)
	composed := query.Compose(filters...)
	ctx := query.WithContext(d.Context(), composed)
	updater := composed.Updater()

	opts.SetReturnDocument(options.After)

	if protected != nil {
		if !protected.X.IsZero() {
			query.Push(composed, protected)
		}
		updater = append(updater, primitive.M{
			"$set": primitive.M{
				"_x": primitive.NewObjectID(),
				"_v": time.Now().Unix(),
			},
		})
	}

	defer func() {
		invokerr := composed.OnClose().Invoke(ctx, target)
		if err == nil {
			err = invokerr
		}

		return
	}()

	result := collection.FindOneAndUpdate(ctx, composed.M(), updater, opts)
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
