package edi

import (
	"fmt"
	"strconv"
	"strings"
)

// IEA represents the IEA EDI segment
type IEA struct {
	NumberOfIncludedFunctionalGroups int
	InterchangeControlNumber         int
}

// String converts IEA to its X12 single line string representation
func (s *IEA) String(delimiter string) string {
	elements := []string{
		"IEA",
		strconv.Itoa(s.NumberOfIncludedFunctionalGroups),
		strconv.Itoa(s.InterchangeControlNumber),
	}
	return strings.Join(elements, delimiter) + "\n"
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
	s.InterchangeControlNumber, err = strconv.Atoi(elements[1])
	return err
}
