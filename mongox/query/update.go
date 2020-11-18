package query

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Updater is a filter to update the data
type Updater interface {
	Update() (update primitive.M)
}

// Update is a simple implementations of the Updater filter
type Update primitive.M

var _ Updater = &Update{}

// Update returns an update command
func (u Update) Update() (update primitive.M) {
	return primitive.M(u)
}
