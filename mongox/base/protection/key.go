package protection

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Key field stores unique document id and version
type Key struct {
	X primitive.ObjectID `bson:"_x" json:"_x" index:",hashed"`
	V int64              `bson:"_v" json:"_v"`
}

// Restate creates a new protection key
func (k *Key) Restate() {
	k.X = primitive.NewObjectID()
	k.V = time.Now().Unix()
}
