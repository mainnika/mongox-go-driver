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

// New creates a new Primary structure with a defined _id
func New(id string) Primary {
	return Primary{ID: id}
}

func GetID(source mongox.StringBased) (id string, err error) {
	id = source.GetID()
	if id != "" {
		return id, nil
	}

	return "", mongox.ErrUninitializedBase
}
