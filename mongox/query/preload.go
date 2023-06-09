package query

// Preloader is a filter to preload the result
type Preloader interface {
	Preload() (preloads []string)
}

// Preload is a simple implementation of the Preloader filter
type Preload []string

var _ Preloader = Preload{}

// Preload returns a preload list
func (l Preload) Preload() (preloads []string) {
	return l
}
