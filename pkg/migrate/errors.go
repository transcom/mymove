package migrate

import (
	"fmt"
)

// ErrInvalidDirection is an error for an invalid direction
type ErrInvalidDirection struct {
	Value string
}

func (e *ErrInvalidDirection) Error() string {
	return fmt.Sprintf("invalid direction %q, expecting up", e.Value)
}

// ErrInvalidFormat is an error for an invalid migration format.
// Only SQL and Fizz are currently supported
type ErrInvalidFormat struct {
	Value string
}

func (e *ErrInvalidFormat) Error() string {
	return fmt.Sprintf("invalid format %q, expecting sql or fizz", e.Value)
}

// ErrInvalidPath is an error for an invalid path
type ErrInvalidPath struct {
	Value string
}

func (e *ErrInvalidPath) Error() string {
	return fmt.Sprintf("invalid path %q, should start with file:// or s3://", e.Value)
}
