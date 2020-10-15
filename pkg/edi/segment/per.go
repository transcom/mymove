package edisegment

import (
	"fmt"
)

// PER represents the PER EDI segment
type PER struct {
	ContactFunctionCode          string `validate:"required,oneof=CN IC EM"`
	Name                         string `validate:"omitempty,min=1,max=60"`
	CommunicationNumberQualifier string `validate:"omitempty,eq=TE"`
	CommunicationNumber          string `validate:"omitempty,min=1,max=80"`
}

// StringArray converts PER to an array of strings
func (s *PER) StringArray() []string {

	return []string{
		"PER",
		s.ContactFunctionCode,
		s.Name,
		s.CommunicationNumberQualifier,
		s.CommunicationNumber,
	}
}

// Parse parses an X12 string that's split into an array into the PER struct
func (s *PER) Parse(parts []string) error {
	numElements := len(parts)
	if numElements < 1 || numElements > 4 {
		return fmt.Errorf("PER: Wrong number of elements, expected between 1 and 4, got %d", numElements)
	}

	s.ContactFunctionCode = parts[0]
	s.Name = parts[1]
	s.CommunicationNumberQualifier = parts[2]
	s.CommunicationNumber = parts[3]
	return nil
}
