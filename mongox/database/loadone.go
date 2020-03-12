package database

import (
	"fmt"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// LoadOne function loads a first single target document by a query
func (d *Database) LoadOne(target interface{}, filters ...interface{}) error {

	composed := query.Compose(append(filters, query.Limit(1))...)
	hasPreloader, _ := composed.Preloader()

	var result *mongox.Cursor
	var err error

	if hasPreloader {
		result, err = d.createAggregateLoad(target, composed)
	} else {
		result, err = d.createSimpleLoad(target, composed)
	}
	if err != nil {
		return fmt.Errorf("can't create find result: %w", err)
	}

	hasNext := result.Next(d.Context())
	if result.Err() != nil {
		return err
	}
	if !hasNext {
		return mongox.ErrNoDocuments
	}

	base.Reset(target)

	return result.Decode(target)
}
