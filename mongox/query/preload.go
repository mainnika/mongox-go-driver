package query

// Preloader is a filter to skip the result
type Preloader interface {
	Preload() []string
}

// Preload is a simple implementation of the Skipper filter
type Preload []string

var _ Preloader = Preload{}

// Preload returns a preload list
func (l Preload) Preload() []string {
	return l
}
