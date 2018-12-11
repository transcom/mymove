package invoice

import (
	"fmt"
	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"strings"
	"time"
)

// CreateInvoices is a service object to create new invoices from Shipments
type CreateInvoices struct {
	DB    *pop.Connection
	Clock clock.Clock
}

// Call creates Invoices and updates their ShipmentLineItem associations
func (c CreateInvoices) Call(invoices *models.Invoices, shipments models.Shipments) (*validate.Errors, error) {
	for _, shipment := range shipments {
		invoiceNumber, err := c.createInvoiceNumber(shipment)
		if err != nil {
			return nil, errors.Wrap(err, "Could not create invoice number")
		}

		*invoices = append(*invoices, models.Invoice{
			Status:        models.InvoiceStatusINPROCESS,
			InvoiceNumber: invoiceNumber,
			InvoicedDate:  c.Clock.Now().UTC(),
			ShipmentID:    shipment.ID,
			Shipment:      shipment,
		})
	}
	return c.DB.ValidateAndCreate(invoices)
}

func (c CreateInvoices) createInvoiceNumber(shipment models.Shipment) (string, error) {
	// Assuming invoice numbers are eagerly fetched on shipment.
	if shipment.Invoices == nil {
		return "", errors.New("Invoices is nil")
	}

	// If we have existing invoices, then get the existing base invoice number and add the appropriate suffix,
	// then go ahead and return it.
	invoices := shipment.Invoices
	invoiceCount := len(invoices)
	if invoiceCount > 0 {
		parts := strings.Split(invoices[invoiceCount-1].InvoiceNumber, "-")
		return fmt.Sprintf("%s-%02d", parts[0], invoiceCount), nil
	}

	acceptedOffers := shipment.ShipmentOffers.Accepted()
	numAcceptedOffers := len(acceptedOffers)
	if numAcceptedOffers == 0 {
		return "", errors.New("No accepted shipment offer found")
	} else if numAcceptedOffers > 1 {
		return "", errors.Errorf("Found %d accepted shipment offers", numAcceptedOffers)
	}
	acceptedOffer := acceptedOffers[0]
	if acceptedOffer.TransportationServiceProvider.ID == uuid.Nil {
		return "", errors.New("Accepted shipment offer is missing Transportation Service Provider")
	}

	scac := acceptedOffer.TransportationServiceProvider.StandardCarrierAlphaCode

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
