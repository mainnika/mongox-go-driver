package query

// Skipper is a filter to skip the result
type Skipper interface {
	Skip() *int64
}

// Skip is a simple implementation of the Skipper filter
type Skip int64

var _ Skipper = Skip(0)

// Skip returns a skip number
func (l Skip) Skip() *int64 {

	lim := int64(l)
	if lim <= 0 {
		return nil
	}

	return &lim
}
