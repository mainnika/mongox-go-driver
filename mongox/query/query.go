package query

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Query is an enchanched primitive.M map
type Query struct {
	m         primitive.M
	limiter   Limiter
	sorter    Sorter
	skipper   Skipper
	preloader Preloader
	updater   Updater
	ondecode  Callbacks
	onclose   Callbacks
}

// And function pushes the elem query to the $and array of the query
func (q *Query) And(elem primitive.M) (query *Query) {

	if q.m == nil {
		q.m = primitive.M{}
	}

	queries, exists := q.m["$and"].(primitive.A)

	if !exists {
		q.m["$and"] = primitive.A{elem}
		return q
	}

	q.m["$and"] = append(queries, elem)

	return q
}

// Limiter returns limiter value or nil
func (q *Query) Limiter() (limit *int64) {

	if q.limiter == nil {
		return
	}

	return q.limiter.Limit()
}

// Sorter is a sort rule for a query
func (q *Query) Sorter() (sort interface{}) {

	if q.sorter == nil {
		return
	}

	return q.sorter.Sort()
}

// Skipper is a skipper for a query
func (q *Query) Skipper() (skip *int64) {

	if q.skipper == nil {
		return
	}

	return q.skipper.Skip()
}

// Updater is an update command for a query
func (q *Query) Updater() (update primitive.A) {

	if q.updater == nil {
		return
	}

	return q.updater.Update()
}

// Preloader is a preloader list for a query
func (q *Query) Preloader() (ok bool, preloads []string) {

	if q.preloader == nil {
		return false, nil
	}

	preloads = q.preloader.Preload()
	ok = len(preloads) > 0

	return
}

// OnDecode callback is called after the mongo decode function
func (q *Query) OnDecode() (callbacks Callbacks) {
	return q.ondecode
}

func (q *Query) OnClose() (callbacks Callbacks) {
	return q.onclose
}

// Empty checks the query for any content
func (q *Query) Empty() (isEmpty bool) {
	return len(q.m) == 0
}

// M returns underlying query map
func (q *Query) M() (m primitive.M) {
	return q.m
}
