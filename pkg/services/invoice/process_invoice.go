package invoice

import (
	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ProcessInvoice is a service object to generate/send/record an invoice.
type ProcessInvoice struct {
	DB                    *pop.Connection
	GexSender             services.GexSender
	SendProductionInvoice bool
	ICNSequencer          sequence.Sequencer
}

// Call processes an invoice by generating the EDI, sending the invoice to GEX, and recording the status.
func (p ProcessInvoice) Call(invoice *models.Invoice, shipment models.Shipment) (*string, *validate.Errors, error) {
	ediString, err := p.generateAndSendInvoiceData(invoice, shipment)
	if err != nil {
		// The invoice submission has failed, so we record the failure.
		verrs, err := p.updateInvoiceFailed(invoice, models.InvoiceStatusSUBMISSIONFAILURE, validate.NewErrors(), err)
		return ediString, verrs, err
	}

	// Update invoice record as submitted
	verrs, err := UpdateInvoiceSubmitted{DB: p.DB}.Call(invoice, shipment.ShipmentLineItems)
	if err != nil || verrs.HasAny() {
		// Updating as submitted failed (although the invoice submission succeeded), so we try to mark it as a
		// status update failure to prevent the invoice from being submitted again.
		verrs, err := p.updateInvoiceFailed(invoice, models.InvoiceStatusUPDATEFAILURE, verrs, err)
		return ediString, verrs, err
	}

	// If we get here, everything should be good.
	return ediString, validate.NewErrors(), nil
}

func (p ProcessInvoice) generateAndSendInvoiceData(invoice *models.Invoice, shipment models.Shipment) (*string, error) {
	// pass value into generator --> edi string
	invoice858C, err := ediinvoice.Generate858C(shipment, *invoice, p.DB, p.SendProductionInvoice, p.ICNSequencer, clock.New())
	if err != nil {
		return nil, err
	}

	// send edi through gex post api
	transactionName := "placeholder"
	invoice858CString, err := invoice858C.EDIString()
	if err != nil {
		return nil, err
	}

	resp, err := p.GexSender.SendToGex(invoice858CString, transactionName)
	if err != nil {
		return &invoice858CString, err
	}

	if resp != nil && resp.StatusCode != 200 {
		return &invoice858CString, errors.Errorf("Invoice POST request to GEX failed: response status code %d", resp.StatusCode)
	}

	return &invoice858CString, nil
}

func (p ProcessInvoice) updateInvoiceFailed(invoice *models.Invoice, invoiceStatus models.InvoiceStatus, causeVerrs *validate.Errors, cause error) (*validate.Errors, error) {
	// Update invoice record as failed
	invoice.Status = invoiceStatus
	verrs, err := p.DB.ValidateAndSave(invoice)
	if err != nil || verrs.HasAny() {
		verrs.Append(causeVerrs)
		if err != nil {
			if cause != nil {
				return verrs, multierror.Append(err, cause)
			}
			return verrs, err
		}
		return verrs, cause
	}

	return causeVerrs, cause
}
