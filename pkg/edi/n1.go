package edi

import (
	"fmt"
	"strings"
)

// N1 represents the N1 EDI segment
type N1 struct {
	EntityIdentifierCode        string
	Name                        string
	IdentificationCodeQualifier string
	IdentificationCode          string
}

// String converts N1 to its X12 single line string representation
func (s *N1) String(delimiter string) string {
	elements := []string{
		"N1",
		s.EntityIdentifierCode,
		s.Name,
		s.IdentificationCodeQualifier,
		s.IdentificationCode,
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the N1 struct
func (s *N1) Parse(parts []string) error {
	numElements := len(parts)
	if numElements != 2 && numElements != 4 {
		return fmt.Errorf("N1: Wrong number of elements, expected 2 or 4, got %d", numElements)
	}

	s.EntityIdentifierCode = parts[0]
	s.Name = parts[1]
	if numElements == 4 {
		s.IdentificationCodeQualifier = parts[2]
		s.IdentificationCode = parts[3]
	}
	return nil
}
