package models_test

import (
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
