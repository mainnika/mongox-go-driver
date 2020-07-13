package jsonbased

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

var _ mongox.JSONBased = (*Primary)(nil)

// Primary is a structure with object as an _id field
type Primary struct {
	ID primitive.D `bson:"_id" json:"_id"`
}

// GetID returns an _id
func (p *Primary) GetID() (id primitive.D) {
	return p.ID
}

// SetID sets an _id
func (p *Primary) SetID(id primitive.D) {
	p.ID = id
}
