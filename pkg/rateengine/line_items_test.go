package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) TestCreateBaseShipmentLineItems() {
	engine := NewRateEngine(suite.DB(), suite.logger, route.NewTestingPlanner(1044))

	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), 1, 1, []int{1}, []models.ShipmentStatus{models.ShipmentStatusINTRANSIT})
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	assertions.FuelEIADieselPrice.EIAPricePerGallonMillicents = 320700
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	// Refetching shipments from database to get all needed eagerly fetched relationships.
	dbShipment, err := models.FetchShipmentByTSP(suite.DB(), tspUser.TransportationServiceProviderID, shipment.ID)
	suite.FatalNoError(err)

	shipmentCost, err := engine.HandleRunOnShipment(*dbShipment)
	suite.FatalNoError(err)

	lineItems, err := CreateBaseShipmentLineItems(suite.DB(), shipmentCost)
	suite.FatalNoError(err)

	// There are 6 Base Shipment line items:
	// origin fee, destination fee, linehaul, pack, unpack, fuel surcharge
	suite.Len(lineItems, 6)

	itemLHS := suite.findLineItem(lineItems, "LHS")
	if itemLHS != nil {
		suite.validateLineItemFields(*itemLHS, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(1044), models.ShipmentLineItemLocationORIGIN, unit.Cents(260858), unit.Millicents(0))
	}

	item135A := suite.findLineItem(lineItems, "135A")
	if item135A != nil {
		suite.validateLineItemFields(*item135A, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationORIGIN, unit.Cents(10230), unit.Millicents(511000))
	}

	item135B := suite.findLineItem(lineItems, "135B")
	if item135B != nil {
		suite.validateLineItemFields(*item135B, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationDESTINATION, unit.Cents(11524), unit.Millicents(576000))
	}

	item105A := suite.findLineItem(lineItems, "105A")
	if item105A != nil {
		suite.validateLineItemFields(*item105A, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationORIGIN, unit.Cents(88625), unit.Millicents(4431000))
	}

	item105C := suite.findLineItem(lineItems, "105C")
	if item105C != nil {
		suite.validateLineItemFields(*item105C, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationDESTINATION, unit.Cents(9305), unit.Millicents(465280))
	}

	item16A := suite.findLineItem(lineItems, "16A")
	if item105C != nil {
		suite.validateLineItemFields(*item16A, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(1044), models.ShipmentLineItemLocationORIGIN, unit.Cents(15651), unit.Millicents(320700))
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

func (suite *RateEngineSuite) validateLineItemFields(lineItem models.ShipmentLineItem, quantity1 unit.BaseQuantity, quantity2 unit.BaseQuantity, location models.ShipmentLineItemLocation, amountCents unit.Cents, appliedRate unit.Millicents) {
	suite.Equal(quantity1, lineItem.Quantity1)
	suite.Equal(quantity2, lineItem.Quantity2)
	suite.Equal(location, lineItem.Location)
	suite.Equal(amountCents, *lineItem.AmountCents)
	suite.Equal(appliedRate, *lineItem.AppliedRate)

	suite.Equal(models.ShipmentLineItemStatusSUBMITTED, lineItem.Status)
}
