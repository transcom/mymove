package edisegment

import (
	"fmt"
)

// FA2 represents the FA2 EDI segment
type FA2 struct {
	BreakdownStructureDetailCode string `validate:"oneof=ZZ TA A1 A2 A3 A4 A5 A6 B1 B2 B3 C1 C2 D1 D4 D6 D7 E1 E2 E3 F1 F3 G2 I1 J1 K6 L1 M1 N1 P5"`
	FinancialInformationCode     string `validate:"min=1,max=80"`
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
