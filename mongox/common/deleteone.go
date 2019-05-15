package common

import (
	"time"

	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/base"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeleteOne removes a document from a database and then returns it into target
func DeleteOne(db *mongox.Database, target interface{}, filters ...interface{}) error {

	collection := db.GetCollectionOf(target)
	opts := &options.FindOneAndDeleteOptions{}
	composed := query.Compose(filters...)
	protected := base.GetProtection(target)

	opts.Sort = composed.Sorter()

	if target != nil {
		composed.And(primitive.M{"_id": base.GetID(target)})
	}

	if protected != nil {
		query.Push(composed, protected)
		protected.X = primitive.NewObjectID()
		protected.V = time.Now().Unix()
	}

	result := collection.FindOneAndDelete(db.Context(), composed.M(), opts)
	if result.Err() != nil {
		return errors.InternalErrorf("can't create find one and delete result: %s", result.Err())
	}

	err := result.Decode(target)
	if err == mongo.ErrNoDocuments {
		return errors.NotFoundErrorf("%s", err)
	}
	if err != nil {
		return errors.InternalErrorf("can't decode result: %s", err)
	}

	return nil
}
