package edisegment

import (
	"fmt"
	"strings"
)

// FA2 represents the FA2 EDI segment
type FA2 struct {
	BreakdownStructureDetailCode string
	FinancialInformationCode     string
}

// String converts FA2 to its X12 single line string representation
func (s *FA2) String(delimiter string) string {
	return strings.Join([]string{"FA2", s.BreakdownStructureDetailCode, s.FinancialInformationCode}, delimiter)
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
