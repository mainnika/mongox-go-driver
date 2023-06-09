package query

import (
	"fmt"

	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base/protection"
)

type applyFilterFunc func(query *Query, filter interface{}) (ok bool)

// Compose is a function to compose filters into a single query
func Compose(filters ...interface{}) (query *Query, err error) {
	query = &Query{}
	for _, filter := range filters {
		ok, err := Push(query, filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter %v, %w", filter, err)
		}
		if !ok {
			panic(fmt.Errorf("unknown filter %v", filter))
		}
	}

	return query, nil
}

// Push applies single filter to a query
func Push(query *Query, filter interface{}) (ok bool, err error) {
	emptyFilter := reflect2.IsNil(filter)
	if emptyFilter {
		return true, nil
	}

	validator, hasValidator := filter.(Validator)
	if hasValidator {
		err := validator.Validate()
		if err != nil {
			return false, fmt.Errorf("while validating filter %v, %w", filter, err)
		}
	}

	ok = false // true if at least one filter was applied
	for _, applier := range []applyFilterFunc{
		applyBson,
		applyLimit,
		applySort,
		applySkip,
		applyProtection,
		applyPreloader,
		applyUpdater,
		applyCallbacks,
	} {
		ok = applier(query, filter) || ok
	}

	return ok, nil
}

// applyBson is a fallback for a custom primitive.M
func applyBson(query *Query, filter interface{}) (ok bool) {
	if filter, ok := filter.(primitive.M); ok {
		query.And(filter)
		return true
	}

	return false
}

// applyLimits extends query with a limiter
func applyLimit(query *Query, filter interface{}) (ok bool) {
	if filter, ok := filter.(Limiter); ok {
		query.limiter = filter
		return true
	}

	return false
}

// applySort extends query with a sort rule
func applySort(query *Query, filter interface{}) (ok bool) {
	if filter, ok := filter.(Sorter); ok {
		query.sorter = filter
		return true
	}

	return false
}

// applySkip extends query with a skip number
func applySkip(query *Query, filter interface{}) (ok bool) {
	if filter, ok := filter.(Skipper); ok {
		query.skipper = filter
		return true
	}

	return false
}

func applyProtection(query *Query, filter interface{}) (ok bool) {
	keyDoc := primitive.M{}
	switch filter := filter.(type) {
	case protection.Key:
		filter.Inject(keyDoc)
	case *protection.Key:
		filter.Inject(keyDoc)
	default:
		return false
	}

	query.And(keyDoc)

	return true
}

func applyPreloader(query *Query, filter interface{}) (ok bool) {
	if filter, ok := filter.(Preloader); ok {
		query.preloader = filter
		return true
	}

	return false
}

func applyUpdater(query *Query, filter interface{}) (ok bool) {
	if filter, ok := filter.(Updater); ok {
		query.updater = filter
		return true
	}

	return false
}

func applyCallbacks(query *Query, filter interface{}) (ok bool) {
	switch callback := filter.(type) {
	case OnDecode:
		query.onDecode = append(query.onDecode, Callback(callback))
	case OnClose:
		query.onClose = append(query.onClose, Callback(callback))
	default:
		return false
	}

	return true
}
