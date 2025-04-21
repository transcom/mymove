package serviceparamvaluelookups

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

func (suite *ServiceParamValueLookupsSuite) TestMTOEarliestRequestedPickup() {
	key := models.ServiceItemParamNameMTOEarliestRequestedPickup

	earliestRequestedPickup := time.Date(2025, time.April, 21, 0, 0, 0, 452487000, time.Local)
	var mtoServiceItem models.MTOServiceItem
	var paymentRequest models.PaymentRequest
	var paramLookup *ServiceItemParamKeyData

	setupTestData := func() {
		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: models.TimePointer(time.Now()),
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &earliestRequestedPickup,
				},
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model:    mtoServiceItem.MTOShipment,
				LinkOnly: true,
			},
		}, nil)

		var err error
		paramLookup, err = ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
	}

	suite.Run("golden path", func() {
		setupTestData()

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := earliestRequestedPickup.Format(ghcrateengine.TimestampParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.Run("bogus MoveTaskOrderID", func() {
		setupTestData()

		// Pass in a non-existent MoveTaskOrderID
		invalidMoveTaskOrderID := uuid.Must(uuid.NewV4())
		badParamLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, invalidMoveTaskOrderID, nil)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Nil(badParamLookup)
	})
}
