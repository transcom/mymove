package edisegment

import (
	"fmt"
	"strconv"
)

// SE represents the SE EDI segment
type SE struct {
	NumberOfIncludedSegments    int
	TransactionSetControlNumber string
}

// StringArray converts SE to an array of strings
func (s *SE) StringArray() []string {
	return []string{
		"SE",
		strconv.Itoa(s.NumberOfIncludedSegments),
		s.TransactionSetControlNumber,
	}
}

// Parse parses an X12 string that's split into an array into the SE struct
func (s *SE) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("SE: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.NumberOfIncludedSegments, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}
	s.TransactionSetControlNumber = elements[1]
	return nil
}
