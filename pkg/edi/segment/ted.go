package edisegment

import (
	"fmt"
)

// TED represents the TED EDI segment
type TED struct {
	ApplicationErrorConditionCode string `validate:"oneof=007 812 832 DUP IID INC K MJ PPD T ZZZ"`
	FreeFormMessage               string `validate:"omitempty,max=60"`
}

// StringArray converts TED to an array of strings
func (s *TED) StringArray() []string {
	return []string{
		"TED",
		s.ApplicationErrorConditionCode,
		s.FreeFormMessage,
	}
}

// Parse parses an X12 string that's split into an array into the TED struct
func (s *TED) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("TED: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.ApplicationErrorConditionCode = elements[0]
	s.FreeFormMessage = elements[1]

	return nil
}
