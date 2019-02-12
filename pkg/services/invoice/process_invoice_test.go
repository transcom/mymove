package invoice

import (
	"errors"
	"net/http"
	"testing"

	"github.com/facebookgo/clock"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type TestGexSender struct {
	sendError error
}

func (g *TestGexSender) SendToGex(edi string, transactionName string) (resp *http.Response, err error) {
	return nil, g.sendError
}

func (suite *InvoiceServiceSuite) TestProcessInvoiceCall() {
	shipment := helperDeliveredShipment(suite)
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	gexSenderSuccess := TestGexSender{nil}
	gexSenderFail := TestGexSender{errors.New("test error")}

	processInvoice := ProcessInvoice{
		DB:                    suite.DB(),
		SendProductionInvoice: false,
		ICNSequencer:          sequence.NewDatabaseSequencer(suite.DB(), ediinvoice.ICNSequenceName),
		// GexSender set by each test below.
	}

	suite.T().Run("process invoice fails due to bad GEX response code", func(t *testing.T) {
		invoice := helperCreateInvoiceForShipment(suite, shipment, officeUser)

		processInvoice.GexSender = &gexSenderFail
		_, verrs, err := processInvoice.Call(&invoice, shipment)
		suite.Empty(verrs.Errors)
		suite.Error(err)

		// Make sure the invoice status notes failure.
		helperCheckInvoiceStatus(suite, invoice.ID, models.InvoiceStatusSUBMISSIONFAILURE)

		// Make sure the line items aren't linked to invoice.
		helperCheckLineItemInvoiceID(suite, &shipment.ID, nil)
	})

	suite.T().Run("process invoice fails due to Generate858C failure", func(t *testing.T) {
		invoice := helperCreateInvoiceForShipment(suite, shipment, officeUser)

		// Temporarily make the GBLNumber nil to force an error.
		savedGBLNumber := shipment.GBLNumber
		shipment.GBLNumber = nil

		processInvoice.GexSender = &gexSenderFail
		_, verrs, err := processInvoice.Call(&invoice, shipment)
		suite.Empty(verrs.Errors)
		suite.Error(err)

		// Set gblNumber back for next test.
		shipment.GBLNumber = savedGBLNumber

		// Make sure the invoice status notes failure.
		helperCheckInvoiceStatus(suite, invoice.ID, models.InvoiceStatusSUBMISSIONFAILURE)

		// Make sure the line items aren't linked to invoice.
		helperCheckLineItemInvoiceID(suite, &shipment.ID, nil)

	})

	suite.T().Run("process invoice fails due to invoice number validation failure", func(t *testing.T) {
		invoice := helperCreateInvoiceForShipment(suite, shipment, officeUser)

		// Set the invoice number to incorrect length range to force a validation error.
		invoice.InvoiceNumber = "12345"

		processInvoice.GexSender = &gexSenderSuccess
		_, verrs, err := processInvoice.Call(&invoice, shipment)
		suite.NotEmpty(verrs.Errors)
		// We should have two invoice number errors since we first try to set the invoice to
		// success, then when that fails, we attempt to set it to failure, which fails too.
		suite.Len(verrs.Get("invoice_number"), 2)
		suite.NoError(err)

		// The invoice will be stuck in IN_PROCESS status since we couldn't update it.
		helperCheckInvoiceStatus(suite, invoice.ID, models.InvoiceStatusINPROCESS)

		// Make sure the line items aren't linked to invoice.
		helperCheckLineItemInvoiceID(suite, &shipment.ID, nil)
	})

	suite.T().Run("process invoice succeeds", func(t *testing.T) {
		invoice := helperCreateInvoiceForShipment(suite, shipment, officeUser)

		processInvoice.GexSender = &gexSenderSuccess
		ediString, verrs, err := processInvoice.Call(&invoice, shipment)
		suite.Empty(verrs.Errors)
		suite.NoError(err)
		suite.NotEmpty(ediString)

		// Make sure the invoice status notes success.
		helperCheckInvoiceStatus(suite, invoice.ID, models.InvoiceStatusSUBMITTED)

		// Make sure the line items are linked to invoice.
		helperCheckLineItemInvoiceID(suite, &shipment.ID, &invoice.ID)
	})
}

func helperDeliveredShipment(suite *InvoiceServiceSuite) models.Shipment {
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), 1, 1, []int{1}, []models.ShipmentStatus{models.ShipmentStatusDELIVERED})
	suite.NoError(err)

	shipment := shipments[0]

	amountCents := unit.Cents(12325)
	testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:    shipment,
			Quantity1:   unit.BaseQuantityFromInt(2000),
			AmountCents: &amountCents,
		},
	})

	// Refetching shipment from database to get all needed eagerly fetched relationships (like line items).
	shipment, err = FetchShipmentForInvoice{DB: suite.DB()}.Call(shipment.ID)
	suite.FatalNoError(err)

	return shipment
}

func helperCreateInvoiceForShipment(suite *InvoiceServiceSuite, shipment models.Shipment, approver models.OfficeUser) models.Invoice {
	var invoice models.Invoice
	verrs, err := CreateInvoice{DB: suite.DB(), Clock: clock.NewMock()}.Call(approver, &invoice, shipment)
	suite.Empty(verrs.Errors)
	suite.FatalNoError(err)
	suite.Equal(models.InvoiceStatusINPROCESS, invoice.Status)

	return invoice
}

func helperCheckInvoiceStatus(suite *InvoiceServiceSuite, invoiceID uuid.UUID, expectedStatus models.InvoiceStatus) {
	var invoice models.Invoice
	err := suite.DB().Find(&invoice, invoiceID)
	suite.FatalNoError(err)

	suite.Equal(expectedStatus, invoice.Status)
}

func helperCheckLineItemInvoiceID(suite *InvoiceServiceSuite, shipmentID *uuid.UUID, invoiceID *uuid.UUID) {
	lineItems, err := models.FetchLineItemsByShipmentID(suite.DB(), shipmentID)
	suite.FatalNoError(err)

	for _, lineItem := range lineItems {
		if invoiceID == nil {
			suite.Nil(lineItem.InvoiceID)
		} else {
			suite.NotNil(lineItem.InvoiceID)
			if lineItem.InvoiceID != nil {
				suite.Equal(*invoiceID, *lineItem.InvoiceID)
			}
		}
	}
}
