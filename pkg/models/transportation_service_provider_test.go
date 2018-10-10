package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func CreateTestTsp(suite *ModelSuite) TransportationServiceProvider {
	tsp := TransportationServiceProvider{
		StandardCarrierAlphaCode: "TRSS",
		Enrolled:                 true,
	}
	suite.mustSave(&tsp)
	return tsp
}

func (suite *ModelSuite) Test_TransportationServiceProvider() {
	tsp := &TransportationServiceProvider{}

	expErrors := map[string][]string{
		"standard_carrier_alpha_code": []string{"StandardCarrierAlphaCode can not be blank."},
	}

	suite.verifyValidationErrors(tsp, expErrors)
}
