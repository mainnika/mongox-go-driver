package database

import (
	"fmt"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// LoadStream function loads documents one by one into a target channel
func (d *Database) LoadStream(target interface{}, filters ...interface{}) (loader mongox.StreamLoader, err error) {

	var cursor *mongox.Cursor

	composed := query.Compose(filters...)
	hasPreloader, _ := composed.Preloader()

	if hasPreloader {
		cursor, err = d.createAggregateLoad(target, composed)
	} else {
		cursor, err = d.createSimpleLoad(target, composed)
	}
	if err != nil {
		err = fmt.Errorf("can't create find result: %w", err)
		return
	}

	loader = &StreamLoader{cur: cursor, ctx: d.Context(), target: target}

	return
}
