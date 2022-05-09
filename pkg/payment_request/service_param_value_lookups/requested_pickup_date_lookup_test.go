package serviceparamvaluelookups

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestRequestedPickupDateLookup() {
	key := models.ServiceItemParamNameRequestedPickupDate

	requestedPickupDate := time.Date(testdatagen.TestYear, time.May, 18, 0, 0, 0, 0, time.UTC)

	var mtoServiceItem models.MTOServiceItem

	setupTestData := func() {
		mtoServiceItem = testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate,
				},
			})

		// Don't need a payment request for this test.
	}

	suite.Run("golden path", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := requestedPickupDate.Format(ghcrateengine.DateParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.Run("nil requested pickup date", func() {
		setupTestData()

		// Set the requested pickup date to nil
		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.RequestedPickupDate = nil
		suite.MustSave(&mtoShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("", valueStr)
	})
}
