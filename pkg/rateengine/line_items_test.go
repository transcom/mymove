package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) TestCreateBaseShipmentLineItems() {
	engine := NewRateEngine(suite.db, suite.logger, route.NewTestingPlanner(1044))

	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, 1, 1, []int{1}, []models.ShipmentStatus{models.ShipmentStatusINTRANSIT})
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	// Refetching shipments from database to get all needed eagerly fetched relationships.
	dbShipment, err := models.FetchShipmentByTSP(suite.db, tspUser.TransportationServiceProviderID, shipment.ID)
	suite.NoError(err)

	shipmentCost, err := engine.HandleRunOnShipment(*dbShipment)
	suite.NoError(err)

	lineItems, err := CreateBaseShipmentLineItems(suite.db, shipmentCost)
	suite.NoError(err)

	suite.Len(lineItems, 4)

	itemLHS := suite.findLineItem(lineItems, "LHS")
	if itemLHS != nil {
		suite.validateLineItemFields(*itemLHS, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(1044), models.ShipmentLineItemLocationNEITHER, unit.Cents(260858))
	}

	item135A := suite.findLineItem(lineItems, "135A")
	if item135A != nil {
		suite.validateLineItemFields(*item135A, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationORIGIN, unit.Cents(10230))
	}

	item135B := suite.findLineItem(lineItems, "135B")
	if item135B != nil {
		suite.validateLineItemFields(*item135B, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationDESTINATION, unit.Cents(11524))
	}

	item105A := suite.findLineItem(lineItems, "105A")
	if item105A != nil {
		suite.validateLineItemFields(*item105A, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationORIGIN, unit.Cents(97930))
	}
}

func (suite *RateEngineSuite) findLineItem(lineItems []models.ShipmentLineItem, itemCode string) *models.ShipmentLineItem {
	for _, lineItem := range lineItems {
		if itemCode == lineItem.Tariff400ngItem.Code {
			return &lineItem
		}
	}

	suite.T().Errorf("Could not find shipment line item for %s", itemCode)
	return nil
}

func (suite *RateEngineSuite) validateLineItemFields(lineItem models.ShipmentLineItem, quantity1 unit.BaseQuantity, quantity2 unit.BaseQuantity, location models.ShipmentLineItemLocation, amountCents unit.Cents) {
	suite.Equal(quantity1, lineItem.Quantity1)
	suite.Equal(quantity2, lineItem.Quantity2)
	suite.Equal(location, lineItem.Location)
	suite.Equal(amountCents, *lineItem.AmountCents)

	suite.Equal(models.ShipmentLineItemStatusSUBMITTED, lineItem.Status)
}
