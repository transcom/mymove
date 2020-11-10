package edisegment

import (
	"fmt"
)

// BX represents the BX EDI segment
type BX struct {
	TransactionSetPurposeCode    string `validate:"eq=00"`
	TransactionMethodTypeCode    string `validate:"eq=J"`
	ShipmentMethodOfPayment      string `validate:"eq=PP"`
	ShipmentIdentificationNumber string `validate:"min=1,max=30"`
	StandardCarrierAlphaCode     string `validate:"alpha,min=2,max=4"`
	WeightUnitCode               string `validate:"isdefault"` // not used
	ShipmentQualifier            string `validate:"eq=4"`
}

// StringArray converts BX to an array of strings
func (s *BX) StringArray() []string {
	return []string{
		"BX",
		s.TransactionSetPurposeCode,
		s.TransactionMethodTypeCode,
		s.ShipmentMethodOfPayment,
		s.ShipmentIdentificationNumber,
		s.StandardCarrierAlphaCode,
		s.WeightUnitCode,
		s.ShipmentQualifier,
	}
}

// Parse parses an X12 string that's split into an array into the BX struct
func (s *BX) Parse(elements []string) error {
	expectedNumElements := 7
	if len(elements) != expectedNumElements {
		return fmt.Errorf("BX: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.TransactionSetPurposeCode = elements[0]
	s.TransactionSetPurposeCode = elements[1]
	s.ShipmentMethodOfPayment = elements[2]
	s.ShipmentIdentificationNumber = elements[3]
	s.StandardCarrierAlphaCode = elements[4]
	s.WeightUnitCode = elements[5]
	s.ShipmentQualifier = elements[6]
	return nil
}
