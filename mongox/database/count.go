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
		return
	}

	collection := d.GetCollectionOf(target)
	ctx := query.WithContext(d.Context(), composed)

	opts := options.Count()
	opts.Limit = composed.Limiter()
	opts.Skip = composed.Skipper()

	result, err = collection.CountDocuments(ctx, composed.M(), opts)

	_ = composed.OnClose().Invoke(ctx, target)

	return
}
