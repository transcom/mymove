package models_test

import (
	"github.com/gofrs/uuid"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestOrganizationCreation() {
	t := suite.T()

	email := "test@truss.works"
	phone := "9144825484"

	newOrganization := Organization{
		Name:     "Truss",
		PocEmail: &email,
		PocPhone: &phone,
	}

	if verrs, err := suite.DB().ValidateAndCreate(&newOrganization); err != nil || verrs.HasAny() {
		t.Fatal("Didn't create admin user in db.")
	}

	if newOrganization.ID == uuid.Nil {
		t.Error("Didn't get an id back for admin user.")
	}
}

func (suite *ModelSuite) TestOrganizationCreationWithoutValues() {
	newOrganization := &Organization{}

	expErrors := map[string][]string{
		"name": {"Name can not be blank."},
	}

	suite.verifyValidationErrors(newOrganization, expErrors)
}
