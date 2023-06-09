package database

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// DeleteArray removes documents list from a database by their ids
func (d *Database) DeleteArray(target interface{}, filters ...interface{}) (err error) {
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
	targetLen := targetSliceV.Len()

	collection, err := d.GetCollectionOf(zeroElem.Interface())
	if err != nil {
		return err
	}

	composed, err := query.Compose(filters...)
	if err != nil {
		return err
	}

	if targetLen > 0 {
		var ids primitive.A
		for i := 0; i < targetLen; i++ {
			elem := targetSliceV.Index(i)
			elemID, err := base.GetID(elem.Interface())
			if err != nil {
				return err
			}

			ids = append(ids, elemID)
		}
		composed.And(primitive.M{"_id": primitive.M{"$in": ids}})
	}

	ctx := query.WithContext(d.Context(), composed)
	m := composed.M()
	opts := options.Delete()

	defer func() { _ = composed.OnClose().Invoke(ctx, target) }()

	result, err := collection.DeleteMany(ctx, m, opts)
	if err != nil {
		return fmt.Errorf("while deleting array: %w", err)
	}
	if result.DeletedCount != int64(targetLen) {
		return fmt.Errorf("deleted count mismatch %d != %d", result.DeletedCount, targetLen)
	}

	return nil
}
