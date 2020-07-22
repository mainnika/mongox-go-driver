package base

import (
	"fmt"

	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

// GetID returns source document id
func GetID(source interface{}) (id interface{}) {

	switch doc := source.(type) {
	case mongox.OIDBased:
		return getObjectIDOrGenerate(doc)
	case mongox.StringBased:
		return getStringIDOrPanic(doc)
	case mongox.JSONBased:
		return getObjectOrPanic(doc)
	case mongox.InterfaceBased:
		return getInterfaceOrPanic(doc)

	default:
		panic(fmt.Errorf("source contains malformed document, %v", source))
	}
}

func getObjectIDOrGenerate(source mongox.OIDBased) (id primitive.ObjectID) {

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

func getObjectOrPanic(source mongox.JSONBased) (id primitive.D) {

	id = source.GetID()
	if id != nil {
		return id
	}

	panic(fmt.Errorf("source contains malformed document, %v", source))
}

func getInterfaceOrPanic(source mongox.InterfaceBased) (id interface{}) {

	id = source.GetID()
	if !reflect2.IsNil(id) {
		return id
	}

	panic(fmt.Errorf("source contains malformed document, %v", source))
}
