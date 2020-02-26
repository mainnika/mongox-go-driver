package base

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/mongox"
)

var _ mongox.BaseObject = &Object{}

// Object is a structure with object as an _id field
type Object struct {
	ID primitive.D `bson:"_id,omitempty" json:"_id,omitempty"`
}

// GetID returns an _id
func (db *Object) GetID() primitive.D {
	return db.ID
}

// SetID sets an _id
func (db *Object) SetID(id primitive.D) {
	db.ID = id
}
