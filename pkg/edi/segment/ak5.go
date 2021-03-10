package edisegment

import (
	"fmt"
)

// AK5 represents the AK5 EDI segment
// Purpose: To acknowledge acceptance or rejection and report errors in a transaction set (direct language
// from the 997 spec, though we don't expect to process any errors, only expect to get acks here)
type AK5 struct {
	TransactionSetAcknowledgmentCode   string `validate:"min=1,max=1"`     // only expect to use this field
	TransactionSetSyntaxErrorCodeAK502 string `validate:"omitempty,max=3"` // not expecting these fields to be set
	TransactionSetSyntaxErrorCodeAK503 string `validate:"omitempty,max=3"` // not expecting these fields to be set
	TransactionSetSyntaxErrorCodeAK504 string `validate:"omitempty,max=3"` // not expecting these fields to be set
	TransactionSetSyntaxErrorCodeAK505 string `validate:"omitempty,max=3"` // not expecting these fields to be set
	TransactionSetSyntaxErrorCodeAK506 string `validate:"omitempty,max=3"` // not expecting these fields to be set
}

// StringArray converts AK5 to an array of strings
func (s *AK5) StringArray() []string {
	return []string{
		"AK5",
		s.TransactionSetAcknowledgmentCode,
		s.TransactionSetSyntaxErrorCodeAK502,
		s.TransactionSetSyntaxErrorCodeAK503,
		s.TransactionSetSyntaxErrorCodeAK504,
		s.TransactionSetSyntaxErrorCodeAK505,
		s.TransactionSetSyntaxErrorCodeAK506,
	}
}

// Parse parses an X12 string that's split into an array into the AK5 struct
func (s *AK5) Parse(elements []string) error {
	expectedNumElements := 6
	if len(elements) != expectedNumElements {
		return fmt.Errorf("AK5: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.TransactionSetAcknowledgmentCode = elements[0]
	s.TransactionSetSyntaxErrorCodeAK502 = elements[1]
	s.TransactionSetSyntaxErrorCodeAK503 = elements[2]
	s.TransactionSetSyntaxErrorCodeAK504 = elements[3]
	s.TransactionSetSyntaxErrorCodeAK505 = elements[4]
	s.TransactionSetSyntaxErrorCodeAK506 = elements[5]

	return nil
}
