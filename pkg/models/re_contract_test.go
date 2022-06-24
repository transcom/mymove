package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReContractValidations() {
	suite.Run("test valid ReContract", func() {
		validReContract := models.ReContract{
			Code: "ABC",
			Name: "ABC, Inc.",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReContract, expErrors)
	})

	suite.Run("test empty ReContract", func() {
		emptyReContract := models.ReContract{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReContract, expErrors)
	})
}
