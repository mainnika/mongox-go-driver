package oidbased_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox-testing/database"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/oidbased"
)

func Test_GetID(t *testing.T) {
	type DocWithObjectID struct {
		oidbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	doc := &DocWithObjectID{Primary: oidbased.New([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2})}

	assert.Equal(t, primitive.ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}), doc.Primary.ID)
	assert.Equal(t, primitive.ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}), doc.GetID())
}

func Test_SetID(t *testing.T) {
	type DocWithObjectID struct {
		oidbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	doc := &DocWithObjectID{Primary: oidbased.Generate()}
	doc.SetID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2})

	assert.Equal(t, primitive.ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}), doc.Primary.ID)
	assert.Equal(t, primitive.ObjectID([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}), doc.GetID())
}

func Test_SaveLoad(t *testing.T) {
	type DocWithObjectID struct {
		oidbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	db, err := database.NewEphemeral("")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = db.Close() }()

	doc1 := &DocWithObjectID{Primary: oidbased.Generate()}
	doc2 := &DocWithObjectID{Primary: oidbased.Generate()}

	err = db.SaveOne(doc1)
	assert.NoError(t, err)

	err = db.LoadOne(doc2)
	assert.NoError(t, err)

	assert.Equal(t, doc1, doc2)

	bytes1, _ := json.Marshal(doc1)
	bytes2, _ := json.Marshal(doc2)

	assert.Equal(t, bytes1, bytes2)
}

func Test_Marshal(t *testing.T) {
	type DocWithObjectID struct {
		oidbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	id, _ := primitive.ObjectIDFromHex("feadbeeffeadbeeffeadbeef")
	doc := &DocWithObjectID{Primary: oidbased.New(id)}

	bytes, err := json.Marshal(doc)
	assert.NoError(t, err)
	assert.Equal(t, `{"_id":"feadbeeffeadbeeffeadbeef"}`, string(bytes))
}
