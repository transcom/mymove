package edisegment

import (
	"fmt"
)

// BGN represents the BGN EDI segment
type BGN struct {
	TransactionSetPurposeCode string `validate:"eq=11"`
	ReferenceIdentification   string `validate:"min=1,max=30"`
	Date                      string `validate:"datetime=20060102"`
}

// StringArray converts BGN to an array of strings
func (s *BGN) StringArray() []string {
	return []string{
		"BGN",
		s.TransactionSetPurposeCode,
		s.ReferenceIdentification,
		s.Date,
	}
}

// Parse parses an X12 string that's split into an array into the BGN struct
func (s *BGN) Parse(elements []string) error {
	expectedNumElements := 3
	if len(elements) != expectedNumElements {
		return fmt.Errorf("BGN: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.TransactionSetPurposeCode = elements[0]
	s.ReferenceIdentification = elements[1]
	s.Date = elements[2]

	return nil
}
