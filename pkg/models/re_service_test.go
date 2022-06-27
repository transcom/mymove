package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReServiceValidation() {
	suite.Run("test valid ReService", func() {
		validReService := models.ReService{
			Code: "123abc",
			Name: "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReService, expErrors)
	})

	suite.Run("test empty ReService", func() {
		emptyReService := models.ReService{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReService, expErrors)
	})
}
