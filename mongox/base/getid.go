package base

import (
	"fmt"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/docbased"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/ifacebased"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/oidbased"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/stringbased"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

// GetID returns source document id
func GetID(source interface{}) (id interface{}, err error) {
	switch doc := source.(type) {
	case mongox.OIDBased:
		return oidbased.GetID(doc)
	case mongox.StringBased:
		return stringbased.GetID(doc)
	case mongox.DocBased:
		return docbased.GetID(doc)
	case mongox.InterfaceBased:
		return ifacebased.GetID(doc)
	default:
		return nil, fmt.Errorf("%w: unknown base type", mongox.ErrMalformedBase)
	}
}
