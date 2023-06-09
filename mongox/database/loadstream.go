package database

import (
	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// LoadStream function loads documents one by one into a target channel
func (d *Database) LoadStream(target interface{}, filters ...interface{}) (loader mongox.StreamLoader, err error) {
	composed, err := query.Compose(filters...)
	if err != nil {
		return
	}

	ctx := query.WithContext(d.Context(), composed)

	cur, err := d.createCursor(target, composed)
	if err != nil {
		return nil, err
	}

	loader = &StreamLoader{cur: cur, ctx: ctx, query: composed}

	return loader, nil
}
