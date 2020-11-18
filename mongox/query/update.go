package query

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Updater is a filter to update the data
type Updater interface {
	Update() (update bson.A)
}

// Update is a simple implementations of the Updater filter
type Update bson.M

var _ Updater = &Update{}

// Update returns an update command
func (u Update) Update() (update bson.A) {
	return bson.A{bson.M(u)}
}
