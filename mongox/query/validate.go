package query

// Validator is a filter to validate the filter
type Validator interface {
	Validate() (err error)
}
