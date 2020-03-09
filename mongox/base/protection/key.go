package protection

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Key field stores unique document id and version
type Key struct {
	X primitive.ObjectID `bson:"_x" json:"_x" index:",hashed"`
	V int64              `bson:"_v" json:"_v"`
}
