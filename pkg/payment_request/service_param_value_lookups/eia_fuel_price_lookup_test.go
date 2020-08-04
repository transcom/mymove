package serviceparamvaluelookups

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestEIAFuelPriceLookup() {
	key := models.ServiceItemParamNameEIAFuelPrice.String()
	actualPickupDate := time.Date(2020, time.July, 15, 0, 0, 0, 0, time.UTC)

	var firstGHCDieselFuelPrice models.GHCDieselFuelPrice
	firstGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 06, 0, 0, 0, 0, time.UTC)
	firstGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243699)
	_ = suite.DB().Save(&firstGHCDieselFuelPrice)

	var secondGHCDieselFuelPrice models.GHCDieselFuelPrice
	secondGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 13, 0, 0, 0, 0, time.UTC)
	secondGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243799)
	_ = suite.DB().Save(&secondGHCDieselFuelPrice)

	var thirdGHCDieselFuelPrice models.GHCDieselFuelPrice
	thirdGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 20, 0, 0, 0, 0, time.UTC)
	thirdGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243299)
	_ = suite.DB().Save(&thirdGHCDieselFuelPrice)

	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ActualPickupDate: &actualPickupDate,
		},
	})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	suite.T().Run("lookup GHC diesel fuel price successfully", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("243799", valueStr)
	})

	suite.T().Run("No MTO shipment found", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})
		mtoServiceItem.MTOShipmentID = nil
		err := suite.DB().Save(&mtoServiceItem)
		suite.NoError(err)

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("No MTO shipment pickup date found", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})
}
