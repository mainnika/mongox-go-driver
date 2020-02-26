package common

import (
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/base"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
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

	composed := query.Compose(filters...)
	zeroElem := reflect.Zero(targetSliceElemT)
	hasPreloader, _ := composed.Preloader()

	var result *mongo.Cursor
	var err error

	if hasPreloader {
		result, err = createAggregateLoad(db, zeroElem.Interface(), composed)
	} else {
		result, err = createSimpleLoad(db, zeroElem.Interface(), composed)
	}
	if err != nil {
		return errors.InternalErrorf("can't create find result: %s", err)
	}

	defer result.Close(db.Context())
	var i int

	for i = 0; result.Next(db.Context()); {
		if targetSliceV.Len() == i {
			elem := reflect.New(targetSliceElemT.Elem())
			if err = result.Decode(elem.Interface()); err == nil {
				targetSliceV = reflect.Append(targetSliceV, elem)
			} else {
				continue
			}
		} else {
			elem := targetSliceV.Index(i).Interface()
			base.Reset(elem)
			if err = result.Decode(elem); err != nil {
				continue
			}
		}

		i++
	}

	targetSliceV = targetSliceV.Slice(0, i)
	targetV.Elem().Set(targetSliceV)

	return nil
}
