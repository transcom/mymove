package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ReRateAreaValidation() {
	validReRateArea := ReRateArea{
		IsOconus: true,
		Code:     "123abc",
		Name:     "California",
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validReRateArea, expErrors)

	invalidReRateArea := ReRateArea{}

	expErrors = map[string][]string{
		"code": {"Code can not be blank."},
		"name": {"Name can not be blank."},
	}
	suite.verifyValidationErrors(&invalidReRateArea, expErrors)
}

func (suite *ModelSuite) Test_RateAreaCreateAndSave() {
	validReRateArea := ReRateArea{
		IsOconus: true,
		Code:     "123abc",
		Name:     "California",
	}

	suite.MustSave(&validReRateArea)
}
