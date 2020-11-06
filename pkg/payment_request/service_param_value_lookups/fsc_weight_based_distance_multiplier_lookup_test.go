package serviceparamvaluelookups

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestFSCWeightBasedDistanceMultiplierLookup() {
	key := models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier

	suite.T().Run("correct weight based distance multiplier is returned for billed actual weight less than 5,000 pounds", func(t *testing.T) {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(3000), unit.Pound(3000), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(key)

		suite.FatalNoError(err)
		suite.Equal("0.000417", valueStr)
	})

	suite.T().Run("correct weight based distance multiplier is returned for billed actual weight greater than 5,000 pounds but less than 10,001 pounds", func(t *testing.T) {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(9500), unit.Pound(9500), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(key)

		suite.FatalNoError(err)
		suite.Equal("0.0006255", valueStr)
	})

	suite.T().Run("correct weight based distance multiplier is returned for billed actual weight greater than 10,000 pounds but less than 24,001 pounds", func(t *testing.T) {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(14750), unit.Pound(14750), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(key)

		suite.FatalNoError(err)
		suite.Equal("0.000834", valueStr)
	})

	suite.T().Run("correct weight based distance multiplier is returned for billed actual weight greater than 24,000 pounds", func(t *testing.T) {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(32225), unit.Pound(32225), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(key)

		suite.FatalNoError(err)
		suite.Equal("0.00139", valueStr)
	})

	suite.T().Run("correct weight based distance multiplier is returned for billed actual weight greater than 24,000 pounds", func(t *testing.T) {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(32225), unit.Pound(32225), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(key)

		suite.FatalNoError(err)
		suite.Equal("0.00139", valueStr)
	})
}
