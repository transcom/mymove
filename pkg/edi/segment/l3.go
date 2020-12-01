package edisegment

import (
	"fmt"
	"strconv"
)

// L3 represents the L3 EDI segment
type L3 struct {
	Weight          float64 `validate:"required_with=WeightQualifier,min=0,max=9999999999"`
	WeightQualifier string  `validate:"required_with=Weight,omitempty,eq=B"`
	PriceCents      int64   `validate:"min=-999999999999,max=999999999999"` // Supports negative values for FSC price
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
		"", //not used
		"", //not used
		strconv.FormatInt(s.PriceCents, 10),
	}
}

// Parse parses an X12 string that's split into an array into the L3 struct
func (s *L3) Parse(parts []string) error {
	if len(parts) != 5 {
		return fmt.Errorf("L3: Wrong number of elements, expected 5, got %d", len(parts))
	}

	var err error

	s.Weight, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return err
	}
	s.WeightQualifier = parts[1]
	s.PriceCents, err = strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return err
	}

	return nil
}
