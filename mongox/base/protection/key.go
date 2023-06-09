package protection

import (
	"reflect"
	"time"

	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Key field stores unique document id and version
type Key struct {
	X primitive.ObjectID `bson:"_x" json:"_x" index:",hashed"`
	V int64              `bson:"_v" json:"_v"`
}

// Inject extends the doc with protection key values
func (k *Key) Inject(doc primitive.M) {
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

// Get finds protection field in the source document otherwise returns nil
func Get(source interface{}) (key *Key) {
	v := reflect.ValueOf(source)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil
	}

	el := v.Elem()
	numField := el.NumField()

	for i := 0; i < numField; i++ {
		field := el.Field(i)
		if !field.CanInterface() {
			continue
		}

		switch field.Interface().(type) {
		case *Key:
			key = field.Interface().(*Key)
		case Key:
			ptr := field.Addr()
			key = ptr.Interface().(*Key)
		default:
			continue
		}

		return key
	}

	return nil
}
