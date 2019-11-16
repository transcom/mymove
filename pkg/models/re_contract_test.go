package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReContractValidations() {
	suite.T().Run("test valid ReContract", func(t *testing.T) {
		validReContract := models.ReContract{
			Code: "ABC",
			Name: "ABC, Inc.",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReContract, expErrors)
	})

	suite.T().Run("test empty ReContract", func(t *testing.T) {
		emptyReContract := models.ReContract{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReContract, expErrors)
	})
}
