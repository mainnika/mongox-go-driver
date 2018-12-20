package errors

import "fmt"

// Malformed error
type Malformed string

// Error message
func (m Malformed) Error() string {
	return fmt.Sprintf("Malformed, %s", string(m))
}

// Malformedf creates an instance of Malformed
func Malformedf(format string, params ...interface{}) error {
	return Malformed(fmt.Sprintf(format, params...))
}
