package serviceparamvaluelookups

import (
	"fmt"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestRequestedPickupDateLookup() {
	key := models.ServiceItemParamNameRequestedPickupDate

	requestedPickupDate := time.Date(testdatagen.TestYear, time.May, 18, 0, 0, 0, 0, time.UTC)
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				RequestedPickupDate: &requestedPickupDate,
			},
		})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: mtoServiceItem.MoveTaskOrder,
		})

	suite.T().Run("golden path", func(t *testing.T) {
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		paramLookup, err := ServiceParamLookupInitialize(appCfg, suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		expected := requestedPickupDate.Format(ghcrateengine.DateParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.T().Run("nil requested pickup date", func(t *testing.T) {
		// Set the requested pickup date to nil
		mtoShipment := mtoServiceItem.MTOShipment
		oldRequestedPickupDate := mtoShipment.RequestedPickupDate
		mtoShipment.RequestedPickupDate = nil
		suite.MustSave(&mtoShipment)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		paramLookup, err := ServiceParamLookupInitialize(appCfg, suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find a requested pickup date for MTOShipmentID [%s]", mtoShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)

		mtoShipment.RequestedPickupDate = oldRequestedPickupDate
		suite.MustSave(&mtoShipment)
	})
}
