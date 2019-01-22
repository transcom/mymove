package edisegment

import (
	"fmt"
	"strconv"
)

// GE represents the GE EDI segment
type GE struct {
	NumberOfTransactionSetsIncluded int
	GroupControlNumber              int64
}

// StringArray converts GE to an array of strings
func (s *GE) StringArray() []string {
	return []string{
		"GE",
		strconv.Itoa(s.NumberOfTransactionSetsIncluded),
		strconv.FormatInt(s.GroupControlNumber, 10),
	}
}

// Parse parses an X12 string that's split into an array into the GE struct
func (s *GE) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("GE: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.NumberOfTransactionSetsIncluded, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}
	s.GroupControlNumber, err = strconv.ParseInt(elements[1], 10, 64)
	return err
}
