package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ShippingAgentValidations() {
	shippingAgent := &ShippingAgent{}

	expErrors := map[string][]string{
		"shipment_id":  {"ShipmentID can not be blank."},
		"agent_type":   {"AgentType can not be blank."},
		"name":         {"Name can not be blank."},
		"phone_number": {"PhoneNumber can not be blank."},
		"email":        {"Email can not be blank."},
	}

	suite.verifyValidationErrors(shippingAgent, expErrors)
}
