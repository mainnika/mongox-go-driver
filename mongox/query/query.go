package query

import (
	"github.com/mongodb/mongo-go-driver/bson"

	"reflect"
)

// Query is an enchanched bson.M map
type Query struct {
	m       bson.M
	limiter Limiter
	sorter  Sorter
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

// Limiter is a limit function for a query
func (q *Query) Limiter() Limiter {

	return q.limiter
}

// Sorter is a sort rule for a query
func (q *Query) Sorter() Sorter {

	return q.sorter
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