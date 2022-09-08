package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestMTOAvailableToPrimeLookup() {
	key := models.ServiceItemParamNameMTOAvailableToPrimeAt

	availableToPrimeAt := time.Date(testdatagen.TestYear, time.June, 3, 12, 57, 33, 123, time.UTC)
	var mtoServiceItem models.MTOServiceItem
	var paymentRequest models.PaymentRequest
	var paramLookup *ServiceItemParamKeyData

	setupTestData := func() {
		mtoServiceItem = testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				Move: models.Move{
					AvailableToPrimeAt: &availableToPrimeAt,
				},
			})

		paymentRequest = testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: mtoServiceItem.MoveTaskOrder,
			})

		var err error
		paramLookup, err = ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
	}

	suite.Run("golden path", func() {
		setupTestData()

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := availableToPrimeAt.Format(ghcrateengine.TimestampParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.Run("nil AvailableToPrimeAt", func() {
		setupTestData()

		// Set the AvailableToPrimeAt to nil
		moveTaskOrder := paymentRequest.MoveTaskOrder
		oldAvailableToPrimeAt := moveTaskOrder.AvailableToPrimeAt
		moveTaskOrder.AvailableToPrimeAt = nil
		suite.MustSave(&moveTaskOrder)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.IsType(&apperror.BadDataError{}, errors.Unwrap(err))
		expected := fmt.Sprintf("Data received from requester is bad: %s: This move task order is not available to prime", apperror.BadDataCode)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)

		moveTaskOrder.AvailableToPrimeAt = oldAvailableToPrimeAt
		suite.MustSave(&moveTaskOrder)
	})

	suite.Run("bogus MoveTaskOrderID", func() {
		setupTestData()

		// Pass in a non-existent MoveTaskOrderID
		invalidMoveTaskOrderID := uuid.Must(uuid.NewV4())
		badParamLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, invalidMoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := badParamLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", valueStr)
	})
}
