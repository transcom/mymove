package models_test

import (
	"github.com/gofrs/uuid"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestAdminUserCreation() {
	t := suite.T()

	user := testdatagen.MakeDefaultUser(suite.DB())

	newAdminUser := AdminUser{
		FirstName: "Leo",
		LastName:  "Spaceman",
		UserID:    &user.ID,
		Role:      "SYSTEM_ADMIN",
		Email:     "leo@gmail.com",
	}

	if verrs, err := suite.DB().ValidateAndCreate(&newAdminUser); err != nil || verrs.HasAny() {
		t.Fatal("Didn't create admin user in db.")
	}

	if newAdminUser.ID == uuid.Nil {
		t.Error("Didn't get an id back for admin user.")
	}
}

func (suite *ModelSuite) TestAdminUserCreationWithoutValues() {
	newAdminUser := &AdminUser{}

	expErrors := map[string][]string{
		"first_name": {"FirstName can not be blank."},
		"last_name":  {"LastName can not be blank."},
		"email":      {"Email can not be blank."},
		"role":       {"Role is not in the list [SYSTEM_ADMIN PROGRAM_ADMIN]."},
	}

	suite.verifyValidationErrors(newAdminUser, expErrors)
}
