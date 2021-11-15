package ghcimport

import (
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_importREContract() {
	suite.T().Run("import success", func(t *testing.T) {
		gre := &GHCRateEngineImporter{
			ContractCode: testContractCode,
			ContractName: testContractName,
		}

		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)
		suite.helperCheckContractName(testContractName)
		suite.NotNil(gre.ContractID)
	})

	suite.T().Run("no contract code", func(t *testing.T) {
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

	suite.T().Run("import success, but no contract name", func(t *testing.T) {
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)
		suite.helperCheckContractName(testContractCode)
		suite.NotNil(gre.ContractID)
	})

	suite.T().Run("try to run again with same contract code", func(t *testing.T) {
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
