package ghcimport

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceArea() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	// Prerequisite tables must be loaded.
	err := gre.importREContract(suite.DB())
	suite.NoError(err)

	suite.T().Run("import success", func(t *testing.T) {
		err = gre.importREDomesticServiceArea(suite.DB())
		suite.NoError(err)
		suite.helperVerifyServiceAreaCount(testContractCode)
		suite.NotNil(gre.serviceAreaToIDMap)

		// Spot check a service area
		suite.helperCheckServiceAreaValue(testContractCode)
	})

	suite.T().Run("Run a second time with one changed row, should still succeed", func(t *testing.T) {
		// Get contract UUID.
		var contract models.ReContract
		err = suite.DB().Where("code = ?", testContractCode).First(&contract)
		suite.NoError(err)

		// Change a service area and remove a zip and see if they return as they were before.
		var serviceArea models.ReDomesticServiceArea
		err = suite.DB().
			Where("contract_id = ?", contract.ID).
			Where("service_area = '452'").
			First(&serviceArea)
		suite.NoError(err)
		serviceArea.BasePointCity = "New City"
		suite.MustSave(&serviceArea)

		var zip3 models.ReZip3
		err = suite.DB().
			Where("contract_id = ?", contract.ID).
			Where("zip3 = '647'").
			First(&zip3)
		suite.NoError(err)
		suite.MustDestroy(&zip3)

		err = gre.importREDomesticServiceArea(suite.DB())
		suite.NoError(err)
		suite.helperVerifyServiceAreaCount(testContractCode)
		suite.NotNil(gre.serviceAreaToIDMap)

		// Check to see if data changed above has been reverted.
		suite.helperCheckServiceAreaValue(testContractCode)
	})

	gre2 := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode2,
	}

	// Prerequisite tables must be loaded.
	err = gre2.importREContract(suite.DB())
	suite.NoError(err)

	suite.T().Run("Run with a different contract code, should add new records", func(t *testing.T) {
		err = gre2.importREDomesticServiceArea(suite.DB())
		suite.NoError(err)
		suite.helperVerifyServiceAreaCount(testContractCode2)
		suite.NotNil(gre2.serviceAreaToIDMap)

		// Spot check a service area
		suite.helperCheckServiceAreaValue(testContractCode2)

		// Make sure the other contract's records are still there too.
		suite.helperVerifyServiceAreaCount(testContractCode)
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyServiceAreaCount(contractCode string) {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", contractCode).First(&contract)
	suite.NoError(err)

	// Domestic service areas count
	count, err := suite.DB().Where("contract_id = ?", contract.ID).Count(&models.ReDomesticServiceArea{})
	suite.NoError(err)
	suite.Equal(4, count)

	// Zip3s count
	count, err = suite.DB().Where("contract_id = ?", contract.ID).Count(&models.ReZip3{})
	suite.NoError(err)
	suite.Equal(18, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckServiceAreaValue(contractCode string) {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", contractCode).First(&contract)
	suite.NoError(err)

	var serviceArea models.ReDomesticServiceArea
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_area = '452'").
		First(&serviceArea)
	suite.NoError(err)

	suite.Equal("Butler", serviceArea.BasePointCity)
	suite.Equal("MO", serviceArea.State)
	suite.Equal(1, serviceArea.ServicesSchedule)
	suite.Equal(3, serviceArea.SITPDSchedule)

	expectedZip3s := []string{"647", "648", "656", "657", "658"}

	var zip3s models.ReZip3s
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("domestic_service_area_id = ?", serviceArea.ID).
		All(&zip3s)
	suite.NoError(err)
	for _, zip3 := range zip3s {
		found := false
		for _, expectedZip3 := range expectedZip3s {
			if zip3.Zip3 == expectedZip3 {
				found = true
			}
		}
		suite.True(found, "Could not find zip3 [%s]", zip3.Zip3)
	}

	suite.Len(zip3s, len(expectedZip3s), "Zip3 lengths did not match")
}
