package common

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/base"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
)

// DeleteArray removes documents list from a database by their ids
func DeleteArray(db *mongox.Database, target interface{}) error {

	targetV := reflect.ValueOf(target)
	targetT := targetV.Type()

	targetK := targetV.Kind()
	if targetK != reflect.Ptr {
		panic(errors.Malformedf("target is not a ptr"))
	}

	targetSliceV := targetV.Elem()
	targetSliceT := targetT.Elem()
	if targetSliceT.Kind() != reflect.Slice {
		panic(errors.Malformedf("target should be a ptr to a slice"))
	}

	targetSliceElemT := targetSliceT.Elem()
	if targetSliceElemT.Kind() != reflect.Ptr {
		panic(errors.Malformedf("target slice should contain ptrs"))
	}

	zeroElem := reflect.Zero(targetSliceElemT)
	targetLen := targetSliceV.Len()
	collection := db.GetCollectionOf(zeroElem.Interface())
	opts := options.Delete()
	ids := primitive.A{}

	for i := 0; i < targetLen; i++ {
		elem := targetSliceV.Index(i)
		ids = append(ids, base.GetID(elem.Interface()))
	}

	if len(ids) == 0 {
		return errors.Malformedf("can't delete zero elements")
	}

	result, err := collection.DeleteMany(db.Context(), primitive.M{"_id": primitive.M{"$in": ids}}, opts)
	if err != nil {
		return errors.NotFoundErrorf("can't create find and delete result: %s", err)
	}
	if result.DeletedCount != int64(targetLen) {
		return errors.InternalErrorf("can't verify delete result: removed count mismatch %d != %d", result.DeletedCount, targetLen)
	}

	return nil
}
