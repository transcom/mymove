package edisegment

import (
	"fmt"
	"strings"
)

// B3 represents the B3 EDI segment
type B3 struct {
	InvoiceNumber                string
	ShipmentIdentificationNumber string
	ShipmentMethodOfPayment      string
	Date                         string
	NetAmountDue                 float64
	CorrectionIndicator          string
	DeliveryDate                 string
	DateTimeQualifier            string
	StandardCarrierAlphaCode     string
}

// String converts B3 to its X12 single line string representation
func (s *B3) String(delimiter string) string {
	elements := []string{
		"B3",
		"",
		s.InvoiceNumber,
		s.ShipmentIdentificationNumber,
		s.ShipmentMethodOfPayment,
		"",
		s.Date,
		FloatToNx(s.NetAmountDue, 2),
		s.CorrectionIndicator,
		s.DeliveryDate,
		s.DateTimeQualifier,
		s.StandardCarrierAlphaCode,
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the B3 struct
func (s *B3) Parse(elements []string) error {
	expectedNumElements := 11
	if len(elements) != expectedNumElements {
		return fmt.Errorf("%s: Wrong number of elements, expected %d, got %d", "B3", expectedNumElements, len(elements))
	}

	var err error
	s.InvoiceNumber = elements[1]
	s.ShipmentIdentificationNumber = elements[2]
	s.ShipmentMethodOfPayment = elements[3]
	s.Date = elements[5]
	s.NetAmountDue, err = NxToFloat(elements[6], 2)
	if err != nil {
		return err
	}
	s.CorrectionIndicator = elements[7]
	s.DeliveryDate = elements[8]
	s.DateTimeQualifier = elements[9]
	s.StandardCarrierAlphaCode = elements[10]

	return nil
}
