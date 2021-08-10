package serviceparamvaluelookups

import (
	"testing"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestZipAddressLookup() {
	pickupKey := models.ServiceItemParamNameZipPickupAddress
	destKey := models.ServiceItemParamNameZipDestAddress

	mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
	paramLookup, err := ServiceParamLookupInitialize(appCfg, suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	suite.T().Run("zip code for the pickup address is present on MTO Shipment", func(t *testing.T) {
		valueStr, err := paramLookup.ServiceParamValue(appCfg, pickupKey)
		suite.FatalNoError(err)
		suite.Equal(mtoServiceItem.MTOShipment.PickupAddress.PostalCode, valueStr)
	})

	suite.T().Run("zip code for the destination address is present on MTO Shipment", func(t *testing.T) {
		valueStr, err := paramLookup.ServiceParamValue(appCfg, destKey)
		suite.FatalNoError(err)
		suite.Equal(mtoServiceItem.MTOShipment.DestinationAddress.PostalCode, valueStr)
	})
}
