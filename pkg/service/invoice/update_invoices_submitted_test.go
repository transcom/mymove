package invoice

import (
	"github.com/facebookgo/clock"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *InvoiceServiceSuite) TestUpdateInvoicesCall() {
	shipmentLineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
	suite.db.Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)

	createInvoices := CreateInvoices{
		suite.db,
		clock.NewMock(),
	}
	var invoices models.Invoices
	verrs, err := createInvoices.Call(&invoices, models.Shipments{shipmentLineItem.Shipment})
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)
	updateInvoicesSubmitted := UpdateInvoicesSubmitted{
		DB: suite.db,
	}
	shipmentLineItems := models.ShipmentLineItems{shipmentLineItem}

	verrs, err = updateInvoicesSubmitted.Call(invoices, shipmentLineItems)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	suite.Equal(1, len(invoices))
	suite.Equal(models.InvoiceStatusSUBMITTED, invoices[0].Status)
	suite.Equal(invoices[0].ID, *shipmentLineItems[0].InvoiceID)
}
