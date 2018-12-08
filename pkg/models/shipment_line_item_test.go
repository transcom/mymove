package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchLineItem() {
	//Setup
	lineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
	//make more items that don't relate to the first
	testdatagen.MakeDefaultShipmentLineItem(suite.db)
	testdatagen.MakeDefaultShipmentLineItem(suite.db)

	//Do
	accs, err := models.FetchLineItemsByShipmentID(suite.db, &lineItem.ShipmentID)

	//Test
	suite.NoError(err)
	suite.Equal(1, len(accs))
	suite.Equal(lineItem.ID, accs[0].ID)
}

func (suite *ModelSuite) TestFetchApprovedPreapprovalRequestsByShipment() {
	shipment, err := testdatagen.MakeShipmentForPricing(suite.db, testdatagen.Assertions{})
	suite.FatalNoError(err)

	lineItem := testdatagen.MakeCompleteShipmentLineItem(suite.db, testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:   shipment,
			ShipmentID: shipment.ID,
			Status:     models.ShipmentLineItemStatusAPPROVED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	returnedItems, err := models.FetchApprovedPreapprovalRequestsByShipment(suite.db, shipment)

	if suite.NoError(err) && suite.Len(returnedItems, 1) {
		// We should get back the line item
		suite.Equal(lineItem.ID, returnedItems[0].ID)
		// And the Tariff400ngItem association should be populated
		suite.NotEqual(uuid.Nil, returnedItems[0].Tariff400ngItemID)
		suite.Equal(lineItem.Tariff400ngItemID, returnedItems[0].Tariff400ngItemID)
		suite.Equal(returnedItems[0].Tariff400ngItemID, returnedItems[0].Tariff400ngItem.ID)
	}
}
