package edi

import (
	"fmt"
	"strconv"
	"strings"
)

// L0 represents the L0 EDI segment
type L0 struct {
	LadingLineItemNumber   int
	BilledRatedAsQuantity  float64
	BilledRatedAsQualifier string
	Weight                 float64
	WeightQualifier        string
	WeightUnitCode         string
}

// String converts L0 to its X12 single line string representation
func (s *L0) String(delimiter string) string {
	elements := []string{
		"L0",
		strconv.Itoa(s.LadingLineItemNumber),
		strconv.FormatFloat(s.BilledRatedAsQuantity, 'f', 3, 64),
		s.BilledRatedAsQualifier,
		strconv.FormatFloat(s.Weight, 'f', 3, 64),
		s.WeightQualifier,
		"",
		"",
		"",
		"",
		"",
		s.WeightUnitCode,
	}
	return strings.Join(elements, delimiter) + "\n"
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

	if numElements == 11 {
		s.Weight, err = strconv.ParseFloat(parts[3], 64)
		if err != nil {
			return err
		}
		s.WeightQualifier = parts[4]
		s.WeightUnitCode = parts[10]
	}

	return nil
}
