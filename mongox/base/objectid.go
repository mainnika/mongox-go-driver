package base

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// ObjectID is a structure with objectId as an _id field
type ObjectID struct {
	ID primitive.ObjectID `bson:"_id" json:"_id"`
}

// GetID returns an _id
func (db *ObjectID) GetID() primitive.ObjectID {
	return db.ID
}

// SetID sets an _id
func (db *ObjectID) SetID(id primitive.ObjectID) {
	db.ID = id
}
