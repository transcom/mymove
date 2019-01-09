package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func CreateTestTsp(suite *ModelSuite) TransportationServiceProvider {
	name := "Truss Transit Services"
	pocGeneralName := "Joey Joe-Joe Schabadoo"
	pocGeneralEmail := "joey@example.com"
	pocGeneralPhone := "(555) 867-5309"
	pocClaimsName := "Claimy Claimer"
	pocClaimsEmail := "claims@example.com"
	pocClaimsPhone := "(555) 123-4567"

	tsp := TransportationServiceProvider{
		StandardCarrierAlphaCode: "TRSS",
		Enrolled:                 true,
		Name:                     &name,
		PocGeneralName:           &pocGeneralName,
		PocGeneralEmail:          &pocGeneralEmail,
		PocGeneralPhone:          &pocGeneralPhone,
		PocClaimsName:            &pocClaimsName,
		PocClaimsEmail:           &pocClaimsEmail,
		PocClaimsPhone:           &pocClaimsPhone,
	}
	suite.MustSave(&tsp)
	return tsp
}

func (suite *ModelSuite) Test_TransportationServiceProvider() {
	tsp := &TransportationServiceProvider{}

	expErrors := map[string][]string{
		"standard_carrier_alpha_code": []string{"StandardCarrierAlphaCode can not be blank."},
	}

	suite.verifyValidationErrors(tsp, expErrors)
}
