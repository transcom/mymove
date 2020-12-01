package ghcimport

import (
	"testing"

	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_importREContractYears() {
	gre := &GHCRateEngineImporter{
		Logger:            suite.logger,
		ContractCode:      testContractCode,
		ContractStartDate: testContractStartDate,
	}

	suite.T().Run("import success", func(t *testing.T) {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.DB())
		suite.NoError(err)

		err = gre.importREContractYears(suite.DB())
		suite.NoError(err)
		suite.helperVerifyContractYears()
		suite.helperCheckContractYearValue()
	})

	suite.T().Run("run a second time; should fail immediately due to date range constraint", func(t *testing.T) {
		err := gre.importREContractYears(suite.DB())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.ExclusionViolation, "re_contract_years_daterange_excl"))
		}

		// Check to see if anything else changed
		suite.helperVerifyContractYears()
		suite.helperCheckContractYearValue()
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyContractYears() {
	count, err := suite.DB().Count(&models.ReContractYear{})
	suite.NoError(err)
	suite.Equal(8, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckContractYearValue() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = $1", testContractCode).First(&contract)
	suite.NoError(err)

	var basePeriod1 models.ReContractYear
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("name = $2", "Base Period Year 1").
		First(&basePeriod1)
	suite.NoError(err)
	suite.Equal(1.0000, basePeriod1.Escalation)
	suite.Equal(1.0000, basePeriod1.EscalationCompounded)

	var optionPeriod1 models.ReContractYear
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("name = $2", "Option Period 1").
		First(&optionPeriod1)
	suite.NoError(err)
	suite.Equal(1.02140, optionPeriod1.Escalation)
	suite.Equal(1.06298, optionPeriod1.EscalationCompounded)

	var awardTerm2 models.ReContractYear
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("name = $2", "Award Term 2").
		First(&awardTerm2)
	suite.NoError(err)
	suite.Equal(1.01940, awardTerm2.Escalation)
	suite.Equal(1.12848, awardTerm2.EscalationCompounded)
}
