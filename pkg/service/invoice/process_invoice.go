package invoice

import (
	"os"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/gex"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
)

// ProcessInvoice is a service object to generate/send/record an invoice.
type ProcessInvoice struct {
	DB                    *pop.Connection
	GexSender             gex.SendToGex
	SendProductionInvoice bool
}

// Call processes an invoice by generating the EDI, sending the invoice to GEX, and recording the status.
func (p ProcessInvoice) Call(invoice *models.Invoice, shipment models.Shipment) (*string, *validate.Errors, error) {
	ediString, err := p.generateAndSendInvoiceData(invoice, shipment)
	if err != nil {
		verrs, err := p.updateInvoiceFailed(invoice, validate.NewErrors(), err)
		return ediString, verrs, err
	}

	// Update invoice record as submitted
	verrs, err := UpdateInvoiceSubmitted{DB: p.DB}.Call(invoice, shipment.ShipmentLineItems)
	if err != nil || verrs.HasAny() {
		// Updating as submitted failed, so we need to try to mark it as failed (which could fail too).
		verrs, err := p.updateInvoiceFailed(invoice, verrs, err)
		return ediString, verrs, err
	}

	// If we get here, everything should be good.
	return ediString, validate.NewErrors(), nil
}

func (p ProcessInvoice) generateAndSendInvoiceData(invoice *models.Invoice, shipment models.Shipment) (*string, error) {
	// pass value into generator --> edi string
	invoice858C, err := ediinvoice.Generate858C(shipment, *invoice, p.DB, p.SendProductionInvoice, clock.New())
	if err != nil {
		return nil, err
	}

	// to use for demo visual
	// should this have a flag or be taken out?
	ediWriter := edi.NewWriter(os.Stdout)
	ediWriter.WriteAll(invoice858C.Segments())

	// send edi through gex post api
	transactionName := "placeholder"
	invoice858CString, err := invoice858C.EDIString()
	if err != nil {
		return nil, err
	}

	resp, err := p.GexSender.Call(invoice858CString, transactionName)
	if err != nil {
		return &invoice858CString, err
	}

	if resp != nil && resp.StatusCode != 200 {
		return &invoice858CString, errors.Errorf("Invoice POST request to GEX failed: response status code %d", resp.StatusCode)
	}

	return &invoice858CString, nil
}

func (p ProcessInvoice) updateInvoiceFailed(invoice *models.Invoice, causeVerrs *validate.Errors, cause error) (*validate.Errors, error) {
	// Update invoice record as failed
	invoice.Status = models.InvoiceStatusSUBMISSIONFAILURE
	verrs, err := p.DB.ValidateAndSave(invoice)
	if err != nil || verrs.HasAny() {
		verrs.Append(causeVerrs)
		return verrs, multierror.Append(err, cause)
	}

	return causeVerrs, cause
}
