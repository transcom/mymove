package rateengine

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) createShipmentWithServiceArea() models.Shipment {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusDELIVERED}
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	return shipments[0]
}

func (suite *RateEngineSuite) TestAccessorialsPricingPackCrate() {
	itemCode := "105B"
	rateCents := unit.Cents(2275)
	shipment := suite.createShipmentWithServiceArea()
	q1 := 5
	item := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Quantity1: unit.BaseQuantityFromInt(q1),
			Shipment:  shipment,
			Status:    models.ShipmentLineItemStatusAPPROVED,
			Location:  models.ShipmentLineItemLocationORIGIN,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			Code:                itemCode,
			RequiresPreApproval: true,
			DiscountType:        models.Tariff400ngItemDiscountTypeHHG,
		},
	})

	testdatagen.MakeTariff400ngItemRate(suite.DB(), testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			Code:      itemCode,
			RateCents: rateCents,
		},
	})

	engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)
	computedPriceAndRate, err := engine.ComputeShipmentLineItemCharge(item)

	if suite.NoError(err) {
		discountRate := shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.LinehaulRate
		suite.Equal(discountRate.Apply(rateCents.Multiply(q1)), computedPriceAndRate.Fee)
	}
}

// Iterates through all codes that have pricers and make sure they don't explode with sane values
func (suite *RateEngineSuite) TestAccessorialsSmokeTest() {
	rateCents := unit.Cents(100)
	shipment := suite.createShipmentWithServiceArea()

	for code := range tariff400ngItemPricing {
		item := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Quantity1: unit.BaseQuantityFromInt(1),
				Shipment:  shipment,
				Status:    models.ShipmentLineItemStatusAPPROVED,
				Location:  models.ShipmentLineItemLocationORIGIN,
			},
			Tariff400ngItem: models.Tariff400ngItem{
				Code:                code,
				RequiresPreApproval: true,
			},
		})

		rateCode := code
		if newCode, ok := tariff400ngItemRateMap[code]; ok {
			rateCode = newCode
		}

		testdatagen.MakeTariff400ngItemRate(suite.DB(), testdatagen.Assertions{
			Tariff400ngItemRate: models.Tariff400ngItemRate{
				Code:      rateCode,
				RateCents: rateCents,
			},
		})

		engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)
		_, err := engine.ComputeShipmentLineItemCharge(item)

		// Make sure we don't error
		if !suite.NoError(err) {
			fmt.Printf("Failed while running code %v\n", code)
		}
	}
}

func (suite *RateEngineSuite) TestPricePreapprovalRequestsForShipment() {
	codes := []string{"105B", "120A", "130A"}
	var shipment models.Shipment
	for _, code := range codes {
		item := testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Status: models.ShipmentLineItemStatusAPPROVED,
			},
			Tariff400ngItem: models.Tariff400ngItem{
				Code:                code,
				RequiresPreApproval: true,
			},
		})
		shipment = item.Shipment
	}

	engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)
	pricedItems, err := engine.PricePreapprovalRequestsForShipment(shipment)

	// There should be no error
	if suite.NoError(err) {
		// All items should have a populated amount
		for _, pricedItem := range pricedItems {
			suite.NotNil(pricedItem.AmountCents)
			suite.NotNil(pricedItem.AppliedRate)
		}
	}
}

func (suite *RateEngineSuite) TestPricePreapprovalRequest() {

	item := testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Status: models.ShipmentLineItemStatusSUBMITTED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			Code:                "4A",
			RequiresPreApproval: true,
		},
	})

	engine := NewRateEngine(suite.DB(), suite.logger, suite.planner)
	err := engine.PricePreapprovalRequest(&item)

	if suite.NoError(err) {
		suite.NotNil(item.AmountCents)
		suite.NotNil(item.AppliedRate)
	}
}
