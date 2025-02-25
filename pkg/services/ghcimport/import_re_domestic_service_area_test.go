package ghcimport

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceArea() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}

	setupTestData := func() {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importREDomesticServiceArea(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("import success", func() {
		setupTestData()
		suite.helperVerifyServiceAreaCount(testContractCode)
		suite.NotNil(gre.serviceAreaToIDMap)

		// Spot check a service area
		suite.helperCheckServiceAreaValue(testContractCode)
	})

	suite.Run("Run a second time with changed data, should still succeed", func() {
		setupTestData()

		// Get contract UUID.
		var contract models.ReContract
		err := suite.DB().Where("code = ?", testContractCode).First(&contract)
		suite.NoError(err)

		// Change a service area and remove/change a zip and see if they return as they were before.
		var serviceArea models.ReDomesticServiceArea
		err = suite.DB().
			Where("contract_id = ?", contract.ID).
			Where("service_area = '452'").
			First(&serviceArea)
		suite.NoError(err)
		serviceArea.ServicesSchedule = 2
		serviceArea.SITPDSchedule = 2
		suite.MustSave(&serviceArea)

		var zip3ToDelete models.ReZip3
		err = suite.DB().
			Where("contract_id = ?", contract.ID).
			Where("zip3 = '647'").
			First(&zip3ToDelete)
		suite.NoError(err)
		suite.MustDestroy(&zip3ToDelete)

		var zip3ToUpdate models.ReZip3
		err = suite.DB().
			Where("contract_id = ?", contract.ID).
			Where("zip3 = '657'").
			First(&zip3ToUpdate)
		suite.NoError(err)
		zip3ToUpdate.BasePointCity = "New City"
		zip3ToUpdate.State = "XX"
		suite.MustSave(&zip3ToUpdate)

		err = gre.importREDomesticServiceArea(suite.AppContextForTest())
		suite.NoError(err)
		suite.helperVerifyServiceAreaCount(testContractCode)
		suite.NotNil(gre.serviceAreaToIDMap)

		// Check to see if data changed above has been reverted.
		suite.helperCheckServiceAreaValue(testContractCode)
	})

	suite.Run("Run with a different contract code, should add new records", func() {
		setupTestData()

		gre2 := &GHCRateEngineImporter{
			ContractCode: testContractCode2,
		}

		// Prerequisite tables must be loaded.
		err := gre2.importREContract(suite.AppContextForTest())
		suite.NoError(err)
		err = gre2.importREDomesticServiceArea(suite.AppContextForTest())
		suite.NoError(err)

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
	suite.Equal(5, count)

	// Zip3s count
	count, err = suite.DB().Where("contract_id = ?", contract.ID).Count(&models.ReZip3{})
	suite.NoError(err)
	suite.Equal(19, count)
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

	suite.Equal(1, serviceArea.ServicesSchedule)
	suite.Equal(3, serviceArea.SITPDSchedule)

	expectedZip3s := []struct {
		zip3  string
		city  string
		state string
	}{
		{"647", "Butler", "MO"},
		{"648", "Neosho", "MO"},
		{"656", "Springfield", "MO"},
		{"657", "Springfield", "MO"},
		{"658", "Springfield", "MO"},
	}

	var zip3s models.ReZip3s
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("domestic_service_area_id = ?", serviceArea.ID).
		All(&zip3s)
	suite.NoError(err)
	for _, zip3 := range zip3s {
		found := false
		for _, expectedZip3 := range expectedZip3s {
			if zip3.Zip3 == expectedZip3.zip3 {
				found = true
				suite.Equal(expectedZip3.city, zip3.BasePointCity)
				suite.Equal(expectedZip3.state, zip3.State)
			}
		}
		suite.True(found, "Could not find zip3 [%s]", zip3.Zip3)
	}

	suite.Len(zip3s, len(expectedZip3s), "Zip3 lengths did not match")
}
