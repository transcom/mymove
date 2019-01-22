package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchLineItem() {
	//Setup
	lineItem := testdatagen.MakeDefaultShipmentLineItem(suite.DB())
	//make more items that don't relate to the first
	testdatagen.MakeDefaultShipmentLineItem(suite.DB())
	testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	//Do
	accs, err := models.FetchLineItemsByShipmentID(suite.DB(), &lineItem.ShipmentID)

	if suite.NoError(err) {
		//Test
		suite.Equal(1, len(accs))
		suite.Equal(lineItem.ID, accs[0].ID)

		// Test associations
		suite.Equal(lineItem.ShipmentID, accs[0].Shipment.ID)
		suite.Equal(lineItem.Tariff400ngItemID, accs[0].Tariff400ngItem.ID)
	}
}

func (suite *ModelSuite) TestFetchApprovedPreapprovalRequestsByShipment() {
	shipment, err := testdatagen.MakeShipmentForPricing(suite.DB(), testdatagen.Assertions{})
	suite.FatalNoError(err)

	// Given: An approved pre-approval line item
	lineItem := testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:   shipment,
			ShipmentID: shipment.ID,
			Status:     models.ShipmentLineItemStatusAPPROVED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	returnedItems, err := models.FetchApprovedPreapprovalRequestsByShipment(suite.DB(), shipment)

	if suite.NoError(err) && suite.Len(returnedItems, 1) {
		// We should get back the line item
		suite.Equal(lineItem.ID, returnedItems[0].ID)
		// And the Tariff400ngItem association should be populated
		suite.NotEqual(uuid.Nil, returnedItems[0].Tariff400ngItemID)
		suite.Equal(lineItem.Tariff400ngItemID, returnedItems[0].Tariff400ngItemID)
		suite.Equal(returnedItems[0].Tariff400ngItemID, returnedItems[0].Tariff400ngItem.ID)
	}
}

func (suite *ModelSuite) TestFetchShipmentLineItemByID() {
	// Given: A shipment line item
	lineItem := testdatagen.MakeDefaultShipmentLineItem(suite.DB())

	fetchedItem, err := models.FetchShipmentLineItemByID(suite.DB(), &lineItem.ID)

	if suite.NoError(err) {
		suite.Equal(lineItem.ID, fetchedItem.ID)
		suite.Equal(lineItem.Quantity1, fetchedItem.Quantity1)

		suite.Equal(lineItem.ShipmentID, fetchedItem.Shipment.ID)
		suite.Equal(lineItem.Tariff400ngItemID, fetchedItem.Tariff400ngItem.ID)
	}
}

func (suite *ModelSuite) TestApproveShipmentLineItem() {
	// Given: A submitted pre-approval request
	lineItem := testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Status: models.ShipmentLineItemStatusSUBMITTED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	err := lineItem.Approve()

	if suite.NoError(err) {
		suite.Equal(models.ShipmentLineItemStatusAPPROVED, lineItem.Status)
		suite.False(lineItem.ApprovedDate.IsZero())
	}
}

func (suite *ModelSuite) TestApproveShipmentLineItemFails() {
	// Given: An approved pre-approval request
	lineItem := testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Status: models.ShipmentLineItemStatusAPPROVED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	err := lineItem.Approve()

	suite.Error(err)
}

func (suite *ModelSuite) TestDestroyInvoicedShipmentLineItemFails() {
	// Given: An invoice ShipmentLineItem with an invoice ID
	invoice := testdatagen.MakeDefaultInvoice(suite.DB())
	lineItem := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Status:    models.ShipmentLineItemStatusAPPROVED,
			InvoiceID: &invoice.ID,
		},
	})

	// When: The line item is destroyed
	err := suite.DB().Destroy(&lineItem)

	// Then: The destroy action fails
	suite.EqualError(err, models.ErrDestroyForbidden.Error())
}
