package database

import (
	"fmt"
	"time"

	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// DeleteOne removes a document from a database and then returns it into target
func (d *Database) DeleteOne(target interface{}, filters ...interface{}) (err error) {

	collection := d.GetCollectionOf(target)
	opts := &options.FindOneAndDeleteOptions{}
	composed := query.Compose(filters...)
	protected := base.GetProtection(target)
	ctx := query.WithContext(d.Context(), composed)

	opts.Sort = composed.Sorter()

	if !reflect2.IsNil(target) {
		composed.And(primitive.M{"_id": base.GetID(target)})
	}

	if protected != nil {
		query.Push(composed, protected)
		protected.X = primitive.NewObjectID()
		protected.V = time.Now().Unix()
	}

	defer func() {
		invokerr := composed.OnClose().Invoke(ctx, target)
		if err == nil {
			err = invokerr
		}

		return
	}()

	result := collection.FindOneAndDelete(ctx, composed.M(), opts)
	if result.Err() != nil {
		return fmt.Errorf("can't create find one and delete result: %w", result.Err())
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
