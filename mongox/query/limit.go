package query

// Limiter is a filter to limit the result
type Limiter interface {
	Limit() (limit *int64)
}

// Limit is a simple implementation of the Limiter filter
type Limit int64

var _ Limiter = Limit(0)

// Limit returns a limit
func (l Limit) Limit() (limit *int64) {
	if l <= 0 {
		return nil
	}

	limit = new(int64)
	*limit = int64(l)

	return limit
}
