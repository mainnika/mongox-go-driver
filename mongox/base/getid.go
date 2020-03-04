package base

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

// GetID returns source document id
func GetID(source interface{}) (id interface{}) {

	switch doc := source.(type) {
	case mongox.ObjectIDBased:
		return getObjectIDOrGenerate(doc)
	case mongox.StringBased:
		return getStringIDOrPanic(doc)
	case mongox.ObjectBased:
		return getObjectOrPanic(doc)
	default:
		panic(fmt.Errorf("source contains malformed document, %v", source))
	}
}

func getObjectIDOrGenerate(source mongox.ObjectIDBased) (id primitive.ObjectID) {

	id = source.GetID()
	if id != primitive.NilObjectID {
		return id
	}

	id = primitive.NewObjectID()
	source.SetID(id)

	return
}

func getStringIDOrPanic(source mongox.StringBased) (id string) {

	id = source.GetID()
	if id != "" {
		return id
	}

	panic(fmt.Errorf("source contains malformed document, %v", source))
}

func getObjectOrPanic(source mongox.ObjectBased) (id primitive.D) {

	id = source.GetID()
	if id != nil {
		return id
	}

	panic(fmt.Errorf("source contains malformed document, %v", source))
}
