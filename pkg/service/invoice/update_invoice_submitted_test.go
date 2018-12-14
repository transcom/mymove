package invoice

import (
	"github.com/facebookgo/clock"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *InvoiceServiceSuite) TestUpdateInvoicesCall() {
	shipmentLineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
	suite.db.Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)

	createInvoice := CreateInvoice{
		suite.db,
		clock.NewMock(),
	}
	var invoice models.Invoice
	verrs, err := createInvoice.Call(&invoice, shipmentLineItem.Shipment)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)
	updateInvoicesSubmitted := UpdateInvoiceSubmitted{
		DB: suite.db,
	}
	shipmentLineItems := models.ShipmentLineItems{shipmentLineItem}

	verrs, err = updateInvoicesSubmitted.Call(&invoice, shipmentLineItems)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	suite.Equal(models.InvoiceStatusSUBMITTED, invoice.Status)
	suite.Equal(invoice.ID, *shipmentLineItems[0].InvoiceID)
}
