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

	composed, err := query.Compose(filters...)
	if err != nil {
		return
	}

	zeroElem := reflect.Zero(targetSliceElemT)
	_, hasPreloader := composed.Preloader()
	ctx := query.WithContext(d.Context(), composed)

	var result *mongox.Cursor
	var i int

	defer func() {

		if result != nil {
			closerr := result.Close(ctx)
			if err == nil {
				err = closerr
			}
		}

		invokerr := composed.OnClose().Invoke(ctx, target)
		if err == nil {
			err = invokerr
		}

		return
	}()

	if hasPreloader {
		result, err = d.createAggregateLoad(zeroElem.Interface(), composed)
	} else {
		result, err = d.createSimpleLoad(zeroElem.Interface(), composed)
	}
	if err != nil {
		err = fmt.Errorf("can't create find result: %w", err)
		return
	}

	for i = 0; result.Next(ctx); i++ {

		var elem interface{}

		if i == targetSliceV.Len() {
			value := reflect.New(targetSliceElemT.Elem())
			elem = value.Interface()

			err = composed.OnCreate().Invoke(ctx, elem)
			if err != nil {
				return
			}

			err = result.Decode(elem)
			if err != nil {
				return
			}

			targetSliceV = reflect.Append(targetSliceV, value)
		} else {
			elem = targetSliceV.Index(i).Interface()

			if created := base.Reset(elem); created {
				err = composed.OnCreate().Invoke(ctx, elem)
			}
			if err != nil {
				return
			}

			err = result.Decode(elem)
			if err != nil {
				return
			}
		}

		err = composed.OnDecode().Invoke(ctx, elem)
		if err != nil {
			return
		}
	}

	targetSliceV = targetSliceV.Slice(0, i)
	targetV.Elem().Set(targetSliceV)

	return
}
