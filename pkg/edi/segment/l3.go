package edisegment

import (
	"fmt"
	"strconv"
)

// L3 represents the L3 EDI segment
type L3 struct {
	Weight          float64 `validate:"required_with=WeightQualifier"`
	WeightQualifier string  `validate:"required_with=Weight"`
	PriceCents      int
}

// StringArray converts L3 to an array of strings
func (s *L3) StringArray() []string {

	var weight string
	if s.Weight == 0 {
		weight = ""
	} else {
		weight = strconv.FormatFloat(s.Weight, 'f', 3, 64)
	}

	return []string{
		"L3",
		weight,
		s.WeightQualifier,
		strconv.Itoa(s.PriceCents),
	}
}

// Parse parses an X12 string that's split into an array into the L3 struct
func (s *L3) Parse(parts []string) error {
	if len(parts) != 3 {
		return fmt.Errorf("L3: Wrong number of elements, expected 3, got %d", len(parts))
	}

	var err error

	s.Weight, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return err
	}
	s.WeightQualifier = parts[1]
	s.PriceCents, err = strconv.Atoi(parts[2])
	if err != nil {
		return err
	}

	return nil
}
