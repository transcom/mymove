package serviceparamvaluelookups

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestActualPickupDateLookup() {
	key := models.ServiceItemParamNameActualPickupDate

	actualPickupDate := time.Date(testdatagen.TestYear, time.April, 4, 0, 0, 0, 0, time.UTC)

	var mtoServiceItem models.MTOServiceItem

	setupTestData := func() {
		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ActualPickupDate: &actualPickupDate,
				},
			},
		}, nil)

		// Don't need a payment request for this test.
	}

	suite.Run("golden path", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := actualPickupDate.Format(ghcrateengine.DateParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.Run("nil actual pickup date", func() {
		setupTestData()

		// Set the actual pickup date to nil
		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.ActualPickupDate = nil
		suite.MustSave(&mtoShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("", valueStr)
	})
}
