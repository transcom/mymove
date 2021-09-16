package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestWeightAdjustedLookup() {
	key := models.ServiceItemParamNameWeightAdjusted

	suite.Run("adjusted weight is present on MTO Shipment", func() {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithAdjustedWeight(unit.Pound(1000), unit.Pound(2000), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
		suite.FatalNoError(err)
		suite.Equal("1000", valueStr)
	})

	suite.Run("nil AdjustedWeight should not cause an error", func() {
		// Set the adjusted weight to nil
		mtoServiceItem, paymentRequest, _ := suite.setupTestMTOServiceItemWithAdjustedWeight(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.BillableWeightCap = nil
		suite.MustSave(&mtoShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.TestAppContext(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
		suite.NoError(err)
		suite.Equal("", valueStr)
	})
}
