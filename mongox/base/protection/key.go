package protection

import (
	"time"

	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Key field stores unique document id and version
type Key struct {
	X primitive.ObjectID `bson:"_x" json:"_x" index:",hashed"`
	V int64              `bson:"_v" json:"_v"`
}

// PutToDocument extends the doc with protection key values
func (k *Key) PutToDocument(doc primitive.M) {

	if reflect2.IsNil(doc) {
		return
	}

	if k.X.IsZero() {
		doc["_x"] = primitive.M{"$exists": false}
		doc["_v"] = primitive.M{"$exists": false}
	} else {
		doc["_x"] = k.X
		doc["_v"] = k.V
	}
}

// Restate creates a new protection key
func (k *Key) Restate() {
	k.X = primitive.NewObjectID()
	k.V = time.Now().Unix()
}
