package edisegment

import (
	"fmt"
	"strconv"
)

// L10 represents the B3 EDI segment
type L10 struct {
	Weight          float64 `validate:"required"`
	WeightQualifier string  `validate:"eq=B"`
	WeightUnitCode  string  `validate:"eq=L"`
}

// StringArray converts L10 to an array of strings
func (s *L10) StringArray() []string {
	return []string{
		"L10",
		strconv.FormatFloat(s.Weight, 'f', 3, 64),
		s.WeightQualifier,
		s.WeightUnitCode,
	}
}

// Parse parses an X12 string that's split into an array into the L10 struct
func (s *L10) Parse(elements []string) error {
	expectedNumElements := 3
	if len(elements) != expectedNumElements {
		return fmt.Errorf("L10: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.Weight, err = strconv.ParseFloat(elements[0], 64)
	if err != nil {
		return err
	}
	s.WeightQualifier = elements[1]
	s.WeightUnitCode = elements[2]
	return nil
}
