package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestMTOAvailableToPrimeLookup() {
	key := models.ServiceItemParamNameMTOAvailableToPrimeAt

	availableToPrimeAt := time.Date(testdatagen.TestYear, time.June, 3, 12, 57, 33, 123, time.UTC)
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: mtoServiceItem.MoveTaskOrder,
		})

	paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	suite.T().Run("golden path", func(t *testing.T) {
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := availableToPrimeAt.Format(ghcrateengine.TimestampParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.T().Run("nil AvailableToPrimeAt", func(t *testing.T) {
		// Set the AvailableToPrimeAt to nil
		moveTaskOrder := paymentRequest.MoveTaskOrder
		oldAvailableToPrimeAt := moveTaskOrder.AvailableToPrimeAt
		moveTaskOrder.AvailableToPrimeAt = nil
		suite.MustSave(&moveTaskOrder)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(&services.BadDataError{}, errors.Unwrap(err))
		expected := fmt.Sprintf("Data received from requester is bad: %s: This move task order is not available to prime", services.BadDataCode)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)

		moveTaskOrder.AvailableToPrimeAt = oldAvailableToPrimeAt
		suite.MustSave(&moveTaskOrder)
	})

	suite.T().Run("bogus MoveTaskOrderID", func(t *testing.T) {
		// Pass in a non-existent MoveTaskOrderID
		invalidMoveTaskOrderID := uuid.Must(uuid.NewV4())
		badParamLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, invalidMoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := badParamLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", valueStr)
	})
}
