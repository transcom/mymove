package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_TransportationServiceProvider() {
	tsp := &TransportationServiceProvider{}

	expErrors := map[string][]string{
		"standard_carrier_alpha_code": {"StandardCarrierAlphaCode can not be blank."},
	}

	suite.verifyValidationErrors(tsp, expErrors)
}
