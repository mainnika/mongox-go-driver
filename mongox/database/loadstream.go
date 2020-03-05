package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// LoadStream function loads documents one by one into a target channel
func (d *Database) LoadStream(target interface{}, filters ...interface{}) (mongox.StreamLoader, error) {

	var cursor *mongo.Cursor
	var err error

	composed := query.Compose(filters...)
	hasPreloader, _ := composed.Preloader()

	if hasPreloader {
		cursor, err = d.createAggregateLoad(target, composed)
	} else {
		cursor, err = d.createSimpleLoad(target, composed)
	}
	if err != nil {
		return nil, fmt.Errorf("can't create find result: %w", err)
	}

	l := &StreamLoader{Cursor: cursor, ctx: d.Context(), target: target}

	return l, nil
}
