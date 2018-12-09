package common

import (
	"reflect"

	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// LoadStream function loads documents one by one into a target channel
func LoadStream(db *mongox.Database, target interface{}, composed *query.Query) error {

	targetV := reflect.ValueOf(target)
	targetT := targetV.Type()

	targetK := targetV.Kind()
	if targetK != reflect.Chan {
		panic(errors.InternalErrorf("target is not a chan"))
	}
	if targetT.Elem().Kind() != reflect.Ptr {
		panic(errors.InternalErrorf("chan element should be a document ptr"))
	}

	dummy := reflect.Zero(targetT.Elem())
	collection := db.GetCollectionOf(dummy.Interface())
	opts := &options.FindOptions{}

	if composed.Sorter() != nil {
		opts.Sort = composed.Sorter().Sort()
	}
	if composed.Limiter() != nil {
		limit := int64(composed.Limiter().Limit())
		opts.Limit = &limit
	}

	result, err := collection.Find(db.Context(), composed.M(), opts)
	if err != nil {
		return errors.InternalErrorf("can't create find result: %s", err)
	}

	go func() {
		defer result.Close(db.Context())

		for {
			elem, ok := targetV.Recv()
			if !ok {
				break
			}

			if result.Next(db.Context()) != true {
				targetV.Send(dummy)
				break
			}

			if result.Decode(elem.Interface()) != nil {
				targetV.Send(dummy)
				break
			}

			targetV.Send(elem)
		}
	}()

	return nil
}
