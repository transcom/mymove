package edisegment

import (
	"fmt"
	"strings"
)

// B3A represents the B3A EDI segment
type B3A struct {
	TransactionTypeCode string
}

// String converts B3A to its X12 single line string representation
func (s *B3A) String(delimiter string) string {
	return strings.Join([]string{"B3A", s.TransactionTypeCode}, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the B3A struct
func (s *B3A) Parse(elements []string) error {
	expectedNumElements := 1
	if len(elements) != expectedNumElements {
		return fmt.Errorf("B3A: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.TransactionTypeCode = elements[0]
	return nil
}
