package edisegment

import (
	"fmt"
)

// N1 represents the N1 EDI segment
type N1 struct {
	EntityIdentifierCode        string
	Name                        string
	IdentificationCodeQualifier string
	IdentificationCode          string
}

// StringArray converts N1 to an array of strings
func (s *N1) StringArray() []string {
	return []string{
		"N1",
		s.EntityIdentifierCode,
		s.Name,
		s.IdentificationCodeQualifier,
		s.IdentificationCode,
	}
}

// Parse parses an X12 string that's split into an array into the N1SF struct
func (s *N1) Parse(parts []string) error {
	numElements := len(parts)
	if numElements != 2 && numElements != 4 {
		return fmt.Errorf("N1SF: Wrong number of elements, expected 2 or 4, got %d", numElements)
	}

	s.EntityIdentifierCode = parts[0]
	s.Name = parts[1]
	if numElements == 4 {
		s.IdentificationCodeQualifier = parts[2]
		s.IdentificationCode = parts[3]
	}
	return nil
}
