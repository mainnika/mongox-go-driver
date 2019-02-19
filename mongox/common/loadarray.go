package common

import (
	"reflect"

	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// LoadArray loads an array of documents from the database by query
func LoadArray(db *mongox.Database, target interface{}, filters ...interface{}) error {

	targetV := reflect.ValueOf(target)
	targetT := targetV.Type()

	targetK := targetV.Kind()
	if targetK != reflect.Ptr {
		panic(errors.InternalErrorf("target is not a ptr"))
	}

	targetSliceV := targetV.Elem()
	targetSliceT := targetT.Elem()
	if targetSliceT.Kind() != reflect.Slice {
		panic(errors.InternalErrorf("target should be a ptr to a slice"))
	}

	targetSliceElemT := targetSliceT.Elem()
	if targetSliceElemT.Kind() != reflect.Ptr {
		panic(errors.InternalErrorf("target slice should contain ptrs"))
	}

	dummy := reflect.Zero(targetSliceElemT)
	collection := db.GetCollectionOf(dummy.Interface())
	opts := options.Find()
	composed := query.Compose(filters...)

	opts.Sort = composed.Sorter()
	opts.Limit = composed.Limiter()
	opts.Skip = composed.Skipper()

	result, err := collection.Find(db.Context(), composed.M(), opts)
	if err != nil {
		return errors.InternalErrorf("can't create find result: %s", err)
	}

	defer result.Close(db.Context())
	var i int

	for i = 0; result.Next(db.Context()); {
		if targetSliceV.Len() == i {
			elem := reflect.New(targetSliceElemT.Elem())
			if result.Decode(elem.Interface()) != nil {
				targetSliceV = reflect.Append(targetSliceV, elem)
			} else {
				continue
			}
		} else {
			result.Decode(targetSliceV.Index(i).Interface())
		}

		i++
	}

	targetSliceV = targetSliceV.Slice(0, i)
	targetV.Elem().Set(targetSliceV)

	return nil
}
