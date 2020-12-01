package edisegment

import (
	"fmt"
	"strconv"
)

// L5 represents the L5 EDI segment
type L5 struct {
	LadingLineItemNumber   int    `validate:"min=1,max=999"`
	LadingDescription      string `validate:"required"`
	CommodityCode          string `validate:"required_with=CommodityCodeQualifier,omitempty,gt=0,lt=11"`
	CommodityCodeQualifier string `validate:"required_with=CommodityCode,omitempty,eq=D"`
}

// StringArray converts L5 to an array of strings
func (s *L5) StringArray() []string {

	return []string{
		"L5",
		strconv.Itoa(s.LadingLineItemNumber),
		s.LadingDescription,
		s.CommodityCode,
		s.CommodityCodeQualifier,
		// TODO: will need to fill in the blank fields if using Marks and Numbers to identify shipments
		"",
		"",
	}
}

// Parse parses an X12 string that's split into an array into the L5 struct
func (s *L5) Parse(parts []string) error {
	numElements := len(parts)
	if numElements != 2 && numElements != 4 && numElements != 6 {
		return fmt.Errorf("L5: Wrong number of elements, expected 4 or 6, got %d", numElements)
	}

	var err error
	s.LadingLineItemNumber, err = strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	s.LadingDescription = parts[1]
	s.CommodityCode = parts[2]
	s.CommodityCodeQualifier = parts[3]
	return nil
}
