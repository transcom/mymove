package ghcimport

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) helperImportRERateArea(action string) {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", testContractCode).First(&contract)
	suite.NoError(err)

	// Update domestic US6B name "Texas-South" to something else and verify it was changed back when done
	var texas *models.ReRateArea
	texas, err = models.FetchReRateAreaItem(suite.DB(), contract.ID, "US68")
	suite.NoError(err)
	suite.Equal(true, suite.NotNil(texas))
	fmt.Printf("\nFetch US6B rate area %v\n\n", texas)
	suite.Equal("Texas-South", texas.Name)

	// Update oconus US8101000 name "Alaska (Zone) I" to something else and verify it was changed back when done
	var alaska *models.ReRateArea
	alaska, err = models.FetchReRateAreaItem(suite.DB(), contract.ID, "US8101000")
	suite.NoError(err)
	suite.Equal(true, suite.NotNil(alaska))
	suite.Equal("Alaska (Zone) I", alaska.Name)

	// Update oconus AS11 name "New South Wales/Australian Capital Territory"
	var wales *models.ReRateArea
	wales, err = models.FetchReRateAreaItem(suite.DB(), contract.ID, "AS11")
	suite.NoError(err)
	suite.Equal(true, suite.NotNil(wales))
	suite.Equal("New South Wales/Australian Capital Territory", wales.Name)

	if action == "setup" {
		modifiedName := "New name"
		texas.Name = modifiedName
		suite.MustSave(texas)
		texas, err = models.FetchReRateAreaItem(suite.DB(), contract.ID, "US68")
		suite.NoError(err)
		suite.Equal(modifiedName, texas.Name)

		modifiedName = "New name 2"
		alaska.Name = modifiedName
		suite.MustSave(alaska)
		alaska, err = models.FetchReRateAreaItem(suite.DB(), contract.ID, "US8101000")
		suite.NoError(err)
		suite.Equal(modifiedName, alaska.Name)

		modifiedName = "New name 3"
		wales.Name = modifiedName
		suite.MustSave(wales)
		wales, err = models.FetchReRateAreaItem(suite.DB(), contract.ID, "AS11")
		suite.NoError(err)
		suite.Equal(modifiedName, wales.Name)
	} else if action == "verify" {
		// nothing to do, verify happens at the top
	}
}

func (suite *GHCRateEngineImportSuite) helperVerifyDomesticRateAreaToIDMap(contractCode string, domesticRateAreaToIDMap map[string]uuid.UUID) {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", contractCode).First(&contract)
	suite.NoError(err)

	suite.NotEqual(map[string]uuid.UUID(nil), domesticRateAreaToIDMap)
	count, dbErr := suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("is_oconus = 'false'").
		Count(models.ReRateArea{})
	suite.NoError(dbErr)

	suite.Equal(4, count)
	suite.Equal(count, len(domesticRateAreaToIDMap))

	var rateArea models.ReRateArea
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("code = 'US68'").
		First(&rateArea)
	suite.NoError(err)

	suite.Equal("Texas-South", rateArea.Name)
	suite.Equal(rateArea.ID, domesticRateAreaToIDMap["US68"])

	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("code = 'US47'").
		First(&rateArea)
	suite.NoError(err)

	suite.Equal("Alabama", rateArea.Name)
	suite.Equal(rateArea.ID, domesticRateAreaToIDMap["US47"])
}

func (suite *GHCRateEngineImportSuite) helperVerifyInternationalRateAreaToIDMap(contractCode string, internationalRateAreaToIDMap map[string]uuid.UUID) {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", contractCode).First(&contract)
	suite.NoError(err)

	suite.NotEqual(map[string]uuid.UUID(nil), internationalRateAreaToIDMap)
	count, dbErr := suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("is_oconus = 'true'").
		Count(models.ReRateArea{})
	suite.NoError(dbErr)

	suite.Equal(5, count)
	suite.Equal(count, len(internationalRateAreaToIDMap))

	var rateArea models.ReRateArea
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("code = 'GE'").
		First(&rateArea)
	suite.NoError(err)

	suite.Equal("Germany", rateArea.Name)
	suite.Equal(rateArea.ID, internationalRateAreaToIDMap["GE"])

	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("code = 'US8101000'").
		First(&rateArea)
	suite.NoError(err)

	suite.Equal("Alaska (Zone) I", rateArea.Name)
	suite.Equal(rateArea.ID, internationalRateAreaToIDMap["US8101000"])
}

func (suite *GHCRateEngineImportSuite) helperImportRERateAreaDropStage(action string) {
	if action == "setup" {
		// drop a staging table that we are depending on to do import
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s;", "stage_conus_to_oconus_prices")
		dropErr := suite.DB().RawQuery(dropQuery).Exec()
		suite.NoError(dropErr)
	}
}

func (suite *GHCRateEngineImportSuite) helperImportRERateAreaVerifyImportComplete(contractCode string) {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", contractCode).First(&contract)
	suite.NoError(err)

	var rateArea models.ReRateArea
	count, countErr := suite.DB().Where("contract_id = ?", contract.ID).Count(&rateArea)

	suite.NoError(countErr)
	suite.Equal(9, count)
}

func (suite *GHCRateEngineImportSuite) TestGHCRateEngineImporter_importRERateArea() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	// Prerequisite tables must be loaded.
	err := gre.importREContract(suite.DB())
	suite.NoError(err)

	suite.T().Run("Successfully run import with staged staging data (empty RE tables)", func(t *testing.T) {
		err = gre.importRERateArea(suite.DB())
		suite.NoError(err)
		suite.helperImportRERateAreaVerifyImportComplete(testContractCode)

		suite.helperVerifyDomesticRateAreaToIDMap(testContractCode, gre.domesticRateAreaToIDMap)
		suite.helperVerifyInternationalRateAreaToIDMap(testContractCode, gre.internationalRateAreaToIDMap)
	})

	suite.T().Run("Successfully run import, 2nd time, with staged staging data and filled in RE tables", func(t *testing.T) {
		err = gre.importRERateArea(suite.DB())
		suite.NoError(err)
		suite.helperImportRERateAreaVerifyImportComplete(testContractCode)

		suite.helperVerifyDomesticRateAreaToIDMap(testContractCode, gre.domesticRateAreaToIDMap)
		suite.helperVerifyInternationalRateAreaToIDMap(testContractCode, gre.internationalRateAreaToIDMap)
	})

	suite.T().Run("Successfully run import, prefilled re_rate_areas, update existing rate area from import", func(t *testing.T) {
		suite.helperImportRERateArea("setup")

		err = gre.importRERateArea(suite.DB())
		suite.NoError(err)
		suite.helperImportRERateAreaVerifyImportComplete(testContractCode)

		suite.helperVerifyDomesticRateAreaToIDMap(testContractCode, gre.domesticRateAreaToIDMap)
		suite.helperVerifyInternationalRateAreaToIDMap(testContractCode, gre.internationalRateAreaToIDMap)
		suite.helperImportRERateArea("verify")
	})

	suite.T().Run("Fail to run import, missing staging table", func(t *testing.T) {
		suite.helperImportRERateAreaDropStage("setup")

		err = gre.importRERateArea(suite.DB())
		suite.Error(err)

		suite.helperSetupStagingTables()
	})

	gre2 := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode2,
	}

	// Prerequisite tables must be loaded.
	err = gre2.importREContract(suite.DB())
	suite.NoError(err)

	suite.T().Run("Run with a different contract code, should add new records", func(t *testing.T) {
		err = gre2.importRERateArea(suite.DB())
		suite.NoError(err)
		suite.helperImportRERateAreaVerifyImportComplete(testContractCode2)

		suite.helperVerifyDomesticRateAreaToIDMap(testContractCode2, gre2.domesticRateAreaToIDMap)
		suite.helperVerifyInternationalRateAreaToIDMap(testContractCode2, gre2.internationalRateAreaToIDMap)

		// Make sure the other contract's records are still there too.
		suite.helperImportRERateAreaVerifyImportComplete(testContractCode)
	})
}
