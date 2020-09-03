package invoice

import (
	"github.com/transcom/mymove/pkg/models"

	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

// GHCPaymentRequestInvoiceGenerator is a service object to turn payment requests into 858s
type GHCPaymentRequestInvoiceGenerator struct {
}

// Generate method takes a payment request and returns an Invoice858C
func (g GHCPaymentRequestInvoiceGenerator) Generate(paymentRequest models.PaymentRequest) (ediinvoice.Invoice858C, error) {
	// TODO: probably need to check if the MTO is loaded on the paymentRequest that is passed in, not sure what is more in line with go standards to error out if it's not there or look it up.
	// TODO: seems ReferenceID is a *string but cannot be saved as nil, do we need to validate it's not nil here

	var edi858 ediinvoice.Invoice858C
	bx := edisegment.BX{
		TransactionSetPurposeCode:    "00",
		TransactionMethodTypeCode:    "J",
		ShipmentMethodOfPayment:      "PP",
		ShipmentIdentificationNumber: *paymentRequest.MoveTaskOrder.ReferenceID,
		StandardCarrierAlphaCode:     "TRUS",
		ShipmentQualifier:            "4",
	}
	edi858.Header = append(edi858.Header, &bx)

	paymentRequestNumberSegment := edisegment.N9{
		ReferenceIdentificationQualifier: "CN",
		ReferenceIdentification:          paymentRequest.PaymentRequestNumber,
	}
	edi858.Header = append(edi858.Header, &paymentRequestNumberSegment)

	return edi858, nil
}
