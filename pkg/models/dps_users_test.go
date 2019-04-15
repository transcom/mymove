package models_test

import (
	"strings"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_DpsUserCreate() {
	t := suite.T()

	dpsUser := models.DpsUser{
		LoginGovEmail: "test@example.com",
	}

	verrs, err := suite.DB().ValidateAndSave(&dpsUser)

	if err != nil {
		t.Fatalf("could not save DPS user: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}
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
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(user.LoginGovEmail, email)
}
