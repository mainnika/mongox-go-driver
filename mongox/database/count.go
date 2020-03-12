package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// Count function counts documents in the database by query
// target is used only to get collection by tag so it'd be better to use nil ptr here
func (d *Database) Count(target interface{}, filters ...interface{}) (int64, error) {

	collection := d.GetCollectionOf(target)
	opts := options.Count()
	composed := query.Compose(filters...)

	opts.Limit = composed.Limiter()
	opts.Skip = composed.Skipper()

	result, err := collection.CountDocuments(d.Context(), composed.M(), opts)
	if err == mongox.ErrNoDocuments {
		return 0, err
	}
	if err != nil {
		return 0, fmt.Errorf("can't decode desult: %w", err)
	}

	return result, nil
}
