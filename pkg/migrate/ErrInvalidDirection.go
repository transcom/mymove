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
