package database

import (
	"fmt"
	"reflect"

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

	zeroElem := reflect.Zero(targetSliceElemT)

	composed, err := query.Compose(filters...)
	if err != nil {
		return err
	}

	ctx := query.WithContext(d.Context(), composed)

	defer func() { _ = composed.OnClose().Invoke(ctx, target) }()

	cur, err := d.createCursor(zeroElem.Interface(), composed)
	if err != nil {
		return fmt.Errorf("can't create find result: %w", err)
	}

	defer func() { _ = cur.Close(ctx) }()

	var i int
	for i = 0; cur.Next(ctx); i++ {
		var elem interface{}
		if i == targetSliceV.Len() {
			value := reflect.New(targetSliceElemT.Elem())
			elem = value.Interface()

			_ = composed.OnCreate().Invoke(ctx, elem)

			err = cur.Decode(elem)
			if err != nil {
				return err
			}

			targetSliceV = reflect.Append(targetSliceV, value)
		} else {
			elem = targetSliceV.Index(i).Interface()

			if created := base.Reset(elem); created {
				_ = composed.OnCreate().Invoke(ctx, elem)
			}

			err = cur.Decode(elem)
			if err != nil {
				return err
			}
		}

		_ = composed.OnDecode().Invoke(ctx, elem)
	}
	err = cur.Err()
	if err != nil {
		return err
	}

	targetSliceV = targetSliceV.Slice(0, i)
	targetV.Elem().Set(targetSliceV)

	return nil
}
