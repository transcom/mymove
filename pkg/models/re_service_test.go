package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReServiceValidation() {
	suite.T().Run("test valid ReService", func(t *testing.T) {
		validReService := models.ReService{
			Code: "123abc",
			Name: "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReService, expErrors)
	})

	suite.T().Run("test empty ReService", func(t *testing.T) {
		invalidReService := models.ReService{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReService, expErrors)
	})
}
