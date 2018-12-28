package invoice

import (
	"testing"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *InvoiceServiceSuite) TestUpdateInvoicesCall() {
	suite.T().Run("invoice updates", func(t *testing.T) {
		shipmentLineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
		suite.db.Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)
		invoice := helperCreateInvoice(suite, shipmentLineItem.Shipment)

		updateInvoicesSubmitted := UpdateInvoiceSubmitted{
			DB: suite.db,
		}
		shipmentLineItems := models.ShipmentLineItems{shipmentLineItem}

		verrs, err := updateInvoicesSubmitted.Call(invoice, shipmentLineItems)
		suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
		suite.NoError(err)

		suite.Equal(models.InvoiceStatusSUBMITTED, invoice.Status)
		suite.Equal(invoice.ID, *shipmentLineItems[0].InvoiceID)
	})

	suite.T().Run("error when save fails", func(t *testing.T) {
		shipmentLineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
		suite.db.Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)
		invoice := helperCreateInvoice(suite, shipmentLineItem.Shipment)

		fakeUUID, err := uuid.NewV4()
		shipmentLineItem.ShipmentID = fakeUUID // create foreign key constraint error
		invoice.ShipmentID = fakeUUID

		updateInvoicesSubmitted := UpdateInvoiceSubmitted{
			DB: suite.db,
		}
		verrs, err := updateInvoicesSubmitted.Call(invoice, models.ShipmentLineItems{})

		suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
		suite.Error(err)
	})

	suite.T().Run("transaction rolls back", func(t *testing.T) {
		shipmentLineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
		suite.db.Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)
		invoice := helperCreateInvoice(suite, shipmentLineItem.Shipment)

		updateInvoicesSubmitted := UpdateInvoiceSubmitted{
			DB: suite.db,
		}
		fakeUUID, err := uuid.NewV4()
		shipmentLineItem.ShipmentID = fakeUUID // create foreign key constraint error
		shipmentLineItems := models.ShipmentLineItems{shipmentLineItem}

		verrs, err := updateInvoicesSubmitted.Call(invoice, shipmentLineItems)
		suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
		suite.Error(err)

		suite.db.Reload(invoice)
		suite.Equal(models.InvoiceStatusINPROCESS, invoice.Status)
		suite.db.Reload(&shipmentLineItem)
		suite.Nil(shipmentLineItem.InvoiceID)
	})
}

func helperCreateInvoice(suite *InvoiceServiceSuite, shipment models.Shipment) *models.Invoice {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	createInvoice := CreateInvoice{
		suite.db,
		clock.NewMock(),
	}
	var invoice models.Invoice
	verrs, err := createInvoice.Call(officeUser, &invoice, shipment)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	return &invoice
}
