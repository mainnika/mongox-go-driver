package query

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sorter is a filter to sort the data before query
type Sorter interface {
	Sort() (sort primitive.M)
}

// Sort is a simple implementations of the Sorter filter
type Sort primitive.M

var _ Sorter = &Sort{}

// Sort returns a slice of fields which have to be sorted
func (f Sort) Sort() (sort primitive.M) {
	return primitive.M(f)
}
