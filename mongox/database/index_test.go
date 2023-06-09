package database_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mainnika/mongox-go-driver/v2/mongox-testing/database"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base/oidbased"
)

func TestDatabase_Ensure(t *testing.T) {

	db, err := database.NewEphemeral("")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testvalues := []struct {
		doc      interface{}
		settings map[string]interface{}
		index    map[string]interface{}
	}{
		{
			doc: &struct {
				oidbased.Primary `bson:",inline" json:",inline" collection:"1"`

				Foobar int `bson:"foobar" json:"foobar" index:"-,unique,allowNull,expireAfter=86400"`
				Foo    int `bson:"foo" json:"foo"`
				Bar    int `bson:"bar" json:"bar"`
			}{},
			index: map[string]interface{}{
				"background":         false,
				"expireAfterSeconds": int32(86400),
				"key": map[string]interface{}{
					"foobar": int32(-1),
				},
				"name": "-,unique,allowNull,expireAfter=86400_foobar_",
				"partialFilterExpression": map[string]interface{}{
					"foobar": map[string]interface{}{"$exists": true},
				},
				"unique": true,
			},
		},
		{
			doc: &struct {
				oidbased.Primary `bson:",inline" json:",inline" collection:"2"`

				Foobar int `bson:"foobar" json:"foobar" index:",unique"`
			}{},
			index: map[string]interface{}{
				"background": false,
				"key": map[string]interface{}{
					"foobar": int32(1),
				},
				"name":   ",unique_foobar_",
				"unique": true,
			},
		},
		{
			doc: &struct {
				oidbased.Primary `bson:",inline" json:",inline" collection:"3"`

				Foobar int `bson:"foobar" json:"foobar" index:"-,+foo,+-bar,unique,allowNull"`
				Foo    int `bson:"foo" json:"foo"`
				Bar    int `bson:"bar" json:"bar"`
			}{},
			index: map[string]interface{}{
				"background": false,
				"key": map[string]interface{}{
					"foobar": int32(-1),
					"foo":    int32(1),
					"bar":    int32(-1),
				},
				"name": "-,+foo,+-bar,unique,allowNull_foobar_",
				"partialFilterExpression": map[string]interface{}{
					"foobar": map[string]interface{}{"$exists": true},
				},
				"unique": true,
			},
		},
		{
			doc: &struct {
				oidbased.Primary `bson:",inline" json:",inline" collection:"4"`

				Foobar int `bson:"foobar" json:"foobar" index:""`
				Foo    int `bson:"foo" json:"foo"`
				Bar    int `bson:"bar" json:"bar"`
			}{},
			index: map[string]interface{}{
				"background": false,
				"key": map[string]interface{}{
					"foobar": int32(1),
				},
				"name": "_foobar_",
			},
		},
		{
			doc: &struct {
				oidbased.Primary `bson:",inline" json:",inline" collection:"5"`

				Foobar int `bson:"foobar" json:"foobar" index:"-"`
				Foo    int `bson:"foo" json:"foo"`
				Bar    int `bson:"bar" json:"bar"`
			}{},
			index: map[string]interface{}{
				"background": false,
				"key": map[string]interface{}{
					"foobar": int32(-1),
				},
				"name": "-_foobar_",
			},
		},
		{
			doc: &struct {
				oidbased.Primary `bson:",inline" json:",inline" collection:"1"`

				Foobar int `bson:"foobar" json:"foobar" index:"-,unique,allowNull,expireAfter={{.Expire}}"`
				Foo    int `bson:"foo" json:"foo"`
				Bar    int `bson:"bar" json:"bar"`
			}{},
			settings: map[string]interface{}{
				"Expire": 86400,
			},
			index: map[string]interface{}{
				"background":         false,
				"expireAfterSeconds": int32(86400),
				"key": map[string]interface{}{
					"foobar": int32(-1),
				},
				"name": "-,unique,allowNull,expireAfter=86400_foobar_",
				"partialFilterExpression": map[string]interface{}{
					"foobar": map[string]interface{}{"$exists": true},
				},
				"unique": true,
			},
		},
	}

	for _, tt := range testvalues {
		err = db.IndexEnsure(tt.settings, tt.doc)
		assert.NoError(t, err)

		collection, err := db.GetCollectionOf(tt.doc)
		require.NoError(t, err)

		indexes, _ := collection.Indexes().List(db.Context())
		index := new(map[string]interface{})

		indexes.Next(db.Context()) // skip _id_
		indexes.Next(db.Context())
		indexes.Decode(&index)

		for k, v := range tt.index {
			assert.Equal(t, v, (*index)[k])
		}
	}
}
