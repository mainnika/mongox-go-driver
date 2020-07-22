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
	collection := d.GetCollectionOf(zeroElem.Interface())
	opts := options.Delete()
	ids := primitive.A{}
	composed := query.Compose(filters...)

	for i := 0; i < targetLen; i++ {
		elem := targetSliceV.Index(i)
		ids = append(ids, base.GetID(elem.Interface()))
	}

	if len(ids) == 0 {
		return fmt.Errorf("can't delete zero elements")
	}

	composed.And(primitive.M{"_id": primitive.M{"$in": ids}})

	result, err := collection.DeleteMany(d.Context(), composed.M(), opts)
	if err != nil {
		return
	}
	if result.DeletedCount != int64(targetLen) {
		err = fmt.Errorf("can't verify delete result: removed count mismatch %d != %d", result.DeletedCount, targetLen)
	}

	return
}
