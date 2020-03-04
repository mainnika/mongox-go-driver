package base

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/errors"
)

// GetID returns source document id
func GetID(source interface{}) (id interface{}) {

	switch doc := source.(type) {
	case mongox.BaseObjectID:
		return getObjectIDOrGenerate(doc)
	case mongox.BaseString:
		return getStringIDOrPanic(doc)
	case mongox.BaseObject:
		return getObjectOrPanic(doc)
	default:
		panic(errors.Malformedf("source contains malformed document, %v", source))
	}
}

func getObjectIDOrGenerate(source mongox.BaseObjectID) (id primitive.ObjectID) {

	id = source.GetID()
	if id != primitive.NilObjectID {
		return id
	}

	id = primitive.NewObjectID()
	source.SetID(id)

	return
}

func getStringIDOrPanic(source mongox.BaseString) (id string) {

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
