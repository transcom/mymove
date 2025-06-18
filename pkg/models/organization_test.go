package models_test

import (
	m "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestOrganizationValidation() {
	email := "test@truss.works"
	phone := "9144825484"

	newOrganization := m.Organization{
		Name:     "Truss",
		PocEmail: &email,
		PocPhone: &phone,
	}

	verrs, err := newOrganization.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) TestOrganizationCreationWithoutValues() {
	newOrganization := &m.Organization{}

	expErrors := map[string][]string{
		"name": {"Name can not be blank."},
	}

	suite.verifyValidationErrors(newOrganization, expErrors, nil)
}
