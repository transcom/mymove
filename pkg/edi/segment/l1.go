package edisegment

import (
	"fmt"
	"strconv"
)

// L1 represents the L1 EDI segment
type L1 struct {
	LadingLineItemNumber     int
	FreightRate              float64
	RateValueQualifier       string
	Charge                   float64
	SpecialChargeDescription string
}

// StringArray converts L1 to an array of strings
func (s *L1) StringArray() []string {
	return []string{
		"L1",
		strconv.Itoa(s.LadingLineItemNumber),
		strconv.FormatFloat(s.FreightRate, 'f', 4, 64),
		s.RateValueQualifier,
		FloatToNx(s.Charge, 2),
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		s.SpecialChargeDescription,
	}
}

// Parse parses an X12 string that's split into an array into the L1 struct
func (s *L1) Parse(elements []string) error {
	expectedNumElements := 12
	if len(elements) != expectedNumElements {
		return fmt.Errorf("L1: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.LadingLineItemNumber, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}
	s.FreightRate, err = strconv.ParseFloat(elements[1], 64)
	if err != nil {
		return err
	}
	s.RateValueQualifier = elements[2]
	s.Charge, err = NxToFloat(elements[3], 2)
	if err != nil {
		return err
	}
	s.SpecialChargeDescription = elements[11]
	return err
}
