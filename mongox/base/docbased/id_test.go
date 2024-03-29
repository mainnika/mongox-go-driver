package docbased_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox-testing/database"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/docbased"
)

func Test_GetID(t *testing.T) {
	type DocWithObject struct {
		docbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	doc := &DocWithObject{Primary: docbased.New(primitive.E{"1", "one"}, primitive.E{"2", "two"})}

	assert.Equal(t, primitive.D{{"1", "one"}, {"2", "two"}}, doc.GetID())
}

func Test_SetID(t *testing.T) {
	type DocWithObject struct {
		docbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	doc := &DocWithObject{Primary: docbased.New(primitive.E{"1", "one"}, primitive.E{"2", "two"})}

	doc.SetID(primitive.D{{"3", "three"}, {"4", "you"}})

	assert.Equal(t, primitive.D{{"3", "three"}, {"4", "you"}}, doc.Primary.ID)
	assert.Equal(t, primitive.D{{"3", "three"}, {"4", "you"}}, doc.GetID())
}

func Test_SaveLoad(t *testing.T) {
	type DocWithObjectID struct {
		docbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	db, err := database.NewEphemeral("")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = db.Close() }()

	doc1 := &DocWithObjectID{Primary: docbased.New(primitive.E{"1", "one"}, primitive.E{"2", "two"})}
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
		docbased.Primary `bson:",inline" json:",inline" collection:"1"`
	}

	doc := &DocWithObjectID{Primary: docbased.New(primitive.E{"1", "one"}, primitive.E{"2", "two"})}

	bytes, err := json.Marshal(doc)
	assert.NoError(t, err)
	assert.Equal(t, `{"_id":[{"Key":"1","Value":"one"},{"Key":"2","Value":"two"}]}`, string(bytes))
}
