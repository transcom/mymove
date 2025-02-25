package ghcimport

import (
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importRETaskOrderFees() {
	gre := &GHCRateEngineImporter{
		ContractCode:      testContractCode,
		ContractStartDate: testContractStartDate,
	}

	setupTestData := func() {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importREContractYears(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.loadServiceMap(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importRETaskOrderFees(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("import success", func() {
		setupTestData()
		suite.helperVerifyTaskOrderFees()
		suite.helperCheckTaskOrderFees()
	})

	suite.Run("run a second time; should fail immediately due to constraint violation", func() {
		setupTestData()
		err := gre.importRETaskOrderFees(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_task_order_fees_unique_key"))
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyTaskOrderFees() {
	count, err := suite.DB().Count(&models.ReTaskOrderFee{})
	suite.NoError(err)
	suite.Equal(16, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckTaskOrderFees() {
	// Get service UUID.
	var serviceMS models.ReService
	err := suite.DB().Where("code = $1", models.ReServiceCodeMS).First(&serviceMS)
	suite.NoError(err)

	var serviceCS models.ReService
	err = suite.DB().Where("code = $1", models.ReServiceCodeCS).First(&serviceCS)
	suite.NoError(err)

	// Get contract year UUID.
	var contractYear models.ReContractYear
	err = suite.DB().Where("name = 'Base Period Year 1'").First(&contractYear)
	suite.NoError(err)

	var taskOrderFeeMS models.ReTaskOrderFee
	err = suite.DB().
		Where("service_id = $1", serviceMS.ID).
		Where("contract_year_id = $2", contractYear.ID).
		First(&taskOrderFeeMS)
	suite.NoError(err)
	suite.Equal(unit.Cents(45115), taskOrderFeeMS.PriceCents)

	var taskOrderFeeCS models.ReTaskOrderFee
	err = suite.DB().
		Where("service_id = $1", serviceCS.ID).
		Where("contract_year_id = $2", contractYear.ID).
		First(&taskOrderFeeCS)
	suite.NoError(err)
	suite.Equal(unit.Cents(22263), taskOrderFeeCS.PriceCents)
}
