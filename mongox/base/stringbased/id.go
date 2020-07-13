package stringbased

import (
	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

var _ mongox.StringBased = (*Primary)(nil)

// Primary is a structure with string as an _id field
type Primary struct {
	ID string `bson:"_id" json:"_id"`
}

// GetID returns an _id
func (p *Primary) GetID() (id string) {
	return p.ID
}

// SetID sets an _id
func (p *Primary) SetID(id string) {
	p.ID = id
}
