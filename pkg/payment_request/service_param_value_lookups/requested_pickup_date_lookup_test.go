package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestRequestedPickupDateLookup() {
	key := "RequestedPickupDate"

	requestedPickupDate := time.Date(testdatagen.TestYear, time.May, 18, 0, 0, 0, 0, time.UTC)
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				RequestedPickupDate: &requestedPickupDate,
			},
		})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	paramLookup := ServiceParamLookupInitialize(suite.DB(), mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

	suite.T().Run("golden path", func(t *testing.T) {
		dateStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := fmt.Sprintf("%d-%02d-%02d", testdatagen.TestYear, time.May, 18)
		suite.Equal(expected, dateStr)
	})

	suite.T().Run("null requested pickup date", func(t *testing.T) {
		// Set the requested pickup date to null
		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.RequestedPickupDate = nil
		suite.MustSave(&mtoShipment)

		dateStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find a requested pickup date for MTOShipmentID [%s]", mtoShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", dateStr)
	})

	suite.T().Run("null MTOShipmentID", func(t *testing.T) {
		// Set the MTOShipmentID to null
		mtoServiceItem.MTOShipmentID = nil
		suite.MustSave(&mtoServiceItem)

		dateStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", dateStr)
	})

	suite.T().Run("bogus MTOServiceItemID", func(t *testing.T) {
		// Pass in a non-existent MTOServiceItemID
		invalidMTOServiceItemID := uuid.Must(uuid.NewV4())
		badParamLookup := ServiceParamLookupInitialize(suite.DB(), invalidMTOServiceItemID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

		dateStr, err := badParamLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", dateStr)
	})
}
