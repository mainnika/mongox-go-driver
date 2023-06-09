package base_test

import (
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/docbased"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/oidbased"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/stringbased"
)

type DocWithCustomInterface struct {
	ID int `bson:"_id" json:"_id" collection:"4"`
}

func (d *DocWithCustomInterface) GetID() interface{} {
	return d.ID
}

func (d *DocWithCustomInterface) SetID(interface{}) {
	panic("not implemented")
}

func TestGetID(t *testing.T) {

	type DocWithObjectID struct {
		oidbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}
	id, err := base.GetID(&DocWithObjectID{Primary: oidbased.Primary{ID: [12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}}})
	require.NoError(t, err)
	assert.Equal(t, primitive.ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}), id)

	type DocWithObject struct {
		docbased.Primary `bson:",inline" json:",inline" collection:"2"`
	}
	id, err = base.GetID(&DocWithObject{Primary: docbased.Primary{ID: primitive.D{{"1", "2"}}}})
	require.NoError(t, err)
	assert.Equal(t, primitive.D{{"1", "2"}}, id)

	type DocWithString struct {
		stringbased.Primary `bson:",inline" json:",inline" collection:"3"`
	}
	id, err = base.GetID(&DocWithString{Primary: stringbased.Primary{ID: "foobar"}})
	assert.Equal(t, "foobar", id)

	id, err = base.GetID(&DocWithCustomInterface{ID: 420})
	assert.NoError(t, err)
	assert.Equal(t, 420, id)
}
