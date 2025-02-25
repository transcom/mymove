package ghcimport

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_importREContract() {
	suite.Run("import success", func() {
		gre := &GHCRateEngineImporter{
			ContractCode: testContractCode,
			ContractName: testContractName,
		}

		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)
		suite.helperCheckContractName(testContractName)
		suite.NotNil(gre.ContractID)
	})

	suite.Run("no contract code", func() {
		gre := &GHCRateEngineImporter{}

		err := gre.importREContract(suite.AppContextForTest())
		if suite.Error(err) {
			suite.Equal("no contract code provided", err.Error())
		}
	})
}

func (suite *GHCRateEngineImportSuite) Test_importREContract_runTwice() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}

	setupTestData := func() {
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("import success, but no contract name", func() {
		setupTestData()
		suite.helperCheckContractName(testContractCode)
		suite.NotNil(gre.ContractID)
	})

	suite.Run("run twice with same contract code - should fail the second time", func() {
		setupTestData()
		err := gre.importREContract(suite.AppContextForTest())
		if suite.Error(err) {
			expected := fmt.Sprintf("the provided contract code [%s] already exists", testContractCode)
			suite.Equal(expected, err.Error())
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperCheckContractName(expectedName string) {
	var contract models.ReContract
	err := suite.DB().Where("code = ?", testContractCode).First(&contract)
	suite.NoError(err)
	suite.Equal(expectedName, contract.Name)
}
