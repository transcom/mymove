package edisegment

import (
	"fmt"
)

// FA2 represents the FA2 EDI segment
type FA2 struct {
	BreakdownStructureDetailCode string
	FinancialInformationCode     string
}

// StringArray converts FA2 to an array of strings
func (s *FA2) StringArray() []string {
	return []string{"FA2", s.BreakdownStructureDetailCode, s.FinancialInformationCode}
}

// Parse parses an X12 string that's split into an array into the FA2 struct
func (s *FA2) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("FA2: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.BreakdownStructureDetailCode = elements[0]
	s.FinancialInformationCode = elements[1]
	return nil
}
