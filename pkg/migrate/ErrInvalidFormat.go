package migrate

import (
	"fmt"
)

// ErrInvalidFormat is an error for an invalid migration format.
// Only SQL and Fizz are currently supported
type ErrInvalidFormat struct {
	Value string
}

func (e *ErrInvalidFormat) Error() string {
	return fmt.Sprintf("invalid format %q, expecting sql or fizz", e.Value)
}
