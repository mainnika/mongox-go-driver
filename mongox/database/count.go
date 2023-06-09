package database

import (
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// Count function counts documents in the database by query
// target is used only to get collection by tag so it'd be better to use nil ptr here
func (d *Database) Count(target interface{}, filters ...interface{}) (result int64, err error) {
	composed, err := query.Compose(filters...)
	if err != nil {
		return -1, err
	}

	collection, err := d.GetCollectionOf(target)
	if err != nil {
		return -1, err
	}
	ctx := query.WithContext(d.Context(), composed)

	m := composed.M()

	opts := options.Count()
	opts.Limit = composed.Limiter()
	opts.Skip = composed.Skipper()

	defer func() { _ = composed.OnClose().Invoke(ctx, target) }()

	result, err = collection.CountDocuments(ctx, m, opts)
	if err != nil {
		return -1, err
	}

	return result, nil
}
