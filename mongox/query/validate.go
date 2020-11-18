package query

// Valider is a filter to validate the filter
type Valider interface {
	Valid() (err error)
}
