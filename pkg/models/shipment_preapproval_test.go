package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

// TestCreateShipmentLineItemCode226A tests that 226A line items are created correctly
func (suite *ModelSuite) TestCreateShipmentLineItemCode226A() {
	// test create 226A preapproval
	item226A := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: models.Tariff400ngItem{
			Code: "226A",
		},
	})

	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	desc := "This is a description"
	reas := "This is the reason"
	actAmt := unit.Cents(1000)
	baseParams := models.BaseShipmentLineItemParams{
		Tariff400ngItemID:   item226A.ID,
		Tariff400ngItemCode: item226A.Code,
		Location:            "ORIGIN",
	}
	additionalParams := models.AdditionalShipmentLineItemParams{
		Description:       &desc,
		Reason:            &reas,
		ActualAmountCents: &actAmt,
	}

	shipmentLineItem, verrs, err := shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantityFromCents(actAmt), shipmentLineItem.Quantity1)
		suite.Equal(item226A.ID.String(), shipmentLineItem.Tariff400ngItem.ID.String())
		suite.Equal(desc, *shipmentLineItem.Description)
		suite.Equal(reas, *shipmentLineItem.Reason)
		suite.Equal(actAmt, *shipmentLineItem.ActualAmountCents)
	}

	// test update 226A values
	baseParamsUpdated := models.BaseShipmentLineItemParams{
		Tariff400ngItemID:   item226A.ID,
		Tariff400ngItemCode: item226A.Code,
		Location:            "ORIGIN",
	}
	descUpdated := "updated description"
	reasUpdated := "updated reason"
	actAmtUpdated := unit.Cents(1500)
	additionalParamsUpdated := models.AdditionalShipmentLineItemParams{
		Description:       &descUpdated,
		Reason:            &reasUpdated,
		ActualAmountCents: &actAmtUpdated,
	}

	verrs, err = shipment.UpdateShipmentLineItem(suite.DB(),
		baseParamsUpdated, additionalParamsUpdated, shipmentLineItem)
	if suite.noValidationErrors(verrs, err) {
		suite.Equal(unit.BaseQuantity(150000), shipmentLineItem.Quantity1)
		suite.Equal(item226A.ID.String(), shipmentLineItem.Tariff400ngItem.ID.String())
		suite.Equal(descUpdated, *shipmentLineItem.Description)
		suite.Equal(reasUpdated, *shipmentLineItem.Reason)
		suite.Equal(actAmtUpdated, *shipmentLineItem.ActualAmountCents)
	}

}
