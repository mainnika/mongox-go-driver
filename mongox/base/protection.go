package base

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Protection field stores unique document id and version
type Protection struct {
	X primitive.ObjectID `bson:"_x" json:"_x" index:",hashed"`
	V int64              `bson:"_v" json:"_v"`
}
