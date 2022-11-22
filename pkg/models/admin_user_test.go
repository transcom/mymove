package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestAdminUserCreation() {
	user := factory.BuildUser(nil, nil, nil)

	newAdminUser := AdminUser{
		FirstName: "Leo",
		LastName:  "Spaceman",
		UserID:    &user.ID,
		Role:      "SYSTEM_ADMIN",
		Email:     "leo@gmail.com",
	}

	verrs, err := newAdminUser.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
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
