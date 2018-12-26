package query

// Limiter is a filter to limit the result
type Limiter interface {
	Limit() *int64
}

// Limit is a simple implementation of the Limiter filter
type Limit int64

var _ Limiter = Limit(0)

// Limit returns a limit
func (l Limit) Limit() *int64 {

	lim := int64(l)
	if lim <= 0 {
		return nil
	}

	return &lim
}
