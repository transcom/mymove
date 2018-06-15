package edi

import (
	"fmt"
	"strconv"
	"strings"
)

// L7 represents the B3 EDI segment
type L7 struct {
	LadingLineItemNumber int
	TariffNumber         string
	TariffItemNumber     string
	TariffDistance       int
}

// String converts L7 to its X12 single line string representation
func (s *L7) String(delimiter string) string {
	elements := []string{
		"L7",
		strconv.Itoa(s.LadingLineItemNumber),
		"",
		s.TariffNumber,
		"",
		s.TariffItemNumber,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		strconv.Itoa(s.TariffDistance),
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the L7 struct
func (s *L7) Parse(elements []string) error {
	numElements := len(elements)
	if numElements != 5 && numElements != 13 {
		return fmt.Errorf("L7: Wrong number of elements, expected 5 or 13, got %d", numElements)
	}

	var err error
	s.LadingLineItemNumber, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}
	s.TariffNumber = elements[2]
	s.TariffItemNumber = elements[4]
	if numElements == 13 {
		s.TariffDistance, err = strconv.Atoi(elements[12])
	}
	return err
}
