package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestWeightBilledLookup() {
	key := models.ServiceItemParamNameWeightBilled

	suite.Run("estimated and original are the same", func() {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1234", valueStr)
	})

	suite.Run("estimated is greater than original", func() {
		// Set the original weight to less than estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1024), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1024", valueStr)
	})

	suite.Run("original is exactly 110% of estimated weight", func() {
		// Set the original weight to exactly 110% of estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(100), unit.Pound(110), models.ReServiceCodeNSTH, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("110", valueStr)
	})

	suite.Run("original is 120% of estimated weight but there is no adjusted weight", func() {
		// Set the original weight to about 120% of estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1481), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1481", valueStr)
	})

	suite.Run("has only original weight", func() {
		// Set the original weight only; no estimated, reweigh, or adjusted weight.
		_, _, paramLookup := suite.setupTestMTOServiceItemWithOriginalWeightOnly(unit.Pound(1755), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1755", valueStr)
	})

	suite.Run("has reweigh weight where reweigh is lower", func() {
		// Set the original weight to greater than the reweigh weight. Lower weight (reweigh) should win.
		_, _, paramLookup := suite.setupTestMTOServiceItemWithReweigh(unit.Pound(1450), unit.Pound(1481), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1450", valueStr)
	})

	suite.Run("has reweigh weight where reweigh is higher", func() {
		// Set the original weight to less than the reweigh weight. Lower weight (original) should win.
		_, _, paramLookup := suite.setupTestMTOServiceItemWithReweigh(unit.Pound(1500), unit.Pound(1480), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1480", valueStr)
	})

	suite.Run("has adjusted weight", func() {
		// Set the original weight to greater than the adjusted weight (adjusted weight should always win)
		adjustedWeight := unit.Pound(1400)
		_, _, paramLookup := suite.setupTestMTOServiceItemWithAdjustedWeight(&adjustedWeight, unit.Pound(1481), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1400", valueStr)
	})

	// Setup data for testing all minimums
	serviceCodesWithMinimum := []struct {
		code            models.ReServiceCode
		originalWeight  unit.Pound
		expectedMinimum string
		shipmentType    models.MTOShipmentType
	}{
		// Domestic
		{models.ReServiceCodeDLH, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDSH, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDOP, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDDP, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDOFSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDDFSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDOASIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDDASIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDOPSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDDDSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDPK, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDNPK, unit.Pound(450), "500", models.MTOShipmentTypeHHGIntoNTSDom},
		{models.ReServiceCodeDUPK, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		// International
		{models.ReServiceCodeIOOLH, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIOOUB, unit.Pound(250), "300", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeICOLH, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeICOUB, unit.Pound(250), "300", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIOCLH, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIOCUB, unit.Pound(250), "300", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIHPK, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIHUPK, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIUBPK, unit.Pound(250), "300", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIUBUPK, unit.Pound(250), "300", models.MTOShipmentTypeInternationalHHG},
		// International SIT
		{models.ReServiceCodeIOFSIT, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIDFSIT, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIOASIT, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIDASIT, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIOPSIT, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIDDSIT, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
	}

	// test minimums are correct
	for _, data := range serviceCodesWithMinimum {
		suite.Run(fmt.Sprintf("original below minimum service code %s", data.code), func() {
			// Set the original weight to below minimum
			_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), data.originalWeight, data.code, data.shipmentType)

			valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
			suite.FatalNoError(err)
			suite.Equal(data.expectedMinimum, valueStr)
		})
	}

	suite.Run("nil PrimeActualWeight", func() {
		// Set the actual weight to nil
		mtoServiceItem, paymentRequest, _ := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
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

	suite.Run("nil PrimeEstimatedWeight", func() {
		// Set the estimated weight to nil
		mtoServiceItem, paymentRequest, _ := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.PrimeEstimatedWeight = nil
		suite.MustSave(&mtoShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
		suite.Equal("500", valueStr)
	})
}

func (suite *ServiceParamValueLookupsSuite) TestShuttleWeightBilledLookup() {
	key := models.ServiceItemParamNameWeightBilled

	suite.Run("estimated and original are the same", func() {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDOSHUT, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1234", valueStr)
	})

	suite.Run("estimated is greater than original", func() {
		// Set the original weight to less than estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1024), models.ReServiceCodeDDSHUT, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1024", valueStr)
	})

	suite.Run("original is 120% of estimated weight", func() {
		// Set the original weight to about 120% of estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1481), models.ReServiceCodeDOSHUT, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1357", valueStr)
	})

	suite.Run("rounds to the nearest whole pound", func() {
		// Set the weights so that a fraction of a pound is returned
		_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1235), unit.Pound(1482), models.ReServiceCodeDDSHUT, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1359", valueStr)
	})

	// Setup data for testing all minimums
	serviceCodesWithMinimum := []struct {
		code            models.ReServiceCode
		originalWeight  unit.Pound
		expectedMinimum string
		shipmentType    models.MTOShipmentType
	}{
		// Domestic Shuttle
		{models.ReServiceCodeDOSHUT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDDSHUT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		// International Shuttle
		{models.ReServiceCodeIOSHUT, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIOSHUT, unit.Pound(250), "300", models.MTOShipmentTypeInternationalUB},
		{models.ReServiceCodeIDSHUT, unit.Pound(450), "500", models.MTOShipmentTypeInternationalHHG},
		{models.ReServiceCodeIDSHUT, unit.Pound(250), "300", models.MTOShipmentTypeInternationalUB},
	}

	// test minimums are correct
	for _, data := range serviceCodesWithMinimum {
		suite.Run(fmt.Sprintf("actual below minimum service code %s", data.code), func() {
			// Set the original weight to below minimum
			_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), data.originalWeight, data.code, data.shipmentType)

			valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
			suite.FatalNoError(err)
			suite.Equal(data.expectedMinimum, valueStr)
		})
	}

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

	suite.Run("nil EstimatedWeight", func() {
		// Set the estimated weight to nil
		mtoServiceItem, paymentRequest, _ := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDOSHUT, models.MTOShipmentTypeHHG)
		mtoServiceItem.EstimatedWeight = nil
		suite.MustSave(&mtoServiceItem)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find estimated weight for MTOServiceItemID [%s]", mtoServiceItem.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})
}
