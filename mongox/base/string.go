package base

import (
	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

var _ mongox.BaseString = &String{}

// String is a structure with string as an _id field
type String struct {
	ID string `bson:"_id,omitempty" json:"_id,omitempty"`
}

// GetID returns an _id
func (db *String) GetID() string {
	return db.ID
}

// SetID sets an _id
func (db *String) SetID(id string) {
	db.ID = id
}
