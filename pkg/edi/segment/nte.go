package edisegment

import (
	"fmt"
)

// NTE represents the NTE EDI segment
type NTE struct {
	NoteReferenceCode string
	Description       string
}

// StringArray converts NTE to an array of strings
func (s *NTE) StringArray() []string {
	return []string{"NTE", s.NoteReferenceCode, s.Description}
}

// Parse parses an X12 string that's split into an array into the NTE struct
func (s *NTE) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("NTE: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.NoteReferenceCode = elements[0]
	s.Description = elements[1]
	return nil
}
