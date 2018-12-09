package query

// Limiter is a filter to limit the result
type Limiter interface {
	Limit() int
}

// Limit is a simple implementation of the Limiter filter
type Limit int

var _ Limiter = Limit(0)

// Limit returns a limit
func (l Limit) Limit() int {

	return int(l)
}
