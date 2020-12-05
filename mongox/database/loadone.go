package database

import (
	"fmt"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// LoadOne function loads a first single target document by a query
func (d *Database) LoadOne(target interface{}, filters ...interface{}) (err error) {

	composed, err := query.Compose(append(filters, query.Limit(1))...)
	if err != nil {
		return
	}

	_, hasPreloader := composed.Preloader()
	ctx := query.WithContext(d.Context(), composed)

	var result *mongox.Cursor

	defer func() {

		if result != nil {
			closerr := result.Close(ctx)
			if err == nil {
				err = closerr
			}
		}

		invokerr := composed.OnClose().Invoke(ctx, target)
		if err == nil {
			err = invokerr
		}

		return
	}()

	if hasPreloader {
		result, err = d.createAggregateLoad(target, composed)
	} else {
		result, err = d.createSimpleLoad(target, composed)
	}
	if err != nil {
		return fmt.Errorf("can't create find result: %w", err)
	}

	hasNext := result.Next(ctx)
	if result.Err() != nil {
		err = result.Err()
		return
	}
	if !hasNext {
		return mongox.ErrNoDocuments
	}

	if created := base.Reset(target); created {
		err = composed.OnCreate().Invoke(ctx, target)
	}
	if err != nil {
		return
	}

	err = result.Decode(target)
	if err != nil {
		return
	}

	err = composed.OnDecode().Invoke(ctx, target)
	if err != nil {
		return
	}

	return
}
