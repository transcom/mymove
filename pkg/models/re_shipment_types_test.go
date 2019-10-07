package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ReShipmentTypeValidation() {
	validReShipmentType := ReShipmentType{
		Code: "123abc",
		Name: "California",
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validReShipmentType, expErrors)

	invalidReShipmentType := ReShipmentType{}

	expErrors = map[string][]string{
		"code": {"Code can not be blank."},
		"name": {"Name can not be blank."},
	}
	suite.verifyValidationErrors(&invalidReShipmentType, expErrors)
}

func (suite *ModelSuite) Test_ReShipmentTypeCreateAndSave() {
	validReShipmentType := ReShipmentType{
		Code: "123abc",
		Name: "California",
	}

	suite.MustSave(&validReShipmentType)
}
