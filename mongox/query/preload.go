package query

// Preloader is a filter to skip the result
type Preloader interface {
	Preload() (preloads []string)
}

// Preload is a simple implementation of the Skipper filter
type Preload []string

var _ Preloader = Preload{}

// Preload returns a preload list
func (l Preload) Preload() (preloads []string) {
	return l
}
