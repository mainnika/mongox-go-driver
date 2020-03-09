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
func (p *Primary) GetID() primitive.ObjectID {
	return p.ID
}

// SetID sets an _id
func (p *Primary) SetID(id primitive.ObjectID) {
	p.ID = id
}
