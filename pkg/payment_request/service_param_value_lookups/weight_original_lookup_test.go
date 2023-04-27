package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestWeightOriginalLookupForShipment() {
	key := models.ServiceItemParamNameWeightOriginal

	suite.Run("actual weight is present on MTO Shipment", func() {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1234", valueStr)
	})

	suite.Run("nil PrimeActualWeight", func() {
		// Set the actual weight to nil
		mtoServiceItem, paymentRequest, _ := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.PrimeActualWeight = nil
		suite.MustSave(&mtoShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find actual weight for MTOShipmentID [%s]", mtoShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})
}

func (suite *ServiceParamValueLookupsSuite) TestWeightOriginalLookupForShuttling() {
	key := models.ServiceItemParamNameWeightOriginal

	suite.Run("actual weight is present on MTO Shipment", func() {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDOSHUT, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1234", valueStr)
	})

	suite.Run("nil ActualWeight", func() {
		// Set the actual weight to nil
		mtoServiceItem, paymentRequest, _ := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDDSHUT, models.MTOShipmentTypeHHG)
		mtoServiceItem.ActualWeight = nil
		suite.MustSave(&mtoServiceItem)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find actual weight for MTOServiceItemID [%s]", mtoServiceItem.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})
}
