package edisegment

import (
	"fmt"
	"strconv"
)

// G62 represents the G62 EDI segment
// Requested Pickup Date 68, Requested Pickup Time 5
// Actual Pickup Date 86, Actual Pickup Time 8
type G62 struct {
	DateQualifier int    `validate:"oneof=10 76 86"`
	Date          string `validate:"datetime=20060102"`
	TimeQualifier int    `validate:"omitempty,oneof=5 8"`
	Time          string `validate:"omitempty,required_with=TimeQualifier,datetime=1504"`
}

// StringArray converts G62 to an array of strings
func (s *G62) StringArray() []string {
	if s.Time == "" {
		return []string{
			"G62",
			strconv.Itoa(s.DateQualifier),
			s.Date,
			"",
			"",
		}
	}
	return []string{
		"G62",
		strconv.Itoa(s.DateQualifier),
		s.Date,
		strconv.Itoa(s.TimeQualifier),
		s.Time,
	}
}

// Parse parses an X12 string that's split into an array into the G62 struct
func (s *G62) Parse(elements []string) error {
	expectedNumElements := 4
	if len(elements) != expectedNumElements {
		return fmt.Errorf("G62: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.DateQualifier, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}
	s.Date = elements[1]
	s.TimeQualifier, err = strconv.Atoi(elements[2])
	if err != nil {
		return err
	}
	s.Time = elements[3]
	return nil
}
