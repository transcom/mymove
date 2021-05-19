package models_test

import (
	"strings"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_DpsUserCreate() {
	dpsUser := models.DpsUser{
		LoginGovEmail: "test@example.com",
	}

	verrs, err := dpsUser.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) Test_DpsUserValidations() {
	document := &models.DpsUser{}

	var expErrors = map[string][]string{
		"login_gov_email": {"LoginGovEmail can not be blank."},
	}

	suite.verifyValidationErrors(document, expErrors)
}

func (suite *ModelSuite) TestFetchDPSUserByEmailCaseSensitivity() {
	email := "Test@example.com"

	dpsUser := models.DpsUser{
		LoginGovEmail: email,
	}

	suite.MustSave(&dpsUser)
	user, err := models.FetchDPSUserByEmail(suite.DB(), strings.ToLower(email))
	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(user.LoginGovEmail, email)
}
