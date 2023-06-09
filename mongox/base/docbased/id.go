package docbased

import (
	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

var _ mongox.DocBased = (*Primary)(nil)

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

// New creates a new Primary structure with a defined _id
func New(e primitive.E, ee ...primitive.E) Primary {
	id := primitive.D{e}
	if len(ee) > 0 {
		id = append(id, ee...)
	}

	return Primary{ID: id}
}

func GetID(source mongox.DocBased) (id primitive.D, err error) {
	id = source.GetID()
	if !reflect2.IsNil(id) {
		return id, nil
	}

	return nil, mongox.ErrUninitializedBase
}
