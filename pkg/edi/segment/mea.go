package edisegment

import (
	"fmt"
	"strconv"
)

// MEA represents the MEA EDI segment
type MEA struct {
	MeasurementReferenceIDCode string
	MeasurementQualifier       string
	MeasurementValue           float64
}

// StringArray converts MEA to an array of strings
func (s *MEA) StringArray() []string {
	return []string{
		"MEA",
		s.MeasurementReferenceIDCode,
		s.MeasurementQualifier,
		strconv.FormatFloat(s.MeasurementValue, 'f', 3, 64),
	}
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
