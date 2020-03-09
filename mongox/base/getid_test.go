package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base/jsonbased"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/oidbased"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/stringbased"
)

type DocWithCustomInterface struct {
	ID int `bson:"_id" json:"_id" collection:"4"`
}

func (d *DocWithCustomInterface) GetID() interface{} {
	return d.ID
}

func (d *DocWithCustomInterface) SetID(id interface{}) {
	panic("not implemented")
}

func TestGetID(t *testing.T) {

	type DocWithObjectID struct {
		oidbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}
	type DocWithObject struct {
		jsonbased.Primary `bson:",inline" json:",inline" collection:"2"`
	}
	type DocWithString struct {
		stringbased.Primary `bson:",inline" json:",inline" collection:"3"`
	}

	assert.Equal(t, primitive.ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}), GetID(&DocWithObjectID{oidbased.Primary{[12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}}}))
	assert.Equal(t, primitive.D{{"1", "2"}}, GetID(&DocWithObject{jsonbased.Primary{primitive.D{{"1", "2"}}}}))
	assert.Equal(t, "foobar", GetID(&DocWithString{stringbased.Primary{"foobar"}}))
	assert.Equal(t, 420, GetID(&DocWithCustomInterface{ID: 420}))
}
