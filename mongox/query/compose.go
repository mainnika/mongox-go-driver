package query

import (
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mongodb/mongo-go-driver/bson"
)

// ComposeQuery is a function to compose filters into a single query
func Compose(filters ...interface{}) *Query {

	q := &Query{}

	for _, f := range filters {

		ok := false
		ok = ok || applyBson(q, f)
		ok = ok || applyLimit(q, f)
		ok = ok || applySort(q, f)
		ok = ok || applySkip(q, f)

		if !ok {
			panic(errors.InternalErrorf("unknown filter %v", f))
		}
	}

	return q
}

// applyBson is a fallback for a custom bson.M
func applyBson(q *Query, f interface{}) bool {

	if f, ok := f.(bson.M); ok {
		q.And(f)
		return true
	}

	return false
}

// applyLimits extends query with a limiter
func applyLimit(q *Query, f interface{}) bool {

	if f, ok := f.(Limiter); ok {
		q.limiter = f
		return true
	}

	return false
}

// applySort extends query with a sort rule
func applySort(q *Query, f interface{}) bool {

	if f, ok := f.(Sorter); ok {
		q.sorter = f
		return true
	}

	return false
}

// applySkip extends query with a skip number
func applySkip(q *Query, f interface{}) bool {

	if f, ok := f.(Skipper); ok {
		q.skipper = f
		return true
	}

	return false
}
