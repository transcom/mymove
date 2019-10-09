package edisegment

import (
	"fmt"
	"strconv"
)

// L7 represents the B3 EDI segment
type L7 struct {
	LadingLineItemNumber int    `validate:"omitempty,min=1,max=999"`
	TariffNumber         string `validate:"omitempty,min=1,max=7"`
	TariffItemNumber     string `validate:"omitempty,min=1,max=16"`
	TariffDistance       int    `validate:"omitempty,min=1,max=99999"`
}

// StringArray converts L7 to an array of strings
func (s *L7) StringArray() []string {
	return []string{
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
