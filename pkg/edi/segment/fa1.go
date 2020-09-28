package edisegment

import (
	"fmt"
)

// FA1 represents the FA1 EDI segment
type FA1 struct {
	AgencyQualifierCode string `validate:"eq=DF"`
}

// StringArray converts FA1 to an array of strings
func (s *FA1) StringArray() []string {
	return []string{"FA1", s.AgencyQualifierCode}
}

// Parse parses an X12 string that's split into an array into the FA1 struct
func (s *FA1) Parse(elements []string) error {
	expectedNumElements := 1
	if len(elements) != expectedNumElements {
		return fmt.Errorf("FA1: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.AgencyQualifierCode = elements[0]
	return nil
}
