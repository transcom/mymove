package edi

import (
	"fmt"
	"strconv"
	"strings"
)

// MEA represents the MEA EDI segment
type MEA struct {
	MeasurementReferenceIDCode string
	MeasurementQualifier       string
	MeasurementValue           float64
}

// String converts MEA to its X12 single line string representation
func (s *MEA) String(delimiter string) string {
	elements := []string{
		"MEA",
		s.MeasurementReferenceIDCode,
		s.MeasurementQualifier,
		strconv.FormatFloat(s.MeasurementValue, 'f', 3, 64),
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the MEA struct
func (s *MEA) Parse(elements []string) error {
	expectedNumElements := 3
	if len(elements) != expectedNumElements {
		return fmt.Errorf("MEA: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.MeasurementReferenceIDCode = elements[0]
	s.MeasurementQualifier = elements[1]
	s.MeasurementValue, err = strconv.ParseFloat(elements[2], 64)
	return err
}
