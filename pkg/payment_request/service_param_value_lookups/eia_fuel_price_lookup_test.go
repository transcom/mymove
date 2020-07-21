package serviceparamvaluelookups

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestEIAFuelPriceLookup() {
	key := models.ServiceItemParamNameEIAFuelPrice.String()
	actualPickupDate := time.Date(2020, time.July, 21, 0, 0, 0, 0, time.UTC)

	var firstGHCDieselFuelPrice models.GHCDieselFuelPrice
	firstGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 13, 0, 0, 0, 0, time.UTC)
	firstGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243799)
	err := suite.DB().Save(&firstGHCDieselFuelPrice)
	suite.NoError(err)

	var secondGHCDieselFuelPrice models.GHCDieselFuelPrice
	secondGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 20, 0, 0, 0, 0, time.UTC)
	secondGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243299)
	err = suite.DB().Save(&secondGHCDieselFuelPrice)
	suite.NoError(err)

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

	paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

	suite.T().Run("lookup GHC diesel fuel price successfully", func(t *testing.T) {
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("243299", valueStr)
	})
}