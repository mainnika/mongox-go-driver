package base

import (
	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetID returns source document id
func GetID(source interface{}) (id interface{}) {

	switch doc := source.(type) {
	case mongox.BaseObjectID:
		return getObjectIdOrGenerate(doc)
	case mongox.BaseString:
		return getStringIdOrPanic(doc)
	case mongox.BaseObject:
		return getObjectOrPanic(doc)
	default:
		panic(errors.Malformedf("source contains malformed document, %v", source))
	}

	return
}

func getObjectIdOrGenerate(source mongox.BaseObjectID) (id primitive.ObjectID) {

	id = source.GetID()
	if id != primitive.NilObjectID {
		return id
	}

	id = primitive.NewObjectID()
	source.SetID(id)

	return
}

func getStringIdOrPanic(source mongox.BaseString) (id string) {

	id = source.GetID()
	if id != "" {
		return id
	}

	panic(errors.Malformedf("victim contains malformed document, %v", source))
}

func getObjectOrPanic(source mongox.BaseObject) (id primitive.D) {

	id = source.GetID()
	if id != nil {
		return id
	}

	panic(errors.Malformedf("victim contains malformed document, %v", source))
}
