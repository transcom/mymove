package serviceparamvaluelookups

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestZipAddressLookup() {
	pickupKey := models.ServiceItemParamNameZipPickupAddress
	destKey := models.ServiceItemParamNameZipDestAddress

	setupTestData := func() (models.PaymentRequest, models.MTOServiceItem) {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove})

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)
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
