package edisegment

import (
	"fmt"
	"strings"
)

// BX represents the BX EDI segment
type BX struct {
	TransactionSetPurposeCode    string
	TransactionMethodTypeCode    string
	ShipmentMethodOfPayment      string
	ShipmentIdentificationNumber string
	StandardCarrierAlphaCode     string
	WeightUnitCode               string
	ShipmentQualifier            string
}

// String converts BX to its X12 single line string representation
func (s *BX) String(delimiter string) string {
	elements := []string{
		"BX",
		s.TransactionSetPurposeCode,
		s.TransactionMethodTypeCode,
		s.ShipmentMethodOfPayment,
		s.ShipmentIdentificationNumber,
		s.StandardCarrierAlphaCode,
		s.WeightUnitCode,
		s.ShipmentQualifier,
	}
	return strings.Join(elements, delimiter) + "\n"
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
