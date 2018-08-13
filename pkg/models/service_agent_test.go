package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ServiceAgentValidations() {
	serviceAgent := &ServiceAgent{}

	expErrors := map[string][]string{
		"shipment_id":      {"ShipmentID can not be blank."},
		"role":             {"Role can not be blank."},
		"point_of_contact": {"Point of Contact can not be blank."},
	}

	suite.verifyValidationErrors(serviceAgent, expErrors)
}
