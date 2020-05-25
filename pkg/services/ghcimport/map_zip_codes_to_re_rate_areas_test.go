package ghcimport

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_mapZipCodesToReRateAreas() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	suite.T().Run("import prerequisite tables", func(t *testing.T) {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.DB())
		suite.NoError(err)

		err = gre.importREDomesticServiceArea(suite.DB())
		suite.NoError(err)

		err = gre.importRERateArea(suite.DB())
		suite.NoError(err)
	})

	suite.T().Run("map ReZip3 record to correct ReRateArea record", func(t *testing.T) {
		var reContract models.ReContract
		err := suite.DB().Where("code = ?", testContractCode).First(&reContract)
		suite.NoError(err)

		var reZip3 models.ReZip3
		err = suite.DB().Where("contract_id = ?", reContract.ID).First(&reZip3)
		suite.NoError(err)

		rateArea, found := zip3ToRateAreaMappings[reZip3.Zip3]
		suite.Assertions.True(found, "failed to find rate area map for zip3 %s in zip3ToRateAreaMappings", reZip3.Zip3)


		var reRateArea models.ReRateArea
		err = suite.DB().
			Where("contract_id = ?", reContract.ID).
			Where("code = ?",rateArea).
			First(&reRateArea)

		err = gre.mapREZip3sToRERateAreas(suite.DB())
		suite.NoError(err)

		var updatedReZip3 models.ReZip3
		err = suite.DB().Where("contract_id = ?", reContract.ID).First(&updatedReZip3)
		suite.NoError(err)

		suite.Assertions.NotNil(updatedReZip3.RateAreaID, "ReZip3 record %s should not have nil rate_area_id", updatedReZip3.ID)
		suite.Assertions.Equal(updatedReZip3.RateAreaID, &reRateArea.ID, "ReZip3 %s record is mapped to ReRateArea record %s, but should be mapped to %s", updatedReZip3.RateAreaID, reRateArea.ID)
	})

	// TODO: Running this test first causes NotNil checks to fail in tests that come after it
	suite.T().Run("map zip codes to ReRateAreas successfully", func(t *testing.T) {
		err := gre.mapZipCodesToRERateAreas(suite.DB())
		suite.NoError(err)
	})
}
