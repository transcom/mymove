package edisegment

import (
	"fmt"
	"strings"
)

// HL represents the HL EDI segment
type HL struct {
	HierarchicalIDNumber       string
	HierarchicalParentIDNumber string
	HierarchicalLevelCode      string
}

// String converts HL to its X12 single line string representation
func (s *HL) String(delimiter string) string {
	elements := []string{
		"HL",
		s.HierarchicalIDNumber,
		s.HierarchicalParentIDNumber,
		s.HierarchicalLevelCode,
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the HL struct
func (s *HL) Parse(elements []string) error {
	expectedNumElements := 3
	if len(elements) != expectedNumElements {
		return fmt.Errorf("HL: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.HierarchicalIDNumber = elements[0]
	s.HierarchicalParentIDNumber = elements[1]
	s.HierarchicalLevelCode = elements[2]
	return nil
}
