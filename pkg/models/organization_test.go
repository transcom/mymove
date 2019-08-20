package models_test

import (
	"github.com/gofrs/uuid"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestOrganizationCreation() {
	t := suite.T()

	newOrganization := Organization{
		Name:     "Truss",
		PocEmail: "test@truss.works",
		PocPhone: "9144825484",
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
		"name":      {"Name can not be blank."},
		"poc_email": {"PocEmail can not be blank."},
		"poc_phone": {"PocPhone can not be blank."},
	}

	suite.verifyValidationErrors(newOrganization, expErrors)
}
