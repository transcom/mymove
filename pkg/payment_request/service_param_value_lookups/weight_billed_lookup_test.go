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

	suite.Run("ISLH is minimum of 500", func() {
		// Set the original weight to exactly 110% of estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(100), unit.Pound(110), models.ReServiceCodeISLH, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("500", valueStr)
	})

	suite.Run("original is 120% of estimated weight but there is no adjusted weight", func() {
		// Set the original weight to about 120% of estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1481), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1357", valueStr)
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
		{models.ReServiceCodeDPK, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDNPK, unit.Pound(450), "500", models.MTOShipmentTypeHHGIntoNTSDom},
		{models.ReServiceCodeDUPK, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		// Domestic SIT
		{models.ReServiceCodeDOFSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDDFSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDOASIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDDASIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDOPSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDDDSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		// International
		{models.ReServiceCodeISLH, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeUBP, unit.Pound(250), "300", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeISLH, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeUBP, unit.Pound(250), "300", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeISLH, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeUBP, unit.Pound(250), "300", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIHPK, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIHUPK, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIUBPK, unit.Pound(250), "300", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIUBUPK, unit.Pound(250), "300", models.MTOShipmentTypeHHG},
		// International SIT
		{models.ReServiceCodeIOFSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIDFSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIOASIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIDASIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIOPSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIDDSIT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		// Misc
		{models.ReServiceCodeFSC, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDLH, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeDSH, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
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

	// Setup data for testing PPM minimums
	serviceCodesForPPM := []struct {
		code            models.ReServiceCode
		originalWeight  unit.Pound
		expectedMinimum string
		shipmentType    models.MTOShipmentType
	}{
		// PPM Service items
		{models.ReServiceCodeDLH, unit.Pound(440), "440", models.MTOShipmentTypePPM},
		{models.ReServiceCodeDOP, unit.Pound(440), "440", models.MTOShipmentTypePPM},
		{models.ReServiceCodeDDP, unit.Pound(440), "440", models.MTOShipmentTypePPM},
		{models.ReServiceCodeDPK, unit.Pound(440), "440", models.MTOShipmentTypePPM},
		{models.ReServiceCodeDUPK, unit.Pound(440), "440", models.MTOShipmentTypePPM},
		{models.ReServiceCodeFSC, unit.Pound(440), "440", models.MTOShipmentTypePPM},
	}

	for _, data := range serviceCodesForPPM {
		suite.Run("returns the original weight for PPM service items", func() {
			estimatedWeight := unit.Pound(400)
			_, _, paramLookup := suite.setupTestMTOServiceItemWithEstimatedWeightForPPM(&estimatedWeight, &data.originalWeight, data.code) // #nosec G601 new in 1.22.2

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
		{models.ReServiceCodeIOSHUT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIOSHUT, unit.Pound(250), "300", models.MTOShipmentTypeUnaccompaniedBaggage},
		{models.ReServiceCodeIDSHUT, unit.Pound(450), "500", models.MTOShipmentTypeHHG},
		{models.ReServiceCodeIDSHUT, unit.Pound(250), "300", models.MTOShipmentTypeUnaccompaniedBaggage},
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
		mtoServiceItem.MTOShipment.PrimeActualWeight = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)
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
		mtoServiceItem, paymentRequest, _ := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1000), models.ReServiceCodeDOSHUT, models.MTOShipmentTypeHHG)
		mtoServiceItem.EstimatedWeight = nil
		suite.MustSave(&mtoServiceItem)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)

		suite.Equal("1000", valueStr)
	})
}

func (suite *ServiceParamValueLookupsSuite) TestWeightBilledLookupDivertedShipments() {
	key := models.ServiceItemParamNameWeightBilled

	suite.Run("single diverted shipment", func() {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeightOnDiversion(unit.Pound(1234), unit.Pound(1000), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG, true, nil)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("1000", valueStr)
	})

	suite.Run("multiple diverted shipments with different weights", func() {
		// Setup multiple diverted shipments with different weights
		parentEstimatedWeight := unit.Pound(800)
		parentActualWeight := unit.Pound(1200)
		childEstimatedWeight := unit.Pound(1600)
		childActualWeight := unit.Pound(2400)

		_, _, _, _, _, childParamLookup := suite.setupTestDivertedShipmentChain(&parentEstimatedWeight, &childEstimatedWeight, &parentActualWeight, &childActualWeight, nil, nil, models.ReServiceCodeDLH, models.MTOShipmentTypeHHGIntoNTSDom)

		// Use the child shipment
		valueStr, err := childParamLookup.ServiceParamValue(suite.AppContextForTest(), key)

		// After looking up the child shipment billable weight we should get the parent's weight because it's lower and in the diversion chain
		suite.FatalNoError(err)
		suite.Equal("880", valueStr)
	})

	suite.Run("multiple diverted shipments with different weights with a prime with lowest", func() {
		// Setup multiple diverted shipments with different weights
		parentEstimatedWeight := unit.Pound(800)
		parentActualWeight := unit.Pound(1200)
		parentReweighWeight := unit.Pound(700)
		childEstimatedWeight := unit.Pound(1600)
		childActualWeight := unit.Pound(2400)

		_, _, _, _, _, childParamLookup := suite.setupTestDivertedShipmentChain(&parentEstimatedWeight, &childEstimatedWeight, &parentActualWeight, &childActualWeight, &parentReweighWeight, nil, models.ReServiceCodeDLH, models.MTOShipmentTypeHHGIntoNTSDom)

		// Use the child shipment
		valueStr, err := childParamLookup.ServiceParamValue(suite.AppContextForTest(), key)

		// After looking up the child shipment billable weight we should get the parent's reweigh weight because it's lower and in the diversion chain
		suite.FatalNoError(err)
		suite.Equal("700", valueStr)
	})

}
