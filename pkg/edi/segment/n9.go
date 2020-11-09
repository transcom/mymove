package edisegment

import (
	"fmt"
)

// N9 represents the N9 EDI segment
type N9 struct {
	ReferenceIdentificationQualifier string `validate:"oneof=DY CN CT PQ OQ 1W ML 3L PO"`
	ReferenceIdentification          string `validate:"min=1,max=30"`
	FreeFormDescription              string `validate:"omitempty,min=1,max=45"`
	Date                             string `validate:"omitempty,datetime=20060102"`
}

// StringArray converts N9 to an array of strings
func (s *N9) StringArray() []string {
	return []string{
		"N9",
		s.ReferenceIdentificationQualifier,
		s.ReferenceIdentification,
		s.FreeFormDescription,
		s.Date,
	}
}

// Parse parses an X12 string that's split into an array into the N9 struct
func (s *N9) Parse(elements []string) error {
	numElements := len(elements)
	if numElements != 2 && numElements != 3 && numElements != 4 {
		return fmt.Errorf("N9: Wrong number of fields, expected 2 or 3 or 4, got %d", numElements)
	}
	s.ReferenceIdentificationQualifier = elements[0]
	s.ReferenceIdentification = elements[1]
	if numElements > 2 {
		s.FreeFormDescription = elements[2]
	}
	if numElements > 3 {
		s.Date = elements[3]
	}
	return nil
}
