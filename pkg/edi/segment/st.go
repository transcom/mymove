package edisegment

import (
	"fmt"
)

// ST represents the ST EDI segment
type ST struct {
	TransactionSetIdentifierCode string `validate:"oneof=858 997 824 810"`
	TransactionSetControlNumber  string `validate:"min=4,max=9"`
}

// StringArray converts ST to an array of strings
func (s *ST) StringArray() []string {
	return []string{
		"ST",
		s.TransactionSetIdentifierCode,
		s.TransactionSetControlNumber,
	}
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
