package database

import (
	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// LoadOne function loads a first single target document by a query
func (d *Database) LoadOne(target interface{}, filters ...interface{}) (err error) {
	composed, err := query.Compose(append(filters, query.Limit(1))...)
	if err != nil {
		return err
	}

	ctx := query.WithContext(d.Context(), composed)

	defer func() { _ = composed.OnClose().Invoke(ctx, target) }()

	cur, err := d.createCursor(target, composed)
	if err != nil {
		return err
	}
	defer func() { _ = cur.Close(ctx) }()

	hasNext := cur.Next(ctx)
	if cur.Err() != nil {
		return cur.Err()
	}
	if !hasNext {
		return mongox.ErrNoDocuments
	}

	if created := base.Reset(target); created {
		_ = composed.OnCreate().Invoke(ctx, target)
	}

	err = cur.Decode(target)
	if err != nil {
		return err
	}

	_ = composed.OnDecode().Invoke(ctx, target)

	return nil
}
