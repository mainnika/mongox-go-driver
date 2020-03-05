package common

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// LoadArray loads an array of documents from the database by query
func LoadArray(db mongox.Database, target interface{}, filters ...interface{}) error {

	targetV := reflect.ValueOf(target)
	targetT := targetV.Type()

	targetK := targetV.Kind()
	if targetK != reflect.Ptr {
		panic(fmt.Errorf("target is not a ptr"))
	}

	targetSliceV := targetV.Elem()
	targetSliceT := targetT.Elem()
	if targetSliceT.Kind() != reflect.Slice {
		panic(fmt.Errorf("target should be a ptr to a slice"))
	}

	targetSliceElemT := targetSliceT.Elem()
	if targetSliceElemT.Kind() != reflect.Ptr {
		panic(fmt.Errorf("target slice should contain ptrs"))
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
		return fmt.Errorf("can't create find result: %w", err)
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
