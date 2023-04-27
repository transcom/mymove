package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestWeightAdjustedLookup() {
	key := models.ServiceItemParamNameWeightAdjusted

	suite.Run("adjusted weight is present on MTO Shipment", func() {
		adjustedWeight := unit.Pound(1000)
		_, _, paramLookup := suite.setupTestMTOServiceItemWithAdjustedWeight(&adjustedWeight, unit.Pound(2000), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1000", valueStr)
	})

	suite.Run("nil AdjustedWeight should not cause an error", func() {
		// Set the adjusted weight to nil
		_, _, paramLookup := suite.setupTestMTOServiceItemWithAdjustedWeight(nil, unit.Pound(450), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
		suite.Equal("", valueStr)
	})
}
