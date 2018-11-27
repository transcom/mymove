package edisegment

import (
	"fmt"
)

// N3 represents the N3 EDI segment
type N3 struct {
	AddressInformation1 string
	AddressInformation2 string
}

// StringArray converts N3 to an array of strings
func (s *N3) StringArray() []string {
	return []string{
		"N3",
		s.AddressInformation1,
		s.AddressInformation2,
	}
}

// Parse parses an X12 string that's split into an array into the N3 struct
func (s *N3) Parse(parts []string) error {
	expectedNumElements := 2
	if len(parts) != expectedNumElements {
		return fmt.Errorf("N3: Wrong number of fields, expected %d, got %d", expectedNumElements, len(parts))
	}

	s.AddressInformation1 = parts[0]
	s.AddressInformation2 = parts[1]
	return nil
}
