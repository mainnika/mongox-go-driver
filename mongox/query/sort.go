package query

import (
	"github.com/mongodb/mongo-go-driver/bson"
)

// Sorter is a filter to sort the data before query
type Sorter interface {
	Sort() bson.M
}

// Sort is a simple implementations of the Sorter filter
type Sort bson.M

var _ Sorter = &Sort{}

// Sort returns a slice of fields which have to be sorted
func (f Sort) Sort() bson.M {
	return bson.M(f)
}
