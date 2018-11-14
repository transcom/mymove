package edisegment

import (
	"fmt"
	"strconv"
	"strings"
)

// L10 represents the B3 EDI segment
type L10 struct {
	Weight          float64
	WeightQualifier string
	WeightUnitCode  string
}

// String converts L10 to its X12 single line string representation
func (s *L10) String(delimiter string) string {
	elements := []string{
		"L10",
		strconv.FormatFloat(s.Weight, 'f', 3, 64),
		s.WeightQualifier,
		s.WeightUnitCode,
	}
	return strings.Join(elements, delimiter)
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
