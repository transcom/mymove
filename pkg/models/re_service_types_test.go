package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ReServiceTypeValidation() {
	validReServiceType := ReServiceType{
		Code: "123abc",
		Name: "California",
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validReServiceType, expErrors)

	invalidReServiceType := ReServiceType{}

	expErrors = map[string][]string{
		"code": {"Code can not be blank."},
		"name": {"Name can not be blank."},
	}
	suite.verifyValidationErrors(&invalidReServiceType, expErrors)
}

func (suite *ModelSuite) Test_ReServiceTypeCreateAndSave() {
	validReServiceType := ReServiceType{
		Code: "123abc",
		Name: "California",
	}

	suite.MustSave(&validReServiceType)
}
