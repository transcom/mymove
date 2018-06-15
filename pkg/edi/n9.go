package edi

import (
	"fmt"
	"strings"
)

// N9 represents the N9 EDI segment
type N9 struct {
	ReferenceIdentificationQualifier string
	ReferenceIdentification          string
	FreeFormDescription              string
}

// String converts N9 to its X12 single line string representation
func (s *N9) String(delimiter string) string {
	elements := []string{
		"N9",
		s.ReferenceIdentificationQualifier,
		s.ReferenceIdentification,
		s.FreeFormDescription,
	}
	return strings.Join(elements, delimiter) + "\n"
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
		s.FreeFormDescription = elements[3]
	}
	return nil
}
