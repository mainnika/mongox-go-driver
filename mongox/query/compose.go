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
		ok = ok || applyLimits(q, f)

		if !ok {
			panic(errors.InternalErrorf("unknown filter %v", f))
		}
	}

	return q
}

// applyBson is a fallback for a custom bson.M
func applyBson(q *Query, f interface{}) bool {

	switch f := f.(type) {
	case bson.M:
		q.And(f)
	default:
		return false
	}

	return true
}

// applyLimits extends query with contol functions
func applyLimits(q *Query, f interface{}) bool {

	switch f := f.(type) {
	case Limiter:
		q.limiter = f
	case Sorter:
		q.sorter = f
	case Skipper:
		q.skipper = f
	default:
		return false
	}

	return true
}
