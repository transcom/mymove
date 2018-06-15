package edi

import (
	"fmt"
	"strconv"
	"strings"
)

// L5 represents the L5 EDI segment
type L5 struct {
	LadingLineItemNumber int
	LadingDescription    string
}

// String converts L5 to its X12 single line string representation
func (s *L5) String(delimiter string) string {
	elements := []string{
		"L5",
		strconv.Itoa(s.LadingLineItemNumber),
		s.LadingDescription,
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the L5 struct
func (s *L5) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("L5: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.LadingLineItemNumber, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}
	s.LadingDescription = elements[1]
	return nil
}
