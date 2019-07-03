package migrate

import (
	"fmt"
)

// ErrInvalidPath is an error for an invalid path
type ErrInvalidPath struct {
	Value string
}

func (e *ErrInvalidPath) Error() string {
	return fmt.Sprintf("invalid path %q, should start with file:// or s3://", e.Value)
}
