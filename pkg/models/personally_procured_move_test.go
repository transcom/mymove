package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPPMValidation() {
	ppm := &PersonallyProcuredMove{}

	expErrors := map[string][]string{
		"status": {"Status can not be blank."},
	}

	suite.verifyValidationErrors(ppm, expErrors)
}
