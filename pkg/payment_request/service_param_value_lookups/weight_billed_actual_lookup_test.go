package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestWeightBilledActualLookup() {
	key := models.ServiceItemParamNameWeightBilledActual

	suite.Run("estimated and actual are the same", func() {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		suite.Equal("1234", valueStr)
	})

	suite.Run("estimated is greater than actual", func() {
		// Set the actual weight to less than estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1024), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		suite.Equal("1024", valueStr)
	})

	suite.Run("actual is exactly 110% of estimated weight", func() {
		// Set the actual weight to exactly 110% of estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(100), unit.Pound(110), models.ReServiceCodeNSTH, models.MTOShipmentTypeHHG)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		suite.Equal("110", valueStr)
	})

	suite.Run("actual is 120% of estimated weight", func() {
		// Set the actual weight to about 120% of estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), unit.Pound(1481), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		suite.Equal("1357", valueStr)
	})

	suite.Run("rounds to the nearest whole pound", func() {
		// Set the weights so that a fraction of a pound is returned
		_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1235), unit.Pound(1482), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		suite.Equal("1359", valueStr)
	})

	// Setup data for testing all minimums
	serviceCodesWithMinimum := []struct {
		code            models.ReServiceCode
		actualWeight    unit.Pound
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
		suite.Run(fmt.Sprintf("actual below minimum service code %s", data.code), func() {
			// Set the actual weight to below minimum
			_, _, paramLookup := suite.setupTestMTOServiceItemWithWeight(unit.Pound(1234), data.actualWeight, data.code, data.shipmentType)

			appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
			valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
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

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		paramLookup, err := ServiceParamLookupInitialize(appCfg, suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
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

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		paramLookup, err := ServiceParamLookupInitialize(appCfg, suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.NoError(err)
		suite.Equal("500", valueStr)
	})
}

func (suite *ServiceParamValueLookupsSuite) TestShuttleWeightBilledActualLookup() {
	key := models.ServiceItemParamNameWeightBilledActual

	suite.Run("estimated and actual are the same", func() {
		_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDOSHUT, models.MTOShipmentTypeHHG)
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		suite.Equal("1234", valueStr)
	})

	suite.Run("estimated is greater than actual", func() {
		// Set the actual weight to less than estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1024), models.ReServiceCodeDDSHUT, models.MTOShipmentTypeHHG)
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		suite.Equal("1024", valueStr)
	})

	suite.Run("actual is 120% of estimated weight", func() {
		// Set the actual weight to about 120% of estimated weight
		_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1481), models.ReServiceCodeDOSHUT, models.MTOShipmentTypeHHG)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		suite.Equal("1357", valueStr)
	})

	suite.Run("rounds to the nearest whole pound", func() {
		// Set the weights so that a fraction of a pound is returned
		_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1235), unit.Pound(1482), models.ReServiceCodeDDSHUT, models.MTOShipmentTypeHHG)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.FatalNoError(err)
		suite.Equal("1359", valueStr)
	})

	// Setup data for testing all minimums
	serviceCodesWithMinimum := []struct {
		code            models.ReServiceCode
		actualWeight    unit.Pound
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
			// Set the actual weight to below minimum
			_, _, paramLookup := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), data.actualWeight, data.code, data.shipmentType)

			appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
			valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
			suite.FatalNoError(err)
			suite.Equal(data.expectedMinimum, valueStr)
		})
	}

	suite.Run("nil ActualWeight", func() {
		// Set the actual weight to nil
		mtoServiceItem, paymentRequest, _ := suite.setupTestMTOServiceItemWithShuttleWeight(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDDSHUT, models.MTOShipmentTypeHHG)
		mtoServiceItem.ActualWeight = nil
		suite.MustSave(&mtoServiceItem)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		paramLookup, err := ServiceParamLookupInitialize(appCfg, suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
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

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		paramLookup, err := ServiceParamLookupInitialize(appCfg, suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(appCfg, key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find estimated weight for MTOServiceItemID [%s]", mtoServiceItem.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})
}
