package edisegment

import (
	"fmt"
	"strings"
)

// ST represents the ST EDI segment
type ST struct {
	TransactionSetIdentifierCode string
	TransactionSetControlNumber  string
}

// String converts ST to its X12 single line string representation
func (s *ST) String(delimiter string) string {
	elements := []string{
		"ST",
		s.TransactionSetIdentifierCode,
		s.TransactionSetControlNumber,
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the ST struct
func (s *ST) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("ST: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}
	s.TransactionSetIdentifierCode = elements[0]
	s.TransactionSetControlNumber = elements[1]
	return nil
}
