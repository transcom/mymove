package edisegment

import (
	"fmt"
	"strconv"
)

// AK9 represents the AK9 EDI segment
type AK9 struct {
	FunctionalGroupAcknowledgeCode      string `validate:"oneof=A E P R"`
	NumberOfTransactionSetsIncluded     int    `validate:"min=1,max=999999"`
	NumberOfReceivedTransactionSets     int    `validate:"min=1,max=999999"`
	NumberOfAcceptedTransactionSets     int    `validate:"min=1,max=999999"`
	FunctionalGroupSyntaxErrorCodeAK905 string `validate:"omitempty,max=3"`
	FunctionalGroupSyntaxErrorCodeAK906 string `validate:"omitempty,max=3"`
	FunctionalGroupSyntaxErrorCodeAK907 string `validate:"omitempty,max=3"`
	FunctionalGroupSyntaxErrorCodeAK908 string `validate:"omitempty,max=3"`
	FunctionalGroupSyntaxErrorCodeAK909 string `validate:"omitempty,max=3"`
}

// StringArray converts AK9 to an array of strings
func (s *AK9) StringArray() []string {
	return []string{
		"AK9",
		s.FunctionalGroupAcknowledgeCode,
		strconv.Itoa(s.NumberOfTransactionSetsIncluded),
		strconv.Itoa(s.NumberOfReceivedTransactionSets),
		strconv.Itoa(s.NumberOfAcceptedTransactionSets),
		s.FunctionalGroupSyntaxErrorCodeAK905,
		s.FunctionalGroupSyntaxErrorCodeAK906,
		s.FunctionalGroupSyntaxErrorCodeAK907,
		s.FunctionalGroupSyntaxErrorCodeAK908,
		s.FunctionalGroupSyntaxErrorCodeAK909,
	}
}

// Parse parses an X12 string that's split into an array into the AK9 struct
func (s *AK9) Parse(elements []string) error {
	expectedNumElements := 9
	if len(elements) != expectedNumElements {
		return fmt.Errorf("AK9: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.FunctionalGroupAcknowledgeCode = elements[0]

	var err error
	s.NumberOfTransactionSetsIncluded, err = strconv.Atoi(elements[1])
	if err != nil {
		return err
	}

	s.NumberOfReceivedTransactionSets, err = strconv.Atoi(elements[2])
	if err != nil {
		return err
	}

	s.NumberOfAcceptedTransactionSets, err = strconv.Atoi(elements[3])
	if err != nil {
		return err
	}

	s.FunctionalGroupSyntaxErrorCodeAK905 = elements[4]
	s.FunctionalGroupSyntaxErrorCodeAK906 = elements[5]
	s.FunctionalGroupSyntaxErrorCodeAK907 = elements[6]
	s.FunctionalGroupSyntaxErrorCodeAK908 = elements[7]
	s.FunctionalGroupSyntaxErrorCodeAK909 = elements[8]

	return nil
}
