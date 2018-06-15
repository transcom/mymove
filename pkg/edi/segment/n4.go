package edisegment

import (
	"fmt"
	"strings"
)

// N4 represents the N4 EDI segment
type N4 struct {
	CityName            string
	StateOrProvinceCode string
	PostalCode          string
	CountryCode         string
	LocationQualifier   string
	LocationIdentifier  string
}

// String converts N4 to its X12 single line string representation
func (s *N4) String(delimiter string) string {
	elements := []string{
		"N4",
		s.CityName,
		s.StateOrProvinceCode,
		s.PostalCode,
		s.CountryCode,
		s.LocationQualifier,
		s.LocationIdentifier,
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the N4 struct
func (s *N4) Parse(parts []string) error {
	expectedNumElements := 6
	if len(parts) != expectedNumElements {
		fmt.Printf("what")
		return fmt.Errorf("N4: Wrong number of fields, expected %d, got %d", expectedNumElements, len(parts))
	}

	s.CityName = parts[0]
	s.StateOrProvinceCode = parts[1]
	s.PostalCode = parts[2]
	s.CountryCode = parts[3]
	s.LocationQualifier = parts[4]
	s.LocationIdentifier = parts[5]
	return nil
}
