package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReRateAreaValidation() {
	suite.T().Run("test valid ReRateArea", func(t *testing.T) {
		validReRateArea := models.ReRateArea{
			IsOconus: true,
			Code:     "123abc",
			Name:     "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReRateArea, expErrors)
	})

	suite.T().Run("test empty ReRateArea", func(t *testing.T) {
		invalidReRateArea := models.ReRateArea{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReRateArea, expErrors)
	})
}
