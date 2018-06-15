package edi

import (
	"fmt"
	"strings"
)

// G62 represents the G62 EDI segment
type G62 struct {
	DateQualifier string
	Date          string
}

// String converts G62 to its X12 single line string representation
func (s *G62) String(delimiter string) string {
	return strings.Join([]string{"G62", s.DateQualifier, s.Date}, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the G62 struct
func (s *G62) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("G62: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.DateQualifier = elements[0]
	s.Date = elements[1]
	return nil
}
