package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestContractorValidation() {
	suite.T().Run("test valid Contractor", func(t *testing.T) {
		newContractor := testdatagen.MakeContractor(suite.DB(), testdatagen.Assertions{})

		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&newContractor, expErrors)
	})

	suite.T().Run("test empty Contractor", func(t *testing.T) {
		emptyContractor := models.Contractor{}
		expErrors := map[string][]string{
			"name":            {"Name can not be blank."},
			"type":            {"Type can not be blank."},
			"contract_number": {"ContractNumber can not be blank."},
		}
		suite.verifyValidationErrors(&emptyContractor, expErrors)
	})
}