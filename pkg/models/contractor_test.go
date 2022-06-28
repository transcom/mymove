package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestContractorValidation() {
	suite.Run("test valid Contractor", func() {
		newContractor := models.Contractor{
			Name:           "Contractor 1",
			Type:           "Prime",
			ContractNumber: "HTC111-11-1-1111",
		}

		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&newContractor, expErrors)
	})

	suite.Run("test empty Contractor", func() {
		emptyContractor := models.Contractor{}
		expErrors := map[string][]string{
			"name":            {"Name can not be blank."},
			"type":            {"Type can not be blank."},
			"contract_number": {"ContractNumber can not be blank."},
		}
		suite.verifyValidationErrors(&emptyContractor, expErrors)
	})
}
