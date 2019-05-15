package query

import (
	"go.mongodb.org/mongo-driver/bson"

	"reflect"
)

// Query is an enchanched bson.M map
type Query struct {
	m       bson.M
	limiter Limiter
	sorter  Sorter
	skipper Skipper
}

// And function pushes the elem query to the $and array of the query
func (q *Query) And(elem bson.M) *Query {

	if q.m == nil {
		q.m = bson.M{}
	}

	queries, exists := q.m["$and"].(bson.A)

	if !exists {
		q.m["$and"] = bson.A{elem}
		return q
	}

	q.m["$and"] = append(queries, elem)

	return q
}

// Limiter returns limiter value or nil
func (q *Query) Limiter() *int64 {

	if q.limiter == nil {
		return nil
	}

	return q.limiter.Limit()
}

// Sorter is a sort rule for a query
func (q *Query) Sorter() interface{} {

	if q.sorter == nil {
		return nil
	}

	return q.sorter.Sort()
}

// Skipper is a skipper for a query
func (q *Query) Skipper() *int64 {

	if q.skipper == nil {
		return nil
	}

	return q.skipper.Skip()
}

// Empty checks the query for any content
func (q *Query) Empty() bool {

	qv := reflect.ValueOf(q)
	keys := qv.MapKeys()

	return len(keys) == 0
}

// M returns underlying query map
func (q *Query) M() bson.M {
	return q.m
}
