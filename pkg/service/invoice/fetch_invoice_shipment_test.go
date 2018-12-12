package invoice

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

var approvalTests = []struct {
	name                string
	requiresPreapproval bool
	approved            models.ShipmentLineItemStatus
	expectedCount       int
}{
	{"preapproval approved", true, models.ShipmentLineItemStatusAPPROVED, 1},
	{"preapproval not approved", true, models.ShipmentLineItemStatusSUBMITTED, 0},
	{"no preapproval", false, models.ShipmentLineItemStatusSUBMITTED, 1},
}

func (suite *InvoiceServiceSuite) TestFetchInvoiceShipmentCall() {
	for _, at := range approvalTests {
		suite.T().Run(at.name, func(t *testing.T) {
			shipment := testdatagen.MakeDefaultShipment(suite.db)
			lineItem := testdatagen.MakeCompleteShipmentLineItem(suite.db, testdatagen.Assertions{
				ShipmentLineItem: models.ShipmentLineItem{
					Shipment:   shipment,
					ShipmentID: shipment.ID,
					Status:     at.approved,
				},
				Tariff400ngItem: models.Tariff400ngItem{
					RequiresPreApproval: at.requiresPreapproval,
				},
			})
			suite.NotEqual(models.ShipmentLineItem{}.ID, lineItem.ID)

			f := FetchInvoiceShipment{suite.db}
			actualShipment, err := f.Call(shipment.ID)
			suite.NoError(err)

			suite.Equal(at.expectedCount, len(actualShipment.ShipmentLineItems))
			if at.expectedCount != 0 {
				suite.Equal(lineItem.ID, actualShipment.ShipmentLineItems[0].ID)
			}
		})
	}

	suite.T().Run("multiple line items", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultShipment(suite.db)
		for _, at := range approvalTests {
			testdatagen.MakeCompleteShipmentLineItem(suite.db, testdatagen.Assertions{
				ShipmentLineItem: models.ShipmentLineItem{
					Shipment:   shipment,
					ShipmentID: shipment.ID,
					Status:     at.approved,
				},
				Tariff400ngItem: models.Tariff400ngItem{
					RequiresPreApproval: at.requiresPreapproval,
				},
			})
		}

		f := FetchInvoiceShipment{suite.db}
		actualShipment, err := f.Call(shipment.ID)
		suite.NoError(err)

		suite.Equal(2, len(actualShipment.ShipmentLineItems))
	})
}
