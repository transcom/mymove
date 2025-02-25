package ghcimport

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_mapZipCodesToReRateAreas() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}

	setupTestData := func() {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importREDomesticServiceArea(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importRERateArea(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("map ReZip3 records to correct ReRateArea records", func() {
		setupTestData()
		reContract, err := suite.helperFetchReContract()
		suite.NoError(err)

		var reZip3 models.ReZip3
		err = suite.DB().
			Where("contract_id = ?", reContract.ID).
			Where("zip3 = ?", "352").
			First(&reZip3)
		suite.NoError(err)

		var reZip3WithMultipleReRateAreas models.ReZip3
		err = suite.DB().
			Where("contract_id = ?", reContract.ID).
			Where("zip3 = ?", "327").
			First(&reZip3WithMultipleReRateAreas)
		suite.NoError(err)

		suite.Nil(reZip3.RateAreaID, "expected ReZip3 record %s to have nil rate_area_id", reZip3.ID)
		suite.Nil(reZip3WithMultipleReRateAreas.RateAreaID, "expected ReZip3 record %s to have nil rate_area_id", reZip3WithMultipleReRateAreas.ID)

		rateAreaCode, found := zip3ToRateAreaMappings[reZip3.Zip3]
		suite.True(found, "failed to find rate area map for zip3 %s in zip3ToRateAreaMappings", reZip3.Zip3)

		zipRateAreaCode, found := zip3ToRateAreaMappings[reZip3WithMultipleReRateAreas.Zip3]
		suite.True(found, "failed to find rate area map for zip3 %s in zip3ToRateAreaMappings", reZip3WithMultipleReRateAreas.Zip3)
		suite.Equal(zipRateAreaCode, "ZIP", "expected rate area code to be ZIP but got %s", zipRateAreaCode)

		reRateArea, err := suite.helperFetchReRateArea(reContract, rateAreaCode)
		suite.NoError(err)

		err = gre.mapREZip3sToRERateAreas(suite.AppContextForTest())
		suite.NoError(err)

		var updatedReZip3 models.ReZip3
		err = suite.DB().
			Where("id = ?", reZip3.ID).
			First(&updatedReZip3)
		suite.NoError(err)

		suite.NotNil(updatedReZip3.RateAreaID, "expected ReZip3 record %s to not have nil rate_area_id", updatedReZip3.ID)
		suite.Equal(*updatedReZip3.RateAreaID, reRateArea.ID, "expected ReZip3 %s record to be mapped to ReRateArea record %s, but got %s", reRateArea.ID, updatedReZip3.RateAreaID)

		var updatedReZip3WithMultipleReRateAreas models.ReZip3
		err = suite.DB().
			Where("id = ?", reZip3WithMultipleReRateAreas.ID).
			First(&updatedReZip3WithMultipleReRateAreas)
		suite.NoError(err)

		suite.Nil(updatedReZip3WithMultipleReRateAreas.RateAreaID, "expected ReZip3 record %s to have nil rate_area_id", updatedReZip3WithMultipleReRateAreas.ID)
		suite.True(updatedReZip3WithMultipleReRateAreas.HasMultipleRateAreas)
	})

	suite.Run("create ReZip5RateArea records and map to correct ReRateArea records", func() {
		setupTestData()
		reContract, err := suite.helperFetchReContract()
		suite.NoError(err)

		reZip5RateAreasCount, err := suite.DB().
			Where("contract_id = ?", reContract.ID).
			Count(&models.ReZip5RateArea{})
		suite.NoError(err)

		suite.Equal(0, reZip5RateAreasCount)

		err = gre.createAndMapREZip5sToRERateAreas(suite.AppContextForTest())
		suite.NoError(err)

		var reZip5RateArea models.ReZip5RateArea
		err = suite.DB().
			Where("contract_id = ?", reContract.ID).
			First(&reZip5RateArea)
		suite.NoError(err)

		rateAreaCode, found := zip5ToRateAreaMappings[reZip5RateArea.Zip5]
		suite.True(found, "failed to find rate area map for zip3 %s in zip3ToRateAreaMappings", reZip5RateArea.Zip5)

		reRateArea, err := suite.helperFetchReRateArea(reContract, rateAreaCode)
		suite.NoError(err)

		suite.Equal(reZip5RateArea.RateAreaID, reRateArea.ID, "expected ReZip3 %s record to be mapped to ReRateArea record %s, but got %s", reRateArea.ID, reZip5RateArea.RateAreaID)
		suite.helperVerifyNumberOfReZip5RateAreasCreated()
	})
}

func (suite *GHCRateEngineImportSuite) helperFetchReContract() (models.ReContract, error) {
	var reContract models.ReContract
	err := suite.DB().
		Where("code = ?", testContractCode).
		First(&reContract)

	return reContract, err
}

func (suite *GHCRateEngineImportSuite) helperFetchReRateArea(reContract models.ReContract, rateAreaCode string) (models.ReRateArea, error) {
	var reRateArea models.ReRateArea

	err := suite.DB().
		Where("contract_id = ?", reContract.ID).
		Where("code = ?", rateAreaCode).
		First(&reRateArea)

	return reRateArea, err
}

func (suite *GHCRateEngineImportSuite) helperVerifyNumberOfReZip5RateAreasCreated() {
	reContract, err := suite.helperFetchReContract()
	suite.NoError(err)

	reZip5RateAreasCount, err := suite.DB().
		Where("contract_id = ?", reContract.ID).
		Count(&models.ReZip5RateArea{})
	suite.NoError(err)

	suite.Equal(922, reZip5RateAreasCount)
}
