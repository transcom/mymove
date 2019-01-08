package invoice

import (
	"testing"

	"github.com/facebookgo/clock"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *InvoiceServiceSuite) TestUpdateInvoicesCall() {
	suite.T().Run("invoice updates", func(t *testing.T) {
		shipment := helperShipment(suite)
		shipmentLineItem := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Shipment: shipment,
			},
		})

		suite.DB().Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)
		invoice := helperCreateInvoice(suite, shipmentLineItem.Shipment)

		updateInvoicesSubmitted := UpdateInvoiceSubmitted{
			DB: suite.DB(),
		}
		shipmentLineItems := models.ShipmentLineItems{shipmentLineItem}

		verrs, err := updateInvoicesSubmitted.Call(invoice, shipmentLineItems)
		suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
		suite.NoError(err)

		suite.Equal(models.InvoiceStatusSUBMITTED, invoice.Status)
		suite.Equal(invoice.ID, *shipmentLineItems[0].InvoiceID)
	})

	suite.T().Run("error when save fails", func(t *testing.T) {
		shipment := helperShipment(suite)
		shipmentLineItem := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Shipment: shipment,
			},
		})

		suite.DB().Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)
		invoice := helperCreateInvoice(suite, shipmentLineItem.Shipment)

		fakeUUID, err := uuid.NewV4()
		shipmentLineItem.ShipmentID = fakeUUID // create foreign key constraint error
		invoice.ShipmentID = fakeUUID

		updateInvoicesSubmitted := UpdateInvoiceSubmitted{
			DB: suite.DB(),
		}
		verrs, err := updateInvoicesSubmitted.Call(invoice, models.ShipmentLineItems{})

		suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
		suite.Error(err)
	})

	suite.T().Run("transaction rolls back", func(t *testing.T) {
		shipment := helperShipment(suite)
		shipmentLineItem := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Shipment: shipment,
			},
		})

		suite.DB().Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)
		invoice := helperCreateInvoice(suite, shipmentLineItem.Shipment)

		updateInvoicesSubmitted := UpdateInvoiceSubmitted{
			DB: suite.DB(),
		}
		fakeUUID, err := uuid.NewV4()
		shipmentLineItem.ShipmentID = fakeUUID // create foreign key constraint error
		shipmentLineItems := models.ShipmentLineItems{shipmentLineItem}

		verrs, err := updateInvoicesSubmitted.Call(invoice, shipmentLineItems)
		suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
		suite.Error(err)

		suite.DB().Reload(invoice)
		suite.Equal(models.InvoiceStatusINPROCESS, invoice.Status)
		suite.DB().Reload(&shipmentLineItem)
		suite.Nil(shipmentLineItem.InvoiceID)
	})
}

func helperCreateInvoice(suite *InvoiceServiceSuite, shipment models.Shipment) *models.Invoice {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	createInvoice := CreateInvoice{
		suite.DB(),
		clock.NewMock(),
	}
	var invoice models.Invoice
	verrs, err := createInvoice.Call(officeUser, &invoice, shipment)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	return &invoice
}
