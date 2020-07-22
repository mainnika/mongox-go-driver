package query

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base/protection"
	"github.com/mainnika/mongox-go-driver/v2/mongox/utils"
)

// Compose is a function to compose filters into a single query
func Compose(filters ...interface{}) (query *Query) {

	query = &Query{}

	for _, f := range filters {
		if !Push(query, f) {
			panic(fmt.Errorf("unknown filter %v", f))
		}
	}

	return
}

// Push applies single filter to a query
func Push(q *Query, f interface{}) (ok bool) {

	if utils.IsNil(f) {
		return true
	}

	ok = false
	ok = ok || applyBson(q, f)
	ok = ok || applyLimit(q, f)
	ok = ok || applySort(q, f)
	ok = ok || applySkip(q, f)
	ok = ok || applyProtection(q, f)
	ok = ok || applyPreloader(q, f)
	ok = ok || applyCallbacks(q, f)

	return ok
}

// applyBson is a fallback for a custom bson.M
func applyBson(q *Query, f interface{}) (ok bool) {

	if f, ok := f.(bson.M); ok {
		q.And(f)
		return true
	}

	return false
}

// applyLimits extends query with a limiter
func applyLimit(q *Query, f interface{}) (ok bool) {

	if f, ok := f.(Limiter); ok {
		q.limiter = f
		return true
	}

	return false
}

// applySort extends query with a sort rule
func applySort(q *Query, f interface{}) (ok bool) {

	if f, ok := f.(Sorter); ok {
		q.sorter = f
		return true
	}

	return false
}

// applySkip extends query with a skip number
func applySkip(q *Query, f interface{}) (ok bool) {

	if f, ok := f.(Skipper); ok {
		q.skipper = f
		return true
	}

	return false
}

func applyProtection(q *Query, f interface{}) (ok bool) {

	var x *primitive.ObjectID
	var v *int64

	switch f := f.(type) {
	case protection.Key:
		x = &f.X
		v = &f.V
	case *protection.Key:
		x = &f.X
		v = &f.V

	default:
		return false
	}

	if x.IsZero() {
		q.And(primitive.M{"_x": primitive.M{"$exists": false}})
		q.And(primitive.M{"_v": primitive.M{"$exists": false}})
	} else {
		q.And(primitive.M{"_x": *x})
		q.And(primitive.M{"_v": *v})
	}

	return true
}

func applyPreloader(q *Query, f interface{}) (ok bool) {

	if f, ok := f.(Preloader); ok {
		q.preloader = f
		return true
	}

	return false
}

func applyCallbacks(q *Query, f interface{}) (ok bool) {

	switch cb := f.(type) {
	case OnDecode:
		q.ondecode = append(q.ondecode, Callback(cb))
	case OnClose:
		q.onclose = append(q.onclose, Callback(cb))
	default:
		return
	}

	ok = true

	return
}
