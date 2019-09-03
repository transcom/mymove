package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

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
