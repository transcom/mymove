package edisegment

import (
	"fmt"
	"strconv"
)

// AK3 represents the AK3 EDI segment
type AK3 struct {
	SegmentIDCode                   string `validate:"min=2,max=3"`
	SegmentPositionInTransactionSet int    `validate:"min=1,max=999999"`
	LoopIdentifierCode              string `validate:"omitempty,min=1,max=6"`
	SegmentSyntaxErrorCode          string `validate:"omitempty,min=1,max=3"`
}

// StringArray converts AK3 to an array of strings
func (s *AK3) StringArray() []string {
	return []string{
		"AK3",
		s.SegmentIDCode,
		strconv.Itoa(s.SegmentPositionInTransactionSet),
		s.LoopIdentifierCode,
		s.SegmentSyntaxErrorCode,
	}
}

// Parse parses an X12 string that's split into an array into the AK3 struct
func (s *AK3) Parse(elements []string) error {
	expectedNumElements := 4
	if len(elements) != expectedNumElements {
		return fmt.Errorf("AK3: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.SegmentIDCode = elements[0]
	s.SegmentPositionInTransactionSet, err = strconv.Atoi(elements[1])
	if err != nil {
		return err
	}
	s.LoopIdentifierCode = elements[2]
	s.SegmentSyntaxErrorCode = elements[3]

	return nil
}
