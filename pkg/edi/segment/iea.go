package edisegment

import (
	"fmt"
	"strconv"
)

// IEA represents the IEA EDI segment
type IEA struct {
	NumberOfIncludedFunctionalGroups int
	InterchangeControlNumber         int64
}

// StringArray converts IEA to an array of strings
func (s *IEA) StringArray() []string {
	return []string{
		"IEA",
		strconv.Itoa(s.NumberOfIncludedFunctionalGroups),
		fmt.Sprintf("%09d", s.InterchangeControlNumber),
	}
}

// Parse parses an X12 string that's split into an array into the IEA struct
func (s *IEA) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("IEA: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.NumberOfIncludedFunctionalGroups, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}
	s.InterchangeControlNumber, err = strconv.ParseInt(elements[1], 10, 64)
	return err
}
