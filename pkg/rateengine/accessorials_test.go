package rateengine

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *RateEngineSuite) createShipmentWithServiceArea(assertions testdatagen.Assertions) models.Shipment {
	shipment := testdatagen.MakeShipment(suite.DB(), assertions)

	zip3 := models.Tariff400ngZip3{
		Zip3:          Zip5ToZip3(shipment.PickupAddress.PostalCode),
		BasepointCity: "Saucier",
		State:         "MS",
		ServiceArea:   "428",
		RateArea:      "US48",
		Region:        "11",
	}
	suite.MustSave(&zip3)

	serviceArea := models.Tariff400ngServiceArea{
		Name:               "Gulfport, MS",
		ServiceArea:        "428",
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		ServicesSchedule:   1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.NonPeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}
	suite.MustSave(&serviceArea)

	return shipment
}

func (suite *RateEngineSuite) TestAccessorialsPricingPackCrate() {
	itemCode := "105B"
	rateCents := unit.Cents(2275)
	netWeight := unit.Pound(1000)
	shipment := suite.createShipmentWithServiceArea(testdatagen.Assertions{
		Shipment: models.Shipment{
			BookDate:  &testdatagen.DateInsidePeakRateCycle,
			NetWeight: &netWeight,
		},
	})
	item := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Quantity1: unit.BaseQuantity(50000),
			Shipment:  shipment,
			Status:    models.ShipmentLineItemStatusAPPROVED,
			Location:  models.ShipmentLineItemLocationORIGIN,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			Code:                itemCode,
			RequiresPreApproval: true,
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
		suite.Equal(rateCents.Multiply(5), computedPriceAndRate.Fee)
	}
}

// Iterates through all codes that have pricers and make sure they don't explode with sane values
func (suite *RateEngineSuite) TestAccessorialsSmokeTest() {
	rateCents := unit.Cents(100)
	netWeight := unit.Pound(1000)
	shipment := suite.createShipmentWithServiceArea(testdatagen.Assertions{
		Shipment: models.Shipment{
			BookDate:  &testdatagen.DateInsidePeakRateCycle,
			NetWeight: &netWeight,
		},
	})

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
