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

func (suite *ServiceParamValueLookupsSuite) TestActualPickupDateLookup() {
	key := models.ServiceItemParamNameActualPickupDate.String()

	actualPickupDate := time.Date(testdatagen.TestYear, time.May, 18, 0, 0, 0, 0, time.UTC)
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ActualPickupDate: &actualPickupDate,
			},
		})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: mtoServiceItem.MoveTaskOrder,
		})

	paramLookup, _ := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

	suite.T().Run("golden path", func(t *testing.T) {
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := actualPickupDate.Format(ghcrateengine.DateParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.T().Run("nil actual pickup date", func(t *testing.T) {
		// Set the actual pickup date to nil
		mtoShipment := mtoServiceItem.MTOShipment
		oldActualPickupDate := mtoShipment.ActualPickupDate
		mtoShipment.ActualPickupDate = nil
		suite.MustSave(&mtoShipment)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find an actual pickup date for MTOShipmentID [%s]", mtoShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)

		mtoShipment.ActualPickupDate = oldActualPickupDate
		suite.MustSave(&mtoShipment)
	})

	suite.T().Run("nil MTOShipmentID", func(t *testing.T) {
		// Set the MTOShipmentID to nil
		oldMTOShipmentID := mtoServiceItem.MTOShipmentID
		mtoServiceItem.MTOShipmentID = nil
		suite.MustSave(&mtoServiceItem)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", valueStr)

		mtoServiceItem.MTOShipmentID = oldMTOShipmentID
		suite.MustSave(&mtoServiceItem)
	})

	suite.T().Run("bogus MTOServiceItemID", func(t *testing.T) {
		// Pass in a non-existent MTOServiceItemID
		invalidMTOServiceItemID := uuid.Must(uuid.NewV4())
		_, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, invalidMTOServiceItemID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})
}
