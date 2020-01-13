package ghcimport

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importRETaskOrderFees() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	suite.T().Run("import success", func(t *testing.T) {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.DB())
		suite.NoError(err)

		err = gre.importREContractYears(suite.DB())
		suite.NoError(err)

		//err = gre.importRERateArea(suite.DB())
		//suite.NoError(err)

		err = gre.loadServiceMap(suite.DB())
		suite.NoError(err)

		err = gre.importRETaskOrderFees(suite.DB())
		suite.NoError(err)
		suite.helperVerifyTaskOrderFees()
		suite.helperCheckTaskOrderFees()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importRETaskOrderFees(suite.DB())
		if suite.Error(err) {
			suite.Contains(err.Error(), "duplicate key value violates unique constraint")
		}

		// Check to see if anything else changed
		suite.helperVerifyTaskOrderFees()
		suite.helperCheckTaskOrderFees()
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyTaskOrderFees() {
	count, err := suite.DB().Count(&models.ReTaskOrderFees{})
	suite.NoError(err)
	suite.Equal(16, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckTaskOrderFees() {
	// Get service UUID.
	var serviceMS models.ReService
	err := suite.DB().Where("code = 'MS'").First(&serviceMS)
	suite.NoError(err)

	// Get contract year UUID.
	var contractYear models.ReContractYear
	err = suite.DB().Where("name = 'Base Period Year 1'").First(&contractYear)
	suite.NoError(err)

	var taskOrderFee models.ReTaskOrderFee
	err = suite.DB().
		Where("service_id = $1", serviceMS.ID).
		Where("contract_year_id = $2", contractYear.ID).
		First(&taskOrderFee)
	suite.NoError(err)
	suite.Equal(unit.Cents(45115), taskOrderFee.PriceCents)
}
