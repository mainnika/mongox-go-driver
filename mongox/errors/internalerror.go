package errors

import "fmt"

// InternalError error
type InternalError string

// Error message
func (ie InternalError) Error() string {
	return fmt.Sprintf("internal error, %s", string(ie))
}

// InternalErrorf function creates an instance of InternalError
func InternalErrorf(format string, params ...interface{}) error {
	return InternalError(fmt.Sprintf(format, params...))
}
