package edisegment

import (
	"fmt"
	"strconv"
)

// L0 represents the L0 EDI segment
type L0 struct {
	LadingLineItemNumber   int     `validate:"min=1,max=999"`                                               // L001
	BilledRatedAsQuantity  float64 `validate:"required_with=BilledRatedAsQualifier"`                        // L002
	BilledRatedAsQualifier string  `validate:"required_with=BilledRatedAsQuantity,omitempty,len=2"`         // L003
	Weight                 float64 `validate:"required_with=WeightQualifier WeightUnitCode"`                // L004
	WeightQualifier        string  `validate:"required_with=Weight WeightUnitCode,omitempty,eq=B"`          // L005
	Volume                 float64 `validate:"required_with=VolumeUnitQualifier"`                           // L006
	VolumeUnitQualifier    string  `validate:"required_with=Volume,omitempty,eq=E"`                         // L007
	LadingQuantity         int     `validate:"required_with=PackagingFormCode,omitempty,min=1,max=9999999"` // L008
	PackagingFormCode      string  `validate:"required_with=LadingQuantity,omitempty,len=3"`                // L009
	WeightUnitCode         string  `validate:"required_with=Weight WeightQualifier,omitempty,eq=L"`         // L011
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

	var ladingQuantity string
	if s.LadingQuantity == 0 {
		ladingQuantity = ""
	} else {
		ladingQuantity = strconv.Itoa(s.LadingQuantity)
	}

	return []string{
		"L0",
		strconv.Itoa(s.LadingLineItemNumber),
		billedRatedAsQuantity,
		s.BilledRatedAsQualifier,
		weight,
		s.WeightQualifier,
		volume,
		s.VolumeUnitQualifier,
		ladingQuantity,
		s.PackagingFormCode,
		"", // Dunnage Description (not used)
		s.WeightUnitCode,
	}
}

// Parse parses an X12 string that's split into an array into the L0 struct
func (s *L0) Parse(parts []string) error {
	numElements := len(parts)
	if numElements != 3 && numElements != 9 && numElements != 11 {
		return fmt.Errorf("L0: Wrong number of elements, expected 3, 9 or 11, got %d", numElements)
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

	if numElements == 9 {
		s.Weight, err = strconv.ParseFloat(parts[3], 64)
		if err != nil {
			return err
		}
		s.WeightQualifier = parts[4]
		s.Volume, err = strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return err
		}
		s.VolumeUnitQualifier = parts[6]
		s.LadingQuantity, err = strconv.Atoi(parts[7])
		if err != nil {
			return err
		}
		s.PackagingFormCode = parts[8]
	}

	if numElements == 11 {
		s.WeightUnitCode = parts[10]
	}

	return nil
}
