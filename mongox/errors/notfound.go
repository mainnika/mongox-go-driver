package errors

import "fmt"

// NotFound error
type NotFound string

// Error message
func (nf NotFound) Error() string {
	return fmt.Sprintf("can not find, %s", string(nf))
}

// NotFoundErrorf function creates an instance of BadRequestError
func NotFoundErrorf(format string, params ...interface{}) error {
	return NotFound(fmt.Sprintf(format, params...))
}
