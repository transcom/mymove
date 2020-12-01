package serviceparamvaluelookups

import (
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestWeightEstimatedLookup() {
	key := models.ServiceItemParamNameWeightEstimated

	suite.T().Run("estimated weight is present on MTO Shipment", func(t *testing.T) {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("1234", valueStr)
	})

	suite.T().Run("nil PrimeEstimatedWeight", func(t *testing.T) {
		// Set the estimated weight to nil
		mtoServiceItem, paymentRequest, _ := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.PrimeEstimatedWeight = nil
		suite.MustSave(&mtoShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find estimated weight for MTOShipmentID [%s]", mtoShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})
}
