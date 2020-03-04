package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestObject_GetID(t *testing.T) {

	type DocWithObject struct {
		Object `bson:"_id" json:"_id" collection:"2"`
	}

	doc := &DocWithObject{Object: Object(primitive.D{{"1", "one"}, {"2", "two"}})}

	assert.Equal(t, primitive.D{{"1", "one"}, {"2", "two"}}, doc.GetID())
}

func TestObject_SetID(t *testing.T) {

	type DocWithObject struct {
		Object `bson:"_id" json:"_id" collection:"2"`
	}

	doc := &DocWithObject{Object: Object(primitive.D{{"1", "one"}, {"2", "two"}})}

	doc.SetID(primitive.D{{"3", "three"}, {"4", "you"}})

	assert.Equal(t, primitive.D{{"3", "three"}, {"4", "you"}}, primitive.D(doc.Object))
	assert.Equal(t, primitive.D{{"3", "three"}, {"4", "you"}}, primitive.D(doc.GetID()))
}

func TestString_GetID(t *testing.T) {

	type DocWithString struct {
		String `bson:"_id" json:"_id" collection:"3"`
	}

	doc := &DocWithString{String: String("foobar")}

	assert.Equal(t, "foobar", doc.GetID())
}

func TestString_SetID(t *testing.T) {

	type DocWithString struct {
		String `bson:"_id" json:"_id" collection:"3"`
	}

	doc := &DocWithString{String: String("foobar")}

	doc.SetID("rockrockrock")

	assert.Equal(t, "rockrockrock", string(doc.String))
	assert.Equal(t, "rockrockrock", doc.GetID())
}

func TestObjectID_GetID(t *testing.T) {

	type DocWithObjectID struct {
		ObjectID `bson:"_id" json:"_id" collection:"1"`
	}

	doc := &DocWithObjectID{ObjectID: [12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}}

	assert.Equal(t, primitive.ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}), doc.GetID())
}

func TestObjectID_SetID(t *testing.T) {

	type DocWithObjectID struct {
		ObjectID `bson:"_id" json:"_id" collection:"1"`
	}

	doc := &DocWithObjectID{}

	doc.SetID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2})

	assert.Equal(t, primitive.ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}), primitive.ObjectID(doc.ObjectID))
	assert.Equal(t, primitive.ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}), doc.GetID())
}
