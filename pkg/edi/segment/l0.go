package edisegment

import (
	"fmt"
	"strconv"
)

// L0 represents the L0 EDI segment
type L0 struct {
	LadingLineItemNumber   int     `validate:"min=1,max=999"`
	BilledRatedAsQuantity  float64 `validate:"required_with=BilledRatedAsQualifier"`
	BilledRatedAsQualifier string  `validate:"required_with=BilledRatedAsQuantity,omitempty,len=2"`
	Weight                 float64 `validate:"required_with=WeightQualifier"`
	WeightQualifier        string  `validate:"required_with=Weight,omitempty,eq=B"`
	Volume                 float64 `validate:"required_with=VolumeUnitQualifier,omitempty"`
	VolumeUnitQualifier    string  `validate:"required_with=Volume,omitempty,eq=E,len=1"`
	LadingQuantity         int     `validate:"omitempty"`
	WeightUnitCode         string  //`validate:"oneof=L CRT"`
}

// StringArray converts L0 to an array of strings
func (s *L0) StringArray() []string {

	var weight string
	if s.Weight == 0 {
		weight = ""
	} else {
		weight = strconv.FormatFloat(s.Weight, 'f', 3, 64)
	}

	var billedRatedAsQuantity string
	if s.BilledRatedAsQuantity == 0 {
		billedRatedAsQuantity = ""
	} else {
		billedRatedAsQuantity = strconv.FormatFloat(s.BilledRatedAsQuantity, 'f', 3, 64)
	}

	var volume string
	if s.Volume == 0 {
		volume = ""
	} else {
		volume = strconv.FormatFloat(s.Volume, 'f', 3, 64)
	}

	return []string{
		"L0",
		strconv.Itoa(s.LadingLineItemNumber),
		billedRatedAsQuantity,
		s.BilledRatedAsQualifier,
		weight,
		s.WeightQualifier,
		// TODO: will need to fill in the blank fields for crating
		volume,
		s.VolumeUnitQualifier,
		"",
		s.WeightUnitCode, // Packaging Form Code
	}
}

// Parse parses an X12 string that's split into an array into the L0 struct
func (s *L0) Parse(parts []string) error {
	numElements := len(parts)
	if numElements != 3 && numElements != 11 {
		return fmt.Errorf("L0: Wrong number of elements, expected 3 or 11, got %d", numElements)
	}

	var err error
	s.LadingLineItemNumber, err = strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	s.BilledRatedAsQuantity, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err
	}
	s.BilledRatedAsQualifier = parts[2]

	// TODO update to add in other parts
	if numElements == 9 {
		s.Weight, err = strconv.ParseFloat(parts[3], 64)
		if err != nil {
			return err
		}
		s.WeightQualifier = parts[5]
		s.WeightUnitCode = parts[9]
	}

	return nil
}
