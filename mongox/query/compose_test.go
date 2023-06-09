package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base/protection"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

func TestPushBSON(t *testing.T) {

	q := &query.Query{}

	ok, err := query.Push(q, primitive.M{"foo": "bar"})

	assert.True(t, ok)
	assert.NoError(t, err)
	assert.NotEmpty(t, q.M())
	assert.Len(t, q.M()["$and"], 1)
	assert.Contains(t, q.M()["$and"], primitive.M{"foo": "bar"})

	ok, err = query.Push(q, primitive.M{"bar": "foo"})

	assert.True(t, ok)
	assert.NoError(t, err)
	assert.NotEmpty(t, q.M())
	assert.Len(t, q.M()["$and"], 2)
	assert.Contains(t, q.M()["$and"], primitive.M{"foo": "bar"})
	assert.Contains(t, q.M()["$and"], primitive.M{"bar": "foo"})
}

func TestPushLimiter(t *testing.T) {

	q := &query.Query{}
	lim := query.Limit(2)

	ok, err := query.Push(q, lim)

	assert.True(t, ok)
	assert.NoError(t, err)
	assert.NotNil(t, q.Limiter())
	assert.EqualValues(t, q.Limiter(), query.Limit(2).Limit())
}

func TestPushSorter(t *testing.T) {

	q := &query.Query{}
	sort := query.Sort{"foo": 1}

	ok, err := query.Push(q, sort)

	assert.True(t, ok)
	assert.NoError(t, err)
	assert.NotNil(t, q.Sorter())
	assert.EqualValues(t, q.Sorter(), primitive.M{"foo": 1})
}

func TestPushSkipper(t *testing.T) {

	q := &query.Query{}
	skip := query.Skip(66)

	ok, err := query.Push(q, skip)

	assert.True(t, ok)
	assert.NoError(t, err)
	assert.NotNil(t, q.Skipper())
	assert.EqualValues(t, q.Skipper(), query.Skip(66).Skip())
}

func TestPushProtection(t *testing.T) {

	t.Run("push protection key pointer", func(t *testing.T) {
		q := &query.Query{}
		protected := &protection.Key{V: 1, X: primitive.ObjectID{2}}

		ok, err := query.Push(q, protected)

		assert.True(t, ok)
		assert.NoError(t, err)
		assert.NotEmpty(t, q.M()["$and"])
		assert.Contains(t, q.M()["$and"], primitive.M{"_x": primitive.ObjectID{2}, "_v": int64(1)})
	})

	t.Run("push protection key struct", func(t *testing.T) {
		q := &query.Query{}
		protected := protection.Key{V: 1, X: primitive.ObjectID{2}}

		ok, err := query.Push(q, protected)

		assert.True(t, ok)
		assert.NoError(t, err)
		assert.NotEmpty(t, q.M()["$and"])
		assert.Contains(t, q.M()["$and"], primitive.M{"_x": primitive.ObjectID{2}, "_v": int64(1)})
	})

	t.Run("protection key is empty", func(t *testing.T) {
		q := &query.Query{}
		protected := &protection.Key{}

		ok, err := query.Push(q, protected)

		assert.True(t, ok)
		assert.NoError(t, err)
		assert.NotEmpty(t, q.M()["$and"])
		assert.Contains(t, q.M()["$and"], primitive.M{"_x": primitive.M{"$exists": false}, "_v": primitive.M{"$exists": false}})
	})
}

func TestPushPreloader(t *testing.T) {

	q := &query.Query{}
	preloader := query.Preload{"a", "b"}

	ok, err := query.Push(q, preloader)

	assert.True(t, ok)
	assert.NoError(t, err)

	p, hasPreloader := q.Preloader()

	assert.NotNil(t, p)
	assert.True(t, hasPreloader)
	assert.EqualValues(t, p, query.Preload{"a", "b"})
}
