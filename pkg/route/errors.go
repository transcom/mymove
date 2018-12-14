package route

import (
	"fmt"
)

// UnsupportedPostalCode represents a postal code that cannot be handled by the application.
type UnsupportedPostalCode struct {
	postalCode string
}

// NewUnsupportedPostalCodeError creates a new UnsupportedPostalCode error.
func NewUnsupportedPostalCodeError(postalCode string) *UnsupportedPostalCode {
	return &UnsupportedPostalCode{
		postalCode: postalCode,
	}
}

func (e *UnsupportedPostalCode) Error() string {
	return fmt.Sprintf("Unsupported postal code (%s)", e.postalCode)
}
