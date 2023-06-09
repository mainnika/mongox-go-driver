package stringbased_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mainnika/mongox-go-driver/v2/mongox-testing/database"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/stringbased"
)

func Test_GetID(t *testing.T) {
	type DocWithString struct {
		stringbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	doc := &DocWithString{Primary: stringbased.New("foobar")}

	assert.Equal(t, "foobar", doc.GetID())
}

func Test_SetID(t *testing.T) {
	type DocWithString struct {
		stringbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	doc := &DocWithString{Primary: stringbased.New("foobar")}

	doc.SetID("rockrockrock")

	assert.Equal(t, "rockrockrock", doc.Primary.ID)
	assert.Equal(t, "rockrockrock", doc.GetID())
}

func Test_SaveLoad(t *testing.T) {
	type DocWithObjectID struct {
		stringbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	db, err := database.NewEphemeral("")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = db.Close() }()

	doc1 := &DocWithObjectID{Primary: stringbased.New("foobar")}
	doc2 := &DocWithObjectID{}

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
		stringbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	doc := &DocWithObjectID{Primary: stringbased.New("foobar")}

	bytes, err := json.Marshal(doc)
	assert.NoError(t, err)
	assert.Equal(t, `{"_id":"foobar"}`, string(bytes))
}
