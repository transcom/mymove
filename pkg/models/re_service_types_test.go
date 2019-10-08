package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReServiceTypeValidation() {
	suite.T().Run("test valid ReServiceType", func(t *testing.T) {
		validReServiceType := models.ReServiceType{
			Code: "123abc",
			Name: "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReServiceType, expErrors)
	})

	suite.T().Run("test empty ReServiceType", func(t *testing.T) {
		invalidReServiceType := models.ReServiceType{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReServiceType, expErrors)
	})
}
