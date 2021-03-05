package edisegment

import (
	"fmt"
)

// AK2 represents the AK2 EDI segment
type AK2 struct {
	TransactionSetIdentifierCode string `validate:"eq=858"`
	TransactionSetControlNumber  string `validate:"min=4,max=9"`
}

// StringArray converts AK2 to an array of strings
func (s *AK2) StringArray() []string {
	return []string{
		"AK2",
		s.TransactionSetIdentifierCode,
		s.TransactionSetControlNumber,
	}
}

// Parse parses an AK2 string that's split into an array into the AK2 struct
func (s *AK2) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("AK2: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}
	s.TransactionSetIdentifierCode = elements[0]
	s.TransactionSetControlNumber = elements[1]
	return nil
}
