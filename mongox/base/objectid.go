package base

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

var _ mongox.BaseObjectID = &ObjectID{}

// ObjectID is a structure with objectId as an _id field
type ObjectID struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
}

// GetID returns an _id
func (db *ObjectID) GetID() primitive.ObjectID {
	return db.ID
}

// SetID sets an _id
func (db *ObjectID) SetID(id primitive.ObjectID) {
	db.ID = id
}
