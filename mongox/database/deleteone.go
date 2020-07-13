package database

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
	"github.com/mainnika/mongox-go-driver/v2/mongox/utils"
)

// DeleteOne removes a document from a database and then returns it into target
func (d *Database) DeleteOne(target interface{}, filters ...interface{}) (err error) {

	collection := d.GetCollectionOf(target)
	opts := &options.FindOneAndDeleteOptions{}
	composed := query.Compose(filters...)
	protected := base.GetProtection(target)

	opts.Sort = composed.Sorter()

	if !utils.IsNil(target) {
		composed.And(primitive.M{"_id": base.GetID(target)})
	}

	if protected != nil {
		query.Push(composed, protected)
		protected.X = primitive.NewObjectID()
		protected.V = time.Now().Unix()
	}

	result := collection.FindOneAndDelete(d.Context(), composed.M(), opts)
	if result.Err() != nil {
		return fmt.Errorf("can't create find one and delete result: %w", result.Err())
	}

	err = result.Decode(target)
	if err == mongox.ErrNoDocuments {
		return err
	}
	if err != nil {
		return fmt.Errorf("can't decode result: %w", err)
	}

	return
}
