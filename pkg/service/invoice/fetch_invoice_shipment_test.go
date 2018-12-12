package invoice

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type testValues struct {
	name                string
	requiresPreapproval bool
	approved            models.ShipmentLineItemStatus
	invoiced            bool
	expectedCount       int
}

var tvs = []testValues{
	{
		"preapproval approved",
		true,
		models.ShipmentLineItemStatusAPPROVED,
		false,
		1,
	},
	{
		"preapproval not approved",
		true,
		models.ShipmentLineItemStatusSUBMITTED,
		false,
		0,
	},
	{
		"no preapproval",
		false,
		models.ShipmentLineItemStatusSUBMITTED,
		false,
		1,
	},
	{
		"already invoiced",
		false,
		models.ShipmentLineItemStatusSUBMITTED,
		true,
		0,
	},
}

func (suite *InvoiceServiceSuite) TestFetchInvoiceShipmentCall() {
	for _, tv := range tvs {
		suite.T().Run(tv.name, func(t *testing.T) {
			shipment := testdatagen.MakeDefaultShipment(suite.db)
			lineItem := helperSetupLineItem(shipment, tv, suite.db)
			suite.NotEqual(models.ShipmentLineItem{}.ID, lineItem.ID)

			f := FetchInvoiceShipment{suite.db}
			actualShipment, err := f.Call(shipment.ID)
			suite.NoError(err)

			suite.Equal(tv.expectedCount, len(actualShipment.ShipmentLineItems))
			if tv.expectedCount != 0 {
				suite.Equal(lineItem.ID, actualShipment.ShipmentLineItems[0].ID)
			}
		})
	}

	suite.T().Run("multiple line items", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultShipment(suite.db)
		for _, tv := range tvs {
			helperSetupLineItem(shipment, tv, suite.db)
		}

		f := FetchInvoiceShipment{suite.db}
		actualShipment, err := f.Call(shipment.ID)
		suite.NoError(err)

		suite.Equal(2, len(actualShipment.ShipmentLineItems))
	})
}

func helperSetupLineItem(shipment models.Shipment, tv testValues, db *pop.Connection) models.ShipmentLineItem {
	assertions := testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:   shipment,
			ShipmentID: shipment.ID,
			Status:     tv.approved,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: tv.requiresPreapproval,
		},
	}
	if tv.invoiced {
		invoice := testdatagen.MakeInvoice(db, testdatagen.Assertions{})
		assertions.ShipmentLineItem.InvoiceID = &invoice.ID
	}
	return testdatagen.MakeCompleteShipmentLineItem(db, assertions)
}
