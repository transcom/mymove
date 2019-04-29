package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

// TestCreateShipmentLineItemCode226A tests that 226A line items are created correctly
func (suite *ModelSuite) TestCreateAndEditShipmentLineItemCode226A() {
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

// TestCreateAndEditShipmentLineItemCode125 tests that 125 line items are created and updated correctly
func (suite *ModelSuite) TestCreateAndEditShipmentLineItemCode125() {
	// test create 125A preapproval
	item125A := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: models.Tariff400ngItem{
			Code: "125A",
		},
	})

	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	address := models.Address{
		StreetAddress1: "123 Test St",
		City:           "City",
		State:          "CA",
		PostalCode:     "94087",
	}

	reas := "This is the reason"
	date := time.Now()
	// free form text
	militaryTime := "0000"

	baseParams := models.BaseShipmentLineItemParams{
		Tariff400ngItemID:   item125A.ID,
		Tariff400ngItemCode: item125A.Code,
		Location:            "ORIGIN",
	}
	additionalParams := models.AdditionalShipmentLineItemParams{
		Reason:  &reas,
		Date:    &date,
		Time:    &militaryTime,
		Address: &address,
	}

	shipmentLineItem, verrs, err := shipment.CreateShipmentLineItem(suite.DB(),
		baseParams, additionalParams)

	if suite.noValidationErrors(verrs, err) {
		// flat rate, quantity 1 should be set to 1. 10000 bq
		suite.EqualValues(unit.BaseQuantityFromInt(1), shipmentLineItem.Quantity1)
		suite.EqualValues(item125A.ID.String(), shipmentLineItem.Tariff400ngItem.ID.String())
		suite.EqualValues(reas, *shipmentLineItem.Reason)
		suite.EqualValues(date, *shipmentLineItem.Date)
		suite.EqualValues(militaryTime, *shipmentLineItem.Time)
		suite.NotNil(shipmentLineItem.Address.ID)
	}

	// test update line item

	// fields to update
	updateReas := "updated reason"
	updateDate := time.Now().Add(time.Duration(1000))
	// Free form text
	// J - Juliet - Local Time
	updateMilitaryTime := "2359J"
	updateAddress := models.Address{
		ID:             *shipmentLineItem.AddressID,
		StreetAddress1: "123 Test St, update test",
		City:           "City Update Test",
		State:          "CA, Update Test",
		PostalCode:     "60000",
	}

	additionalParams = models.AdditionalShipmentLineItemParams{
		Reason:  &updateReas,
		Date:    &updateDate,
		Time:    &updateMilitaryTime,
		Address: &updateAddress,
	}

	verrs, err = shipment.UpdateShipmentLineItem(suite.DB(),
		baseParams, additionalParams, shipmentLineItem)

	if suite.noValidationErrors(verrs, err) {
		// flat rate, quantity 1 should be set to 1 in base quantity. 10000 bq.
		suite.EqualValues(unit.BaseQuantityFromInt(1), shipmentLineItem.Quantity1)
		suite.EqualValues(updateReas, *shipmentLineItem.Reason)
		suite.EqualValues(updateDate, *shipmentLineItem.Date)
		suite.EqualValues(updateMilitaryTime, *shipmentLineItem.Time)
		suite.EqualValues(updateAddress.ID, shipmentLineItem.Address.ID)
		suite.EqualValues(updateAddress.StreetAddress1, shipmentLineItem.Address.StreetAddress1)
		suite.EqualValues(updateAddress.City, shipmentLineItem.Address.City)
		suite.EqualValues(updateAddress.State, shipmentLineItem.Address.State)
		suite.EqualValues(updateAddress.PostalCode, shipmentLineItem.Address.PostalCode)
	}
}
