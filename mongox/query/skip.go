package query

// Skipper is a filter to skip the result
type Skipper interface {
	Skip() (skip *int64)
}

// Skip is a simple implementation of the Skipper filter
type Skip int64

var _ Skipper = Skip(0)

// Skip returns a skip number
func (l Skip) Skip() (skip *int64) {

	if l <= 0 {
		return
	}

	skip = new(int64)
	*skip = int64(l)

	return
}
