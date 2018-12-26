package invoice

import (
	"fmt"
	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"strings"
	"time"
)

// CreateInvoice is a service object to create new invoices from Shipment
type CreateInvoice struct {
	DB    *pop.Connection
	Clock clock.Clock
}

// Call creates Invoice and updates its ShipmentLineItem associations
func (c CreateInvoice) Call(approver models.OfficeUser, invoice *models.Invoice, shipment models.Shipment) (*validate.Errors, error) {
	invoiceNumber, err := c.createInvoiceNumber(shipment)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create invoice number")
	}

	*invoice = models.Invoice{
		ApproverID:    approver.ID,
		Approver:      approver,
		Status:        models.InvoiceStatusINPROCESS,
		InvoiceNumber: invoiceNumber,
		InvoicedDate:  c.Clock.Now().UTC(),
		ShipmentID:    shipment.ID,
		Shipment:      shipment,
	}
	return c.DB.ValidateAndCreate(invoice)
}

func (c CreateInvoice) createInvoiceNumber(shipment models.Shipment) (string, error) {
	// If we have existing invoices, then get the existing base invoice number and add the appropriate suffix,
	// then go ahead and return it.
	invoices, err := models.FetchInvoicesForShipment(c.DB, shipment.ID)
	if err != nil {
		return "", err
	}
	invoiceCount := len(invoices)
	if invoiceCount > 0 {
		parts := strings.Split(invoices[invoiceCount-1].InvoiceNumber, "-")
		return fmt.Sprintf("%s-%02d", parts[0], invoiceCount), nil
	}

	acceptedOffers, err := shipment.ShipmentOffers.Accepted()
	if err != nil {
		return "", err
	}
	numAcceptedOffers := len(acceptedOffers)
	if numAcceptedOffers == 0 {
		return "", errors.New("No accepted shipment offer found")
	} else if numAcceptedOffers > 1 {
		return "", errors.Errorf("Found %d accepted shipment offers", numAcceptedOffers)
	}
	acceptedOffer := acceptedOffers[0]

	scac := acceptedOffer.TransportationServiceProviderPerformance.TransportationServiceProvider.StandardCarrierAlphaCode

	loc, err := time.LoadLocation(models.InvoiceTimeZone)
	if err != nil {
		return "", err
	}
	year := shipment.CreatedAt.In(loc).Year()

	invoiceNumber, err := models.GenerateBaseInvoiceNumber(c.DB, scac, year)
	if err != nil {
		return "", errors.Wrap(err, "Could not generate invoice number")
	}

	return invoiceNumber, nil
}
