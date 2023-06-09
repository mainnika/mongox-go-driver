package oidbased

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

var _ mongox.OIDBased = (*Primary)(nil)

// Primary is a structure with objectId as the primary key
type Primary struct {
	ID primitive.ObjectID `bson:"_id" json:"_id"`
}

// GetID returns an _id
func (p *Primary) GetID() (id primitive.ObjectID) {
	return p.ID
}

// SetID sets an _id
func (p *Primary) SetID(id primitive.ObjectID) {
	p.ID = id
}

// Generate creates a new Primary structure with a new objectId
func Generate() Primary {
	return Primary{ID: primitive.NewObjectID()}
}

// New creates a new Primary structure with a defined objectId
func New(id primitive.ObjectID) Primary {
	return Primary{ID: id}
}

func GetID(source mongox.OIDBased) (id primitive.ObjectID, err error) {
	id = source.GetID()
	if id != primitive.NilObjectID {
		return id, nil
	}

	return primitive.NilObjectID, mongox.ErrUninitializedBase
}
