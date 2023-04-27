package serviceparamvaluelookups

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestZipAddressLookup() {
	pickupKey := models.ServiceItemParamNameZipPickupAddress
	destKey := models.ServiceItemParamNameZipDestAddress

	setupTestData := func() (models.PaymentRequest, models.MTOServiceItem) {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				EndDate: time.Now().Add(24 * time.Hour),
			},
		})
		availableDate := time.Date(testdatagen.TestYear, time.May, 1, 0, 0, 0, 0, time.UTC)
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &availableDate,
			},
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})
		return paymentRequest, mtoServiceItem
	}

	suite.Run("zip code for the pickup address is present on MTO Shipment", func() {
		paymentRequest, mtoServiceItem := setupTestData()
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), pickupKey)
		suite.FatalNoError(err)
		suite.Equal(mtoServiceItem.MTOShipment.PickupAddress.PostalCode, valueStr)
	})

	suite.Run("zip code for the destination address is present on MTO Shipment", func() {
		paymentRequest, mtoServiceItem := setupTestData()
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), destKey)
		suite.FatalNoError(err)
		suite.Equal(mtoServiceItem.MTOShipment.DestinationAddress.PostalCode, valueStr)
	})
}
