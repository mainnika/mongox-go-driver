package database

import (
	"fmt"
	"reflect"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// LoadArray loads an array of documents from the database by query
func (d *Database) LoadArray(target interface{}, filters ...interface{}) (err error) {

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

	var result *mongox.Cursor
	var i int

	if hasPreloader {
		result, err = d.createAggregateLoad(zeroElem.Interface(), composed)
	} else {
		result, err = d.createSimpleLoad(zeroElem.Interface(), composed)
	}
	if err != nil {
		err = fmt.Errorf("can't create find result: %w", err)
		return
	}

	for i = 0; result.Next(d.Context()); {

		var elem interface{}

		if targetSliceV.Len() == i {
			value := reflect.New(targetSliceElemT.Elem())
			err = result.Decode(value.Interface())
			elem = value.Interface()
			if err == nil {
				targetSliceV = reflect.Append(targetSliceV, value)
			}
		} else {
			elem = targetSliceV.Index(i).Interface()
			base.Reset(elem)
			err = result.Decode(elem)
		}
		if err != nil {
			_ = result.Close(d.Context())
			return
		}

		err = composed.OnDecode().Invoke(d.Context(), elem)
		if err != nil {
			_ = result.Close(d.Context())
			return
		}

		i++
	}

	targetSliceV = targetSliceV.Slice(0, i)
	targetV.Elem().Set(targetSliceV)

	return result.Close(d.Context())
}
