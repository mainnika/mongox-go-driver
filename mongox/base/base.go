package base

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

var _ mongox.ObjectBased = (*Object)(nil)
var _ mongox.ObjectIDBased = (*ObjectID)(nil)
var _ mongox.StringBased = (*String)(nil)

// Object is a structure with object as an _id field
type Object primitive.D

// GetID returns an _id
func (db *Object) GetID() primitive.D {
	return primitive.D(*db)
}

// SetID sets an _id
func (db *Object) SetID(id primitive.D) {
	*db = Object(id)
}

// ObjectID is a structure with objectId as an _id field
type ObjectID primitive.ObjectID

// GetID returns an _id
func (db *ObjectID) GetID() primitive.ObjectID {
	return primitive.ObjectID(*db)
}

// SetID sets an _id
func (db *ObjectID) SetID(id primitive.ObjectID) {
	*db = ObjectID(id)
}

// String is a structure with string as an _id field
type String string

// GetID returns an _id
func (db *String) GetID() string {
	return string(*db)
}

// SetID sets an _id
func (db *String) SetID(id string) {
	*db = String(id)
}
