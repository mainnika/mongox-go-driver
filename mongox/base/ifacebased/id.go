package ifacebased

import (
	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/modern-go/reflect2"
)

// GetID returns an _id from the source document
func GetID(source mongox.InterfaceBased) (id interface{}, err error) {
	id = source.GetID()
	if !reflect2.IsNil(id) {
		return id, nil
	}

	return nil, mongox.ErrUninitializedBase
}
