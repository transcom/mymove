package edisegment

import (
	"fmt"
	"strconv"
)

// L1 represents the L1 EDI segment
type L1 struct {
	LadingLineItemNumber int     `validate:"required,min=1,max=999"`
	FreightRate          *int    `validate:"omitempty,min=0"`
	RateValueQualifier   string  `validate:"required_with=FreightRate,omitempty,eq=LB"`
	Charge               float64 `validate:"required"`
}

// StringArray converts L1 to an array of strings
func (s *L1) StringArray() []string {
	freightRate := ""
	if s.FreightRate != nil {
		freightRate = strconv.Itoa(*s.FreightRate)
	}
	return []string{
		"L1",
		strconv.Itoa(s.LadingLineItemNumber),
		freightRate,
		s.RateValueQualifier,
		FloatToNx(s.Charge, 2),
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
	freightRate, err := strconv.Atoi(elements[1])
	if err != nil {
		return err
	}
	s.FreightRate = &freightRate
	s.RateValueQualifier = elements[2]
	s.Charge, err = NxToFloat(elements[3], 2)
	if err != nil {
		return err
	}
	return err
}
