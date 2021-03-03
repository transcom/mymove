package edisegment

import (
	"fmt"
	"strconv"
)

// AK1 represents the AK1 EDI segment
type AK1 struct {
	FunctionalIdentifierCode string `validate:"required,eq=SI"`
	GroupControlNumber       int64  `validate:"required,min=1,max=999999999"`
}

// StringArray converts AK1 to an array of strings
func (s *AK1) StringArray() []string {
	return []string{
		"AK1",
		s.FunctionalIdentifierCode,
		strconv.FormatInt(s.GroupControlNumber, 10),
	}
}

// Parse parses an X12 string that's split into an array into the AK1 struct
func (s *AK1) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("AK1: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.FunctionalIdentifierCode = elements[0]
	s.GroupControlNumber, err = strconv.ParseInt(elements[1], 10, 64)
	return err
}
