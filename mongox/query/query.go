package query

import (
	"github.com/modern-go/reflect2"
	"github.com/valyala/bytebufferpool"
	"go.mongodb.org/mongo-driver/bson"
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
	onDecode  Callbacks
	onClose   Callbacks
	onCreate  Callbacks
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
		return nil
	}

	return q.limiter.Limit()
}

// Sorter is a sort rule for a query
func (q *Query) Sorter() (sort interface{}) {
	if q.sorter == nil {
		return nil
	}

	return q.sorter.Sort()
}

// Skipper is a skipper for a query
func (q *Query) Skipper() (skip *int64) {
	if q.skipper == nil {
		return nil
	}

	return q.skipper.Skip()
}

// Updater is an update command for a query
func (q *Query) Updater() (update primitive.M, err error) {
	if q.updater == nil {
		return primitive.M{}, nil
	}

	update = q.updater.Update()
	if reflect2.IsNil(update) {
		return primitive.M{}, nil
	}

	buffer := bytebufferpool.Get()
	defer bytebufferpool.Put(buffer)

	// convert update document to bson map values
	buffer.Reset()
	bsonBytes, err := bson.MarshalAppend(buffer.B, update)
	if err != nil {
		return primitive.M{}, err
	}
	update = primitive.M{} // reset update map and unmarshal bson bytes to it again
	err = bson.Unmarshal(bsonBytes, update)
	if err != nil {
		return primitive.M{}, err
	}

	return update, nil
}

// Preloader is a preloader list for a query
func (q *Query) Preloader() (preloads []string, ok bool) {
	if q.preloader == nil {
		return nil, false
	}

	preloads = q.preloader.Preload()

	return preloads, len(preloads) > 0
}

// OnDecode callback is called after the mongo decode function
func (q *Query) OnDecode() (callbacks Callbacks) {
	return q.onDecode
}

// OnClose callback is called after the mongox ends a loading procedure
func (q *Query) OnClose() (callbacks Callbacks) {
	return q.onClose
}

// OnCreate callback is called if the mongox creates a new document instance during loading
func (q *Query) OnCreate() (callbacks Callbacks) {
	return q.onClose
}

// Empty checks the query for any content
func (q *Query) Empty() (isEmpty bool) {
	return len(q.m) == 0
}

// M returns underlying query map
func (q *Query) M() (m primitive.M) {
	return q.m
}

// New creates a new query
func New() (query *Query) {
	return &Query{}
}
