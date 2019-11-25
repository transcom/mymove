package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestContractorInstantiation() {
	contractor := &Contractor{}

	expErrors := map[string][]string{
		"name":            {"Name can not be blank."},
		"contract_number": {"ContractNumber can not be blank."},
		"type":            {"Type can not be blank."},
	}

	suite.verifyValidationErrors(contractor, expErrors)
}
