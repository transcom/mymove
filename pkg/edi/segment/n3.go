package edisegment

import (
	"fmt"
	"strings"
)

// N3 represents the N3 EDI segment
type N3 struct {
	AddressInformation1 string
	AddressInformation2 string
}

// String converts N3 to its X12 single line string representation
func (s *N3) String(delimiter string) string {
	elements := []string{
		"N3",
		s.AddressInformation1,
		s.AddressInformation2,
	}
	return strings.Join(elements, delimiter) + "\n"
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
