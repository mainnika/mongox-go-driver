package base

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
		ObjectID `bson:"_id" json:"_id" collection:"1"`
	}
	type DocWithObject struct {
		Object `bson:"_id" json:"_id" collection:"2"`
	}
	type DocWithString struct {
		String `bson:"_id" json:"_id" collection:"3"`
	}

	GetID(&DocWithObjectID{ObjectID: ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2})})
	GetID(&DocWithObject{Object: Object(primitive.D{{"1", "2"}})})
	GetID(&DocWithString{String: String("foobar")})
	GetID(&DocWithCustomInterface{ID: 420})
}
